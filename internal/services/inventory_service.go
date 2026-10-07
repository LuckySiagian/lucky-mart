package services

import (
	"context"
	"errors"
	"fmt"

	"backend-sembako/internal/dto"
	"backend-sembako/internal/models"
	"backend-sembako/internal/repositories"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrInventoryNotFound = errors.New("inventory record not found")
	ErrNegativeStock     = errors.New("stock adjustment results in negative stock quantity which is strictly forbidden")
)

type InventoryService interface {
	AdjustStock(ctx context.Context, req dto.StockAdjustmentDTO) (*dto.StockMovementResponseDTO, error)
	GetLowStockAlerts(ctx context.Context) ([]dto.InventoryStatusResponseDTO, error)
	GetStockHistory(ctx context.Context, inventoryID uuid.UUID, req dto.PaginationRequest) (dto.PaginatedResponse, error)
}

type inventoryService struct {
	db            *gorm.DB
	inventoryRepo repositories.InventoryRepository
	stockMoveRepo repositories.StockMovementRepository
}

func NewInventoryService(
	db *gorm.DB,
	inventoryRepo repositories.InventoryRepository,
	stockMoveRepo repositories.StockMovementRepository,
) InventoryService {
	return &inventoryService{
		db:            db,
		inventoryRepo: inventoryRepo,
		stockMoveRepo: stockMoveRepo,
	}
}

func (s *inventoryService) AdjustStock(ctx context.Context, req dto.StockAdjustmentDTO) (*dto.StockMovementResponseDTO, error) {
	inv, err := s.inventoryRepo.FindByID(ctx, req.InventoryID)
	if err != nil || inv == nil {
		return nil, ErrInventoryNotFound
	}

	var movement *models.StockMovement

	// Execute DB Transaction for Stock Change + Movement Logging
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		qtyBefore := inv.CurrentStock
		qtyChanged := req.QuantityChanged

		if req.AdjustmentType == "STOCK_OUT" && qtyChanged > 0 {
			qtyChanged = -qtyChanged
		}

		qtyAfter := qtyBefore + qtyChanged

		// ENFORCE BUSINESS RULE: Stock cannot be negative
		if qtyAfter < 0 {
			return ErrNegativeStock
		}

		// Update Inventory current stock
		inv.CurrentStock = qtyAfter
		if err := s.inventoryRepo.UpdateStock(ctx, tx, inv); err != nil {
			return err
		}

		// Update ProductUnit stock quantity
		if err := tx.Model(&models.ProductUnit{}).
			Where("id = ?", req.ProductUnitID).
			Update("stock_quantity", qtyAfter).Error; err != nil {
			return err
		}

		// Create mandatory StockMovement Record
		movement = &models.StockMovement{
			InventoryID:     inv.ID,
			ProductUnitID:   req.ProductUnitID,
			MovementType:    req.AdjustmentType,
			ReferenceNumber: req.ReferenceNumber,
			QuantityBefore:  qtyBefore,
			QuantityChanged: qtyChanged,
			QuantityAfter:   qtyAfter,
			Notes:           req.Notes,
		}

		return s.stockMoveRepo.Create(ctx, tx, movement)
	})

	if err != nil {
		return nil, fmt.Errorf("failed to process stock adjustment: %w", err)
	}

	return &dto.StockMovementResponseDTO{
		ID:              movement.ID,
		InventoryID:     movement.InventoryID,
		MovementType:    movement.MovementType,
		ReferenceNumber: movement.ReferenceNumber,
		QuantityBefore:  movement.QuantityBefore,
		QuantityChanged: movement.QuantityChanged,
		QuantityAfter:   movement.QuantityAfter,
		Notes:           movement.Notes,
		CreatedAt:       movement.CreatedAt,
	}, nil
}

func (s *inventoryService) GetLowStockAlerts(ctx context.Context) ([]dto.InventoryStatusResponseDTO, error) {
	inventories, err := s.inventoryRepo.FindLowStock(ctx)
	if err != nil {
		return nil, err
	}

	dtos := make([]dto.InventoryStatusResponseDTO, len(inventories))
	for i, inv := range inventories {
		whName := "Main Warehouse"
		if inv.Warehouse != nil {
			whName = inv.Warehouse.Name
		}
		prodName := "Product"
		sku := ""
		if inv.Product != nil {
			prodName = inv.Product.Name
			sku = inv.Product.SKU
		}
		unitName := "Unit"
		if inv.ProductUnit != nil {
			unitName = inv.ProductUnit.UnitName
		}

		dtos[i] = dto.InventoryStatusResponseDTO{
			ID:            inv.ID,
			WarehouseName: whName,
			ProductName:   prodName,
			SKU:           sku,
			UnitName:      unitName,
			CurrentStock:  inv.CurrentStock,
			ReservedStock: inv.ReservedStock,
			MinStockAlert: inv.MinStockAlert,
			IsLowStock:    inv.IsLowStock(),
		}
	}

	return dtos, nil
}

func (s *inventoryService) GetStockHistory(ctx context.Context, inventoryID uuid.UUID, req dto.PaginationRequest) (dto.PaginatedResponse, error) {
	movements, total, err := s.stockMoveRepo.FindByInventoryID(ctx, inventoryID, req)
	if err != nil {
		return dto.PaginatedResponse{}, err
	}

	dtos := make([]dto.StockMovementResponseDTO, len(movements))
	for i, m := range movements {
		dtos[i] = dto.StockMovementResponseDTO{
			ID:              m.ID,
			InventoryID:     m.InventoryID,
			MovementType:    m.MovementType,
			ReferenceNumber: m.ReferenceNumber,
			QuantityBefore:  m.QuantityBefore,
			QuantityChanged: m.QuantityChanged,
			QuantityAfter:   m.QuantityAfter,
			Notes:           m.Notes,
			CreatedAt:       m.CreatedAt,
		}
	}

	req.SetDefaults()
	return dto.NewPaginatedResponse(dtos, req.Page, req.Limit, total), nil
}
