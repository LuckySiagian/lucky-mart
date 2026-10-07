package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"backend-sembako/internal/constants"
	"backend-sembako/internal/dto"
	"backend-sembako/internal/models"
	"backend-sembako/internal/repositories"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	ErrOrderNotFound      = errors.New("order not found")
	ErrEmptyCartForOrder  = errors.New("cannot create order from an empty shopping cart")
	ErrOrderCannotCancel  = errors.New("order status cannot be cancelled at its current state")
	ErrUnauthorizedOrder  = errors.New("you do not have permission to view or modify this order")
)

type OrderService interface {
	CreateOrder(ctx context.Context, userID uuid.UUID, userRole constants.UserRole, req dto.CreateOrderDTO) (*dto.OrderResponseDTO, error)
	GetOrders(ctx context.Context, userID uuid.UUID, userRole constants.UserRole, req dto.PaginationRequest) (dto.PaginatedResponse, error)
	GetOrderByID(ctx context.Context, userID uuid.UUID, userRole constants.UserRole, id uuid.UUID) (*dto.OrderResponseDTO, error)
	CancelOrder(ctx context.Context, userID uuid.UUID, userRole constants.UserRole, id uuid.UUID, req dto.CancelOrderDTO) (*dto.OrderResponseDTO, error)
	GetOrderHistory(ctx context.Context, userID uuid.UUID, userRole constants.UserRole, id uuid.UUID) ([]dto.OrderHistoryResponseDTO, error)
	GetInvoiceByID(ctx context.Context, userID uuid.UUID, userRole constants.UserRole, invoiceID uuid.UUID) (*dto.InvoiceResponseDTO, error)
	UpdateOrderStatus(ctx context.Context, userID uuid.UUID, userRole constants.UserRole, id uuid.UUID, req dto.UpdateOrderStatusDTO) (*dto.OrderResponseDTO, error)
}

type orderService struct {
	db              *gorm.DB
	orderRepo       repositories.OrderRepository
	invoiceRepo     repositories.InvoiceRepository
	cartRepo        repositories.CartRepository
	cartItemRepo    repositories.CartItemRepository
	productRepo     repositories.ProductRepository
	inventoryRepo   repositories.InventoryRepository
	stockMoveRepo   repositories.StockMovementRepository
	b2bProfileRepo  repositories.B2BProfileRepository
}

func NewOrderService(
	db *gorm.DB,
	orderRepo repositories.OrderRepository,
	invoiceRepo repositories.InvoiceRepository,
	cartRepo repositories.CartRepository,
	cartItemRepo repositories.CartItemRepository,
	productRepo repositories.ProductRepository,
	inventoryRepo repositories.InventoryRepository,
	stockMoveRepo repositories.StockMovementRepository,
	b2bProfileRepo repositories.B2BProfileRepository,
) OrderService {
	return &orderService{
		db:              db,
		orderRepo:       orderRepo,
		invoiceRepo:     invoiceRepo,
		cartRepo:        cartRepo,
		cartItemRepo:    cartItemRepo,
		productRepo:     productRepo,
		inventoryRepo:   inventoryRepo,
		stockMoveRepo:   stockMoveRepo,
		b2bProfileRepo:  b2bProfileRepo,
	}
}

