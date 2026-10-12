package services

import (
	"context"
	"errors"
	"time"

	"backend-sembako/internal/constants"
	"backend-sembako/internal/dto"
	"backend-sembako/internal/models"
	"backend-sembako/internal/repositories"
	"backend-sembako/pkg/eventbus"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrShipmentNotFound        = errors.New("data pengiriman tidak ditemukan")
	ErrShipmentAlreadyExists   = errors.New("pesanan ini sudah punya data pengiriman")
	ErrShipmentFinished        = errors.New("pengiriman sudah selesai dan tidak bisa diubah lagi")
	ErrOrderNotReadyToShip     = errors.New("pesanan belum siap dikirim (harus sudah dikonfirmasi atau sedang dikemas)")
	ErrUnauthorizedShipment    = errors.New("anda tidak berhak melihat pengiriman ini")
	ErrInvalidStatusTransition = errors.New("perubahan status pengiriman tidak valid")
)

type ShippingService interface {
	GetShippingMethods(ctx context.Context) ([]dto.ShippingMethodDTO, error)
	CalculateShippingRate(ctx context.Context, req dto.CalculateShippingDTO) (*dto.ShippingRateResponseDTO, error)
	GetShipments(ctx context.Context, req dto.PaginationRequest, status string) (dto.PaginatedResponse, error)
	CreateShipment(ctx context.Context, req dto.CreateShipmentDTO) (*dto.ShipmentResponseDTO, error)
	GetShipmentByOrderID(ctx context.Context, userID uuid.UUID, role constants.UserRole, orderID uuid.UUID) (*dto.ShipmentResponseDTO, error)
	UpdateShipmentStatus(ctx context.Context, shipmentID uuid.UUID, req dto.UpdateShipmentStatusDTO) (*dto.ShipmentResponseDTO, error)
}

type shippingService struct {
	db           *gorm.DB
	courierRepo  repositories.CourierRepository
	methodRepo   repositories.ShippingMethodRepository
	rateRepo     repositories.ShippingRateRepository
	shipmentRepo repositories.ShipmentRepository
	orderRepo    repositories.OrderRepository
	bus          eventbus.EventBus
}

func NewShippingService(
	db *gorm.DB,
	courierRepo repositories.CourierRepository,
	methodRepo repositories.ShippingMethodRepository,
	rateRepo repositories.ShippingRateRepository,
	shipmentRepo repositories.ShipmentRepository,
	orderRepo repositories.OrderRepository,
	bus eventbus.EventBus,
) ShippingService {
	return &shippingService{
		db:           db,
		courierRepo:  courierRepo,
		methodRepo:   methodRepo,
		rateRepo:     rateRepo,
		shipmentRepo: shipmentRepo,
		orderRepo:    orderRepo,
		bus:          bus,
	}
}

func (s *shippingService) GetShippingMethods(ctx context.Context) ([]dto.ShippingMethodDTO, error) {
	methods, err := s.methodRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	res := make([]dto.ShippingMethodDTO, len(methods))
	for i, m := range methods {
		res[i] = dto.ShippingMethodDTO{
			ID:            m.ID,
			Code:          m.Code,
			Name:          m.Name,
			Description:   m.Description,
			EstimatedDays: m.EstimatedDays,
			BaseRate:      m.BaseRate,
			CourierID:     m.CourierID,
		}
		if m.Courier != nil {
			res[i].CourierName = m.Courier.Name
		}
	}
	return res, nil
}

func (s *shippingService) CalculateShippingRate(ctx context.Context, req dto.CalculateShippingDTO) (*dto.ShippingRateResponseDTO, error) {
	method, fee, err := quoteShippingFee(ctx, s.methodRepo, s.rateRepo, req.ShippingMethodID, req.City)
	if err != nil {
		return nil, err
	}

	courierName := "Kurir Toko"
	if method.Courier != nil {
		courierName = method.Courier.Name
	}

	return &dto.ShippingRateResponseDTO{
		ShippingMethodID: method.ID,
		CourierName:      courierName,
		MethodName:       method.Name,
		EstimatedDays:    method.EstimatedDays,
		CalculatedFee:    fee,
	}, nil
}

