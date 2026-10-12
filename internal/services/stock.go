package services

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"backend-sembako/internal/models"
	"backend-sembako/pkg/eventbus"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var ErrProductUnitNotFound = errors.New("satuan produk tidak ditemukan")

// applyStockChange changes the sellable stock of one product unit inside tx.
// product_units.stock_quantity is the source of truth used by the cart and
// checkout; the inventories row mirrors it and every change is recorded in
// stock_movements. The unit row is locked so concurrent orders cannot oversell.
func applyStockChange(tx *gorm.DB, unitID uuid.UUID, delta int, movementType, reference, notes string) (*models.StockMovement, error) {
	var unit models.ProductUnit
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", unitID).First(&unit).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrProductUnitNotFound
		}
		return nil, err
	}

	before := unit.StockQuantity
	after := before + delta
	if after < 0 {
		return nil, &InsufficientStockError{Available: before}
	}

	if err := tx.Model(&models.ProductUnit{}).Where("id = ?", unit.ID).Update("stock_quantity", after).Error; err != nil {
		return nil, err
	}

	inv, err := ensureInventory(tx, &unit)
	if err != nil {
		return nil, err
	}
	if err := tx.Model(&models.Inventory{}).Where("id = ?", inv.ID).Update("current_stock", after).Error; err != nil {
		return nil, err
	}

	movement := &models.StockMovement{
		InventoryID:     inv.ID,
		ProductUnitID:   unit.ID,
		MovementType:    movementType,
		ReferenceNumber: reference,
		QuantityBefore:  before,
		QuantityChanged: delta,
		QuantityAfter:   after,
		Notes:           notes,
	}
	if err := tx.Create(movement).Error; err != nil {
		return nil, err
	}
	return movement, nil
}

// ensureInventory returns the inventory row for a product unit, creating it
// (in the default warehouse) when the product was created without one.
func ensureInventory(tx *gorm.DB, unit *models.ProductUnit) (*models.Inventory, error) {
	var inv models.Inventory
	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("product_unit_id = ?", unit.ID).First(&inv).Error
	if err == nil {
		return &inv, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	warehouse, err := defaultWarehouse(tx)
	if err != nil {
		return nil, err
	}
	inv = models.Inventory{
		WarehouseID:   warehouse.ID,
		ProductID:     unit.ProductID,
		ProductUnitID: unit.ID,
		CurrentStock:  unit.StockQuantity,
		MinStockAlert: 10,
	}
	if err := tx.Create(&inv).Error; err != nil {
		return nil, err
	}
	return &inv, nil
}

func defaultWarehouse(tx *gorm.DB) (*models.Warehouse, error) {
	var wh models.Warehouse
	err := tx.Order("created_at ASC").First(&wh).Error
	if err == nil {
		return &wh, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	wh = models.Warehouse{Code: "WH-MAIN-01", Name: "Gudang Utama", Address: "-", IsActive: true}
	if err := tx.Create(&wh).Error; err != nil {
		return nil, err
	}
	return &wh, nil
}

// newReference builds a human-readable, practically unique document number,
// e.g. ORD-20261008-3F9A1C.
func newReference(prefix string) string {
	suffix := strings.ToUpper(strings.ReplaceAll(uuid.NewString(), "-", "")[:6])
	return fmt.Sprintf("%s-%s-%s", prefix, time.Now().Format("20060102"), suffix)
}

// publishEvent sends a domain event; subscribers turn it into an audit log
// entry and, when payload has user_id, an in-app notification. A nil bus is
// allowed so services stay usable in tests.
func publishEvent(bus eventbus.EventBus, name, aggregate string, id uuid.UUID, payload map[string]interface{}) {
	if bus == nil {
		return
	}
	_ = bus.Publish(context.Background(), eventbus.DomainEvent{
		ID:          uuid.New(),
		Name:        name,
		Aggregate:   aggregate,
		AggregateID: id,
		Payload:     payload,
		OccurredAt:  time.Now(),
	})
}

func formatRupiah(amount float64) string {
	s := fmt.Sprintf("%.0f", amount)
	var out []byte
	for i := range s {
		if i > 0 && (len(s)-i)%3 == 0 {
			out = append(out, '.')
		}
		out = append(out, s[i])
	}
	return "Rp" + string(out)
}

// InsufficientStockError reports the stock actually available at the moment
// the row was locked. errors.Is(err, ErrNegativeStock) still matches.
type InsufficientStockError struct {
	Available int
}

func (e *InsufficientStockError) Error() string { return ErrNegativeStock.Error() }
func (e *InsufficientStockError) Unwrap() error { return ErrNegativeStock }