func (s *orderService) CreateOrder(ctx context.Context, userID uuid.UUID, userRole constants.UserRole, req dto.CreateOrderDTO) (*dto.OrderResponseDTO, error) {
	cart, err := s.cartRepo.FindByUserID(ctx, userID)
	if err != nil || cart == nil || len(cart.Items) == 0 {
		return nil, ErrEmptyCartForOrder
	}

	orderNum := generateOrderNumber()
	now := time.Now()
	var createdOrder *models.Order

	// Execute ACID Database Transaction
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var totalAmount float64
		var orderItems []models.OrderItem

		for _, item := range cart.Items {
			// RACE CONDITION PREVENTION: Database Row Locking (SELECT ... FOR UPDATE)
			var productUnit models.ProductUnit
			if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
				Where("id = ?", item.ProductUnitID).
				First(&productUnit).Error; err != nil {
				return fmt.Errorf("failed to lock product unit stock: %w", err)
			}

			// Validate Product Active Status
			var product models.Product
			if err := tx.Where("id = ?", item.ProductID).First(&product).Error; err != nil || !product.IsActive {
				return fmt.Errorf("product '%s' is inactive or deleted", product.Name)
			}

			// Validate Available Stock
			if item.Quantity > productUnit.StockQuantity {
				return fmt.Errorf("insufficient stock for item '%s'. Requested: %d, Available: %d", product.Name, item.Quantity, productUnit.StockQuantity)
			}

			// Determine Price Snapshot (Dynamic Wholesale vs Retail)
			priceSnapshot := productUnit.RetailPrice
			if userRole == constants.RoleB2BCustomer && item.Quantity >= productUnit.MinWholesaleQty {
				priceSnapshot = productUnit.WholesalePrice
			}
			subtotal := float64(item.Quantity) * priceSnapshot
			totalAmount += subtotal

			// 1. Deduct ProductUnit Stock
			newUnitStock := productUnit.StockQuantity - item.Quantity
			if err := tx.Model(&models.ProductUnit{}).
				Where("id = ?", productUnit.ID).
				Update("stock_quantity", newUnitStock).Error; err != nil {
				return err
			}

			// 2. Deduct Inventory Stock
			var inv models.Inventory
			if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
				Where("product_unit_id = ?", productUnit.ID).
				First(&inv).Error; err == nil {
				qtyBefore := inv.CurrentStock
				qtyAfter := qtyBefore - item.Quantity
				inv.CurrentStock = qtyAfter
				_ = tx.Save(&inv)

				// 3. Mandatory StockMovement Log (STOCK_OUT)
				movement := &models.StockMovement{
					InventoryID:     inv.ID,
					ProductUnitID:   productUnit.ID,
					MovementType:    "STOCK_OUT",
					ReferenceNumber: orderNum,
					QuantityBefore:  qtyBefore,
					QuantityChanged: -item.Quantity,
					QuantityAfter:   qtyAfter,
					Notes:           fmt.Sprintf("Penjualan Order %s", orderNum),
				}
				_ = s.stockMoveRepo.Create(ctx, tx, movement)
			}

			// Snapshot Order Item
			orderItems = append(orderItems, models.OrderItem{
				ProductID:     item.ProductID,
				ProductUnitID: item.ProductUnitID,
				Quantity:      item.Quantity,
				UnitPrice:     priceSnapshot,
				Subtotal:      subtotal,
				Notes:         item.Notes,
			})
		}

		shippingFee := 15000.00
		grandTotal := totalAmount + shippingFee

		// Calculate Due Date
		var dueDate *time.Time
		dueTime := now.Add(24 * time.Hour)
		if req.PaymentMethod == constants.PaymentMethodB2BCreditTempo {
			dueTime = now.Add(30 * 24 * time.Hour)
		}
		dueDate = &dueTime

		// Create Header Order Entity
		order := &models.Order{
			OrderNumber:       orderNum,
			UserID:            userID,
			ShippingAddressID: uuid.New(), // Fixed UUID reference for snapshot
			Status:            constants.OrderStatusCreated,
			PaymentMethod:     req.PaymentMethod,
			TotalAmount:       totalAmount,
			ShippingFee:       shippingFee,
			GrandTotal:        grandTotal,
			DueDate:           dueDate,
			OrderDate:         now,
			OrderItems:        orderItems,
		}

		if err := s.orderRepo.Create(ctx, tx, order); err != nil {
			return err
		}

		// Create Order Address Snapshot
		orderAddr := &models.OrderAddress{
			OrderID:       order.ID,
			RecipientName: req.RecipientName,
			Phone:         req.Phone,
			FullAddress:   req.FullAddress,
			City:          req.City,
			PostalCode:    req.PostalCode,
		}
		if err := tx.Create(orderAddr).Error; err != nil {
			return err
		}

		// Generate Invoice
		invoiceNum := fmt.Sprintf("INV-%s", orderNum[4:])
		invoice := &models.Invoice{
			InvoiceNumber: invoiceNum,
			OrderID:       order.ID,
			UserID:        userID,
			TotalAmount:   grandTotal,
			Status:        "UNPAID",
			DueDate:       dueTime,
			IssuedAt:      now,
		}
		if err := s.invoiceRepo.Create(ctx, tx, invoice); err != nil {
			return err
		}

		// Initial Order History Record
		history := &models.OrderHistory{
			OrderID: order.ID,
			Status:  constants.OrderStatusCreated,
			Notes:   "Pesanan berhasil dibuat. Menunggu konfirmasi pembayaran.",
		}
		if err := s.orderRepo.AddHistory(ctx, tx, history); err != nil {
			return err
		}

		// Clear Customer Cart
		_ = s.cartRepo.Clear(ctx, cart.ID)

		createdOrder = order
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("transaction order creation failed: %w", err)
	}

	return s.GetOrderByID(ctx, userID, userRole, createdOrder.ID)
}