func (s *shippingService) GetShipments(ctx context.Context, req dto.PaginationRequest, status string) (dto.PaginatedResponse, error) {
	shipments, total, err := s.shipmentRepo.FindAll(ctx, req, status)
	if err != nil {
		return dto.PaginatedResponse{}, err
	}

	dtos := make([]dto.ShipmentResponseDTO, len(shipments))
	for i := range shipments {
		dtos[i] = mapShipmentToDTO(&shipments[i])
	}

	req.SetDefaults()
	return dto.NewPaginatedResponse(dtos, req.Page, req.Limit, total), nil
}

func (s *shippingService) CreateShipment(ctx context.Context, req dto.CreateShipmentDTO) (*dto.ShipmentResponseDTO, error) {
	order, err := s.orderRepo.FindByID(ctx, req.OrderID)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, ErrOrderNotFound
	}
	if order.Status != constants.OrderStatusConfirmed && order.Status != constants.OrderStatusProcessingPacking && order.Status != constants.OrderStatusPaid {
		return nil, ErrOrderNotReadyToShip
	}
	if existing, err := s.shipmentRepo.FindByOrderID(ctx, order.ID); err != nil {
		return nil, err
	} else if existing != nil {
		return nil, ErrShipmentAlreadyExists
	}

	method, _, err := quoteShippingFee(ctx, s.methodRepo, s.rateRepo, req.ShippingMethodID, "")
	if err != nil {
		return nil, err
	}
	courierID := method.CourierID
	if req.CourierID != nil && *req.CourierID != uuid.Nil {
		courier, err := s.courierRepo.FindByID(ctx, *req.CourierID)
		if err != nil {
			return nil, err
		}
		if courier == nil {
			return nil, errors.New("kurir tidak ditemukan")
		}
		courierID = courier.ID
	}
	fee := order.ShippingFee
	if req.ShippingFee != nil {
		fee = *req.ShippingFee
	}

	note := req.Notes
	if note == "" {
		note = "Pesanan siap diproses untuk pengiriman"
	}

	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		shipment := &models.Shipment{
			TrackingNumber:   newReference("RESI"),
			OrderID:          order.ID,
			CourierID:        courierID,
			ShippingMethodID: method.ID,
			Status:           "READY_TO_SHIP",
			ShippingFee:      fee,
		}
		if err := s.shipmentRepo.Create(ctx, tx, shipment); err != nil {
			return err
		}
		return s.shipmentRepo.AddTracking(ctx, tx, &models.ShipmentTracking{
			ShipmentID: shipment.ID, Status: "READY_TO_SHIP", Location: "Toko", Notes: note,
		})
	})
	if err != nil {
		return nil, err
	}

	shipment, err := s.shipmentRepo.FindByOrderID(ctx, order.ID)
	if err != nil || shipment == nil {
		return nil, ErrShipmentNotFound
	}
	res := mapShipmentToDTO(shipment)
	return &res, nil
}

func (s *shippingService) GetShipmentByOrderID(ctx context.Context, userID uuid.UUID, role constants.UserRole, orderID uuid.UUID) (*dto.ShipmentResponseDTO, error) {
	shipment, err := s.shipmentRepo.FindByOrderID(ctx, orderID)
	if err != nil {
		return nil, err
	}
	if shipment == nil {
		return nil, ErrShipmentNotFound
	}
	if !role.IsStaff() && (shipment.Order == nil || shipment.Order.UserID != userID) {
		return nil, ErrUnauthorizedShipment
	}
	res := mapShipmentToDTO(shipment)
	return &res, nil
}

