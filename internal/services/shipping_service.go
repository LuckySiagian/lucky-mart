package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"backend-sembako/internal/dto"
	"backend-sembako/internal/models"
	"backend-sembako/internal/repositories"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrShipmentNotFound      = errors.New("shipment record not found")
	ErrInvalidStatusTransition = errors.New("invalid shipment status transition")
)

type ShippingService interface {
	GetShippingMethods(ctx context.Context) ([]models.ShippingMethod, error)
	CalculateShippingRate(ctx context.Context, req dto.CalculateShippingDTO) (*dto.ShippingRateResponseDTO, error)
	CreateShipment(ctx context.Context, req dto.CreateShipmentDTO) (*dto.ShipmentResponseDTO, error)
	GetShipmentByOrderID(ctx context.Context, orderID uuid.UUID) (*dto.ShipmentResponseDTO, error)
	UpdateShipmentStatus(ctx context.Context, shipmentID uuid.UUID, req dto.UpdateShipmentStatusDTO) (*dto.ShipmentResponseDTO, error)
}

type shippingService struct {
	db             *gorm.DB
	courierRepo    repositories.CourierRepository
	methodRepo     repositories.ShippingMethodRepository
	rateRepo       repositories.ShippingRateRepository
	shipmentRepo   repositories.ShipmentRepository
	orderRepo      repositories.OrderRepository
}

func NewShippingService(
	db *gorm.DB,
	courierRepo repositories.CourierRepository,
	methodRepo repositories.ShippingMethodRepository,
	rateRepo repositories.ShippingRateRepository,
	shipmentRepo repositories.ShipmentRepository,
	orderRepo repositories.OrderRepository,
) ShippingService {
	return &shippingService{
		db:           db,
		courierRepo:  courierRepo,
		methodRepo:   methodRepo,
		rateRepo:     rateRepo,
		shipmentRepo: shipmentRepo,
		orderRepo:    orderRepo,
	}
}

func (s *shippingService) GetShippingMethods(ctx context.Context) ([]models.ShippingMethod, error) {
	return s.methodRepo.FindAll(ctx)
}

func (s *shippingService) CalculateShippingRate(ctx context.Context, req dto.CalculateShippingDTO) (*dto.ShippingRateResponseDTO, error) {
	method, err := s.methodRepo.FindByID(ctx, req.ShippingMethodID)
	if err != nil || method == nil {
		return nil, errors.New("shipping method not found")
	}

	rate, err := s.rateRepo.FindByMethodAndCity(ctx, req.ShippingMethodID, req.City)
	flatFee := 15000.00
	ratePerKg := 5000.00

	if rate != nil {
		flatFee = rate.FlatRate
		ratePerKg = rate.RatePerKg
	}

	calculatedFee := method.BaseRate + flatFee + (req.TotalWeightKg * ratePerKg)
	courierName := "Kurir Toko"
	if method.Courier != nil {
		courierName = method.Courier.Name
	}

	return &dto.ShippingRateResponseDTO{
		ShippingMethodID: method.ID,
		CourierName:      courierName,
		MethodName:       method.Name,
		EstimatedDays:    method.EstimatedDays,
		CalculatedFee:    calculatedFee,
	}, nil
}

func (s *shippingService) CreateShipment(ctx context.Context, req dto.CreateShipmentDTO) (*dto.ShipmentResponseDTO, error) {
	order, err := s.orderRepo.FindByID(ctx, req.OrderID)
	if err != nil || order == nil {
		return nil, errors.New("order not found")
	}

	trackingNum := fmt.Sprintf("RESI-%s-%04d", time.Now().Format("20060102"), time.Now().UnixNano()%10000)
	var shipment *models.Shipment

	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		shipment = &models.Shipment{
			TrackingNumber:   trackingNum,
			OrderID:          req.OrderID,
			CourierID:        req.CourierID,
			ShippingMethodID: req.ShippingMethodID,
			Status:           "PENDING",
			ShippingFee:      req.ShippingFee,
		}

		if err := s.shipmentRepo.Create(ctx, tx, shipment); err != nil {
			return err
		}

		initialTracking := &models.ShipmentTracking{
			ShipmentID: shipment.ID,
			Status:     "PENDING",
			Location:   "Gudang Pusat Sembako",
			Notes:      "Pesanan siap diproses untuk pengiriman",
		}
		return s.shipmentRepo.AddTracking(ctx, tx, initialTracking)
	})

	if err != nil {
		return nil, err
	}

	return s.GetShipmentByOrderID(ctx, req.OrderID)
}

func (s *shippingService) GetShipmentByOrderID(ctx context.Context, orderID uuid.UUID) (*dto.ShipmentResponseDTO, error) {
	shipment, err := s.shipmentRepo.FindByOrderID(ctx, orderID)
	if err != nil || shipment == nil {
		return nil, ErrShipmentNotFound
	}
	res := mapShipmentToDTO(shipment)
	return &res, nil
}

func (s *shippingService) UpdateShipmentStatus(ctx context.Context, shipmentID uuid.UUID, req dto.UpdateShipmentStatusDTO) (*dto.ShipmentResponseDTO, error) {
	shipment, err := s.shipmentRepo.FindByID(ctx, shipmentID)
	if err != nil || shipment == nil {
		return nil, ErrShipmentNotFound
	}

	now := time.Now()
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		shipment.Status = req.Status
		if req.Status == "IN_TRANSIT" || req.Status == "OUT_FOR_DELIVERY" {
			shipment.ShippedAt = &now
		} else if req.Status == "DELIVERED" {
			shipment.DeliveredAt = &now
		}

		if err := s.shipmentRepo.Update(ctx, tx, shipment); err != nil {
			return err
		}

		tracking := &models.ShipmentTracking{
			ShipmentID: shipment.ID,
			Status:     req.Status,
			Location:   req.Location,
			Notes:      req.Notes,
		}
		return s.shipmentRepo.AddTracking(ctx, tx, tracking)
	})

	if err != nil {
		return nil, err
	}

	return s.GetShipmentByOrderID(ctx, shipment.OrderID)
}

func mapShipmentToDTO(s *models.Shipment) dto.ShipmentResponseDTO {
	courierName := "Kurir Internal"
	if s.Courier != nil {
		courierName = s.Courier.Name
	}
	methodName := "Pengiriman Reguler"
	if s.ShippingMethod != nil {
		methodName = s.ShippingMethod.Name
	}

	trackings := make([]dto.ShipmentTrackingResponseDTO, len(s.Trackings))
	for i, t := range s.Trackings {
		trackings[i] = dto.ShipmentTrackingResponseDTO{
			ID:        t.ID,
			Status:    t.Status,
			Location:  t.Location,
			Notes:     t.Notes,
			CreatedAt: t.CreatedAt,
		}
	}

	return dto.ShipmentResponseDTO{
		ID:             s.ID,
		TrackingNumber: s.TrackingNumber,
		OrderID:        s.OrderID,
		CourierName:    courierName,
		MethodName:     methodName,
		Status:         s.Status,
		ShippingFee:    s.ShippingFee,
		ShippedAt:      s.ShippedAt,
		DeliveredAt:    s.DeliveredAt,
		Trackings:      trackings,
		CreatedAt:      s.CreatedAt,
	}
}