func (s *orderService) GetOrders(ctx context.Context, userID uuid.UUID, userRole constants.UserRole, req dto.PaginationRequest) (dto.PaginatedResponse, error) {
	var orders []models.Order
	var total int64
	var err error

	if userRole == constants.RoleSuperAdmin || userRole == constants.RoleStoreManager || userRole == constants.RoleWarehouseStaff {
		orders, total, err = s.orderRepo.FindAll(ctx, req)
	} else {
		orders, total, err = s.orderRepo.FindByUserID(ctx, userID, req)
	}

	if err != nil {
		return dto.PaginatedResponse{}, err
	}

	dtos := make([]dto.OrderResponseDTO, len(orders))
	for i, o := range orders {
		dtos[i] = mapOrderToDTO(&o)
	}

	req.SetDefaults()
	return dto.NewPaginatedResponse(dtos, req.Page, req.Limit, total), nil
}

func (s *orderService) GetOrderByID(ctx context.Context, userID uuid.UUID, userRole constants.UserRole, id uuid.UUID) (*dto.OrderResponseDTO, error) {
	order, err := s.orderRepo.FindByID(ctx, id)
	if err != nil || order == nil {
		return nil, ErrOrderNotFound
	}

	if userRole != constants.RoleSuperAdmin && userRole != constants.RoleStoreManager && userRole != constants.RoleWarehouseStaff {
		if order.UserID != userID {
			return nil, ErrUnauthorizedOrder
		}
	}

	res := mapOrderToDTO(order)
	return &res, nil
}

func (s *orderService) CancelOrder(ctx context.Context, userID uuid.UUID, userRole constants.UserRole, id uuid.UUID, req dto.CancelOrderDTO) (*dto.OrderResponseDTO, error) {
	order, err := s.orderRepo.FindByID(ctx, id)
	if err != nil || order == nil {
		return nil, ErrOrderNotFound
	}

	if userRole != constants.RoleSuperAdmin && userRole != constants.RoleStoreManager {
		if order.UserID != userID {
			return nil, ErrUnauthorizedOrder
		}
	}

	if !order.CanBeCancelled() {
		return nil, ErrOrderCannotCancel
	}

	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Restore Stock for all items
		for _, item := range order.OrderItems {
			_ = tx.Model(&models.ProductUnit{}).
				Where("id = ?", item.ProductUnitID).
				Update("stock_quantity", gorm.Expr("stock_quantity + ?", item.Quantity))

			var inv models.Inventory
			if err := tx.Where("product_unit_id = ?", item.ProductUnitID).First(&inv).Error; err == nil {
				qtyBefore := inv.CurrentStock
				qtyAfter := qtyBefore + item.Quantity
				inv.CurrentStock = qtyAfter
				_ = tx.Save(&inv)

				_ = s.stockMoveRepo.Create(ctx, tx, &models.StockMovement{
					InventoryID:     inv.ID,
					ProductUnitID:   item.ProductUnitID,
					MovementType:    "STOCK_IN",
					ReferenceNumber: fmt.Sprintf("CANCEL-%s", order.OrderNumber),
					QuantityBefore:  qtyBefore,
					QuantityChanged: item.Quantity,
					QuantityAfter:   qtyAfter,
					Notes:           fmt.Sprintf("Pengembalian stok pembatalan order %s", order.OrderNumber),
				})
			}
		}

		order.Status = constants.OrderStatusCancelled
		if err := s.orderRepo.Update(ctx, tx, order); err != nil {
			return err
		}

		if order.Invoice != nil {
			_ = s.invoiceRepo.UpdateStatus(ctx, tx, order.Invoice.ID, "CANCELLED")
		}

		return s.orderRepo.AddHistory(ctx, tx, &models.OrderHistory{
			OrderID: order.ID,
			Status:  constants.OrderStatusCancelled,
			Notes:   fmt.Sprintf("Pesanan dibatalkan: %s", req.Reason),
		})
	})

	if err != nil {
		return nil, fmt.Errorf("failed to cancel order: %w", err)
	}

	return s.GetOrderByID(ctx, userID, userRole, id)
}

func (s *orderService) GetOrderHistory(ctx context.Context, userID uuid.UUID, userRole constants.UserRole, id uuid.UUID) ([]dto.OrderHistoryResponseDTO, error) {
	order, err := s.GetOrderByID(ctx, userID, userRole, id)
	if err != nil {
		return nil, err
	}
	return order.Histories, nil
}

func (s *orderService) GetInvoiceByID(ctx context.Context, userID uuid.UUID, userRole constants.UserRole, invoiceID uuid.UUID) (*dto.InvoiceResponseDTO, error) {
	inv, err := s.invoiceRepo.FindByID(ctx, invoiceID)
	if err != nil || inv == nil {
		return nil, errors.New("invoice not found")
	}

	if userRole != constants.RoleSuperAdmin && userRole != constants.RoleStoreManager {
		if inv.UserID != userID {
			return nil, ErrUnauthorizedOrder
		}
	}

	return &dto.InvoiceResponseDTO{
		ID:            inv.ID,
		InvoiceNumber: inv.InvoiceNumber,
		OrderID:       inv.OrderID,
		TotalAmount:   inv.TotalAmount,
		Status:        inv.Status,
		DueDate:       inv.DueDate,
		IssuedAt:      inv.IssuedAt,
		PaidAt:        inv.PaidAt,
	}, nil
}

