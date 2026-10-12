package services

import (
	"context"
	"errors"
	"fmt"

	"backend-sembako/internal/constants"
	"backend-sembako/internal/dto"
	"backend-sembako/internal/models"
	"backend-sembako/internal/repositories"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrInventoryNotFound = errors.New("data stok tidak ditemukan")
	ErrNegativeStock     = errors.New("stok tidak boleh kurang dari nol")
)

type InventoryService interface {
	GetInventories(ctx context.Context, req dto.PaginationRequest) (dto.PaginatedResponse, error)
	AdjustStock(ctx context.Context, req dto.StockAdjustmentDTO) (*dto.StockMovementResponseDTO, error)
	GetLowStockAlerts(ctx context.Context) ([]dto.InventoryStatusResponseDTO, error)
	GetStockHistory(ctx context.Context, inventoryID uuid.UUID, req dto.PaginationRequest) (dto.PaginatedResponse, error)
	GetMovements(ctx context.Context, req dto.PaginationRequest, movementType string) (dto.PaginatedResponse, error)
	UpdateMinStockAlert(ctx context.Context, inventoryID uuid.UUID, minStock int) (*dto.InventoryStatusResponseDTO, error)
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

func (s *inventoryService) GetInventories(ctx context.Context, req dto.PaginationRequest) (dto.PaginatedResponse, error) {
	inventories, total, err := s.inventoryRepo.FindAll(ctx, req)
	if err != nil {
		return dto.PaginatedResponse{}, err
	}

	dtos := make([]dto.InventoryStatusResponseDTO, len(inventories))
	for i := range inventories {
		dtos[i] = mapInventoryToDTO(&inventories[i])
	}

	req.SetDefaults()
	return dto.NewPaginatedResponse(dtos, req.Page, req.Limit, total), nil
}

func (s *inventoryService) AdjustStock(ctx context.Context, req dto.StockAdjustmentDTO) (*dto.StockMovementResponseDTO, error) {
	delta := req.QuantityChanged
	switch req.AdjustmentType {
	case constants.MovementStockIn:
		if delta < 0 {
			delta = -delta
		}
	case constants.MovementStockOut:
		if delta > 0 {
			delta = -delta
		}
	}
	if delta == 0 {
		return nil, errors.New("jumlah perubahan stok tidak boleh nol")
	}

	reference := req.ReferenceNumber
	if reference == "" {
		reference = newReference("ADJ")
	}

	var movement *models.StockMovement
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var err error
		movement, err = applyStockChange(tx, req.ProductUnitID, delta, req.AdjustmentType, reference, req.Notes)
		return err
	})
	if err != nil {
		return nil, fmt.Errorf("gagal menyesuaikan stok: %w", err)
	}

	res := mapMovementToDTO(movement)
	return &res, nil
}

func (s *inventoryService) GetLowStockAlerts(ctx context.Context) ([]dto.InventoryStatusResponseDTO, error) {
	inventories, err := s.inventoryRepo.FindLowStock(ctx)
	if err != nil {
		return nil, err
	}

	dtos := make([]dto.InventoryStatusResponseDTO, len(inventories))
	for i := range inventories {
		dtos[i] = mapInventoryToDTO(&inventories[i])
	}
	return dtos, nil
}

func (s *inventoryService) GetStockHistory(ctx context.Context, inventoryID uuid.UUID, req dto.PaginationRequest) (dto.PaginatedResponse, error) {
	movements, total, err := s.stockMoveRepo.FindByInventoryID(ctx, inventoryID, req)
	if err != nil {
		return dto.PaginatedResponse{}, err
	}

	dtos := make([]dto.StockMovementResponseDTO, len(movements))
	for i := range movements {
		dtos[i] = mapMovementToDTO(&movements[i])
	}

	req.SetDefaults()
	return dto.NewPaginatedResponse(dtos, req.Page, req.Limit, total), nil
}

func (s *inventoryService) GetMovements(ctx context.Context, req dto.PaginationRequest, movementType string) (dto.PaginatedResponse, error) {
	movements, total, err := s.stockMoveRepo.FindAll(ctx, req, movementType)
	if err != nil {
		return dto.PaginatedResponse{}, err
	}

	dtos := make([]dto.StockMovementResponseDTO, len(movements))
	for i := range movements {
		dtos[i] = mapMovementToDTO(&movements[i])
	}

	req.SetDefaults()
	return dto.NewPaginatedResponse(dtos, req.Page, req.Limit, total), nil
}

func (s *inventoryService) UpdateMinStockAlert(ctx context.Context, inventoryID uuid.UUID, minStock int) (*dto.InventoryStatusResponseDTO, error) {
	inv, err := s.inventoryRepo.FindByID(ctx, inventoryID)
	if err != nil || inv == nil {
		return nil, ErrInventoryNotFound
	}
	inv.MinStockAlert = minStock
	if err := s.inventoryRepo.UpdateStock(ctx, nil, inv); err != nil {
		return nil, err
	}
	res := mapInventoryToDTO(inv)
	return &res, nil
}

func mapInventoryToDTO(inv *models.Inventory) dto.InventoryStatusResponseDTO {
	res := dto.InventoryStatusResponseDTO{
		ID:            inv.ID,
		ProductID:     inv.ProductID,
		ProductUnitID: inv.ProductUnitID,
		WarehouseName: "Gudang Utama",
		ProductName:   "Produk",
		UnitName:      "Unit",
		CurrentStock:  inv.CurrentStock,
		ReservedStock: inv.ReservedStock,
		MinStockAlert: inv.MinStockAlert,
		IsLowStock:    inv.IsLowStock(),
		UpdatedAt:     inv.UpdatedAt,
	}
	if inv.Warehouse != nil {
		res.WarehouseName = inv.Warehouse.Name
	}
	if inv.Product != nil {
		res.ProductName = inv.Product.Name
		res.SKU = inv.Product.SKU
		res.IsProductActive = inv.Product.IsActive
		if inv.Product.Category != nil {
			res.CategoryName = inv.Product.Category.Name
		}
	}
	if inv.ProductUnit != nil {
		res.UnitName = inv.ProductUnit.UnitName
		res.RetailPrice = inv.ProductUnit.RetailPrice
	}
	return res
}

func mapMovementToDTO(m *models.StockMovement) dto.StockMovementResponseDTO {
	res := dto.StockMovementResponseDTO{
		ID:              m.ID,
		InventoryID:     m.InventoryID,
		ProductUnitID:   m.ProductUnitID,
		MovementType:    m.MovementType,
		ReferenceNumber: m.ReferenceNumber,
		QuantityBefore:  m.QuantityBefore,
		QuantityChanged: m.QuantityChanged,
		QuantityAfter:   m.QuantityAfter,
		Notes:           m.Notes,
		CreatedAt:       m.CreatedAt,
	}
	if m.ProductUnit != nil {
		res.UnitName = m.ProductUnit.UnitName
	}
	if m.Inventory != nil && m.Inventory.Product != nil {
		res.ProductName = m.Inventory.Product.Name
		res.SKU = m.Inventory.Product.SKU
	}
	return res
}