func (s *shippingService) UpdateShipmentStatus(ctx context.Context, shipmentID uuid.UUID, req dto.UpdateShipmentStatusDTO) (*dto.ShipmentResponseDTO, error) {
	shipment, err := s.shipmentRepo.FindByID(ctx, shipmentID)
	if err != nil {
		return nil, err
	}
	if shipment == nil {
		return nil, ErrShipmentNotFound
	}
	if shipment.Status == "DELIVERED" || shipment.Status == "RETURNED" {
		return nil, ErrShipmentFinished
	}
	if shipment.Status == req.Status {
		return nil, ErrInvalidStatusTransition
	}

	order, err := s.orderRepo.FindByID(ctx, shipment.OrderID)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, ErrOrderNotFound
	}
	if order.Status == constants.OrderStatusCancelled {
		return nil, errors.New("pesanan sudah dibatalkan, pengiriman tidak bisa diperbarui")
	}

	location := req.Location
	if location == "" {
		location = "-"
	}

	now := time.Now()
	completed := false
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		shipment.Status = req.Status
		switch req.Status {
		case "PICKED_UP", "IN_TRANSIT", "OUT_FOR_DELIVERY":
			if shipment.ShippedAt == nil {
				shipment.ShippedAt = &now
			}
		case "DELIVERED":
			shipment.DeliveredAt = &now
		}
		if err := s.shipmentRepo.Update(ctx, tx, shipment); err != nil {
			return err
		}
		if err := s.shipmentRepo.AddTracking(ctx, tx, &models.ShipmentTracking{
			ShipmentID: shipment.ID, Status: req.Status, Location: location, Notes: req.Notes,
		}); err != nil {
			return err
		}

		// Keep the order status in step with the shipment.
		switch req.Status {
		case "PICKED_UP", "IN_TRANSIT", "OUT_FOR_DELIVERY":
			if order.Status == constants.OrderStatusConfirmed || order.Status == constants.OrderStatusProcessingPacking || order.Status == constants.OrderStatusPaid {
				return setOrderStatusTx(tx, order, constants.OrderStatusOutForDelivery, "Pesanan sedang diantar kurir")
			}
		case "DELIVERED":
			if order.Status != constants.OrderStatusCompleted {
				completed = true
				return completeOrderTx(tx, order, "Pesanan telah diterima pelanggan")
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	if completed {
		publishEvent(s.bus, "ORDER_COMPLETED", "ORDER", order.ID, map[string]interface{}{
			"user_id":      order.UserID.String(),
			"order_number": order.OrderNumber,
			"title":        "Pesanan " + order.OrderNumber + " selesai",
			"message":      "Pesanan Anda sudah diterima. Terima kasih sudah berbelanja di Lucky Mart!",
		})
	} else if req.Status == "OUT_FOR_DELIVERY" {
		publishEvent(s.bus, "SHIPMENT_OUT_FOR_DELIVERY", "SHIPMENT", shipment.ID, map[string]interface{}{
			"user_id":      order.UserID.String(),
			"order_number": order.OrderNumber,
			"title":        "Pesanan " + order.OrderNumber + " sedang diantar",
			"message":      "Kurir sedang menuju alamat Anda. Nomor resi: " + shipment.TrackingNumber,
		})
	}

	updated, err := s.shipmentRepo.FindByID(ctx, shipmentID)
	if err != nil || updated == nil {
		return nil, ErrShipmentNotFound
	}
	res := mapShipmentToDTO(updated)
	return &res, nil
}

func mapShipmentToDTO(s *models.Shipment) dto.ShipmentResponseDTO {
	res := dto.ShipmentResponseDTO{
		ID:             s.ID,
		TrackingNumber: s.TrackingNumber,
		OrderID:        s.OrderID,
		CourierName:    "Kurir Toko",
		MethodName:     "Pengiriman Reguler",
		Status:         s.Status,
		ShippingFee:    s.ShippingFee,
		ShippedAt:      s.ShippedAt,
		DeliveredAt:    s.DeliveredAt,
		CreatedAt:      s.CreatedAt,
	}
	if s.Courier != nil {
		res.CourierName = s.Courier.Name
	}
	if s.ShippingMethod != nil {
		res.MethodName = s.ShippingMethod.Name
	}
	if s.Order != nil {
		res.OrderNumber = s.Order.OrderNumber
		res.OrderStatus = string(s.Order.Status)
		if s.Order.User != nil {
			res.CustomerName = s.Order.User.FullName
		}
		if a := s.Order.OrderAddress; a != nil {
			res.RecipientName = a.RecipientName
			res.RecipientPhone = a.Phone
			res.Address = a.FullAddress + ", " + a.City
		}
	}

	res.Trackings = make([]dto.ShipmentTrackingResponseDTO, len(s.Trackings))
	for i, t := range s.Trackings {
		res.Trackings[i] = dto.ShipmentTrackingResponseDTO{
			ID: t.ID, Status: t.Status, Location: t.Location, Notes: t.Notes, CreatedAt: t.CreatedAt,
		}
	}
	return res
}