func (s *orderService) UpdateOrderStatus(ctx context.Context, userID uuid.UUID, userRole constants.UserRole, id uuid.UUID, req dto.UpdateOrderStatusDTO) (*dto.OrderResponseDTO, error) {
	order, err := s.orderRepo.FindByID(ctx, id)
	if err != nil || order == nil {
		return nil, ErrOrderNotFound
	}

	if userRole != constants.RoleSuperAdmin && userRole != constants.RoleStoreManager && userRole != constants.RoleWarehouseStaff {
		return nil, ErrUnauthorizedOrder
	}

	notes := req.Notes
	if notes == "" {
		notes = fmt.Sprintf("Status pesanan diperbarui menjadi %s oleh Admin", req.Status)
	}

	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		order.Status = req.Status
		if err := s.orderRepo.Update(ctx, tx, order); err != nil {
			return err
		}

		if req.Status == constants.OrderStatusPaid && order.Invoice != nil {
			now := time.Now()
			order.Invoice.Status = "PAID"
			order.Invoice.PaidAt = &now
			_ = tx.Save(order.Invoice)
		}

		return s.orderRepo.AddHistory(ctx, tx, &models.OrderHistory{
			OrderID: order.ID,
			Status:  req.Status,
			Notes:   notes,
		})
	})

	if err != nil {
		return nil, fmt.Errorf("failed to update order status: %w", err)
	}

	return s.GetOrderByID(ctx, userID, userRole, id)
}


func mapOrderToDTO(o *models.Order) dto.OrderResponseDTO {
	res := dto.OrderResponseDTO{
		ID:            o.ID,
		OrderNumber:   o.OrderNumber,
		UserID:        o.UserID,
		Status:        o.Status,
		PaymentMethod: o.PaymentMethod,
		TotalAmount:   o.TotalAmount,
		ShippingFee:   o.ShippingFee,
		GrandTotal:    o.GrandTotal,
		DueDate:       o.DueDate,
		OrderDate:     o.OrderDate,
	}

	if o.OrderAddress != nil {
		res.Address = &dto.OrderAddressResponseDTO{
			RecipientName: o.OrderAddress.RecipientName,
			Phone:         o.OrderAddress.Phone,
			FullAddress:   o.OrderAddress.FullAddress,
			City:          o.OrderAddress.City,
			PostalCode:    o.OrderAddress.PostalCode,
		}
	}

	if o.Invoice != nil {
		res.Invoice = &dto.InvoiceResponseDTO{
			ID:            o.Invoice.ID,
			InvoiceNumber: o.Invoice.InvoiceNumber,
			OrderID:       o.Invoice.OrderID,
			TotalAmount:   o.Invoice.TotalAmount,
			Status:        o.Invoice.Status,
			DueDate:       o.Invoice.DueDate,
			IssuedAt:      o.Invoice.IssuedAt,
			PaidAt:        o.Invoice.PaidAt,
		}
	}

	items := make([]dto.OrderItemResponseDTO, len(o.OrderItems))
	for i, item := range o.OrderItems {
		prodName := "Product"
		sku := ""
		if item.Product != nil {
			prodName = item.Product.Name
			sku = item.Product.SKU
		}
		unitName := "Unit"
		if item.ProductUnit != nil {
			unitName = item.ProductUnit.UnitName
		}

		items[i] = dto.OrderItemResponseDTO{
			ID:            item.ID,
			ProductID:     item.ProductID,
			ProductName:   prodName,
			SKU:           sku,
			ProductUnitID: item.ProductUnitID,
			UnitName:      unitName,
			Quantity:      item.Quantity,
			UnitPrice:     item.UnitPrice,
			Subtotal:      item.Subtotal,
			Notes:         item.Notes,
		}
	}
	res.Items = items

	histories := make([]dto.OrderHistoryResponseDTO, len(o.Histories))
	for i, h := range o.Histories {
		histories[i] = dto.OrderHistoryResponseDTO{
			ID:        h.ID,
			Status:    h.Status,
			Notes:     h.Notes,
			CreatedAt: h.CreatedAt,
		}
	}
	res.Histories = histories

	return res
}

func generateOrderNumber() string {
	now := time.Now()
	return fmt.Sprintf("ORD-%s-%04d", now.Format("20060102"), now.UnixNano()%10000)
}
