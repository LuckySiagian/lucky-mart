package services

import (
	"context"
	"errors"
	"strings"
	"time"

	"backend-sembako/internal/constants"
	"backend-sembako/internal/models"
	"backend-sembako/internal/repositories"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Shared order-lifecycle steps used by the order, payment and shipping
// services so that, for example, a delivered shipment and an admin marking
// the order completed have exactly the same effect.

func setOrderStatusTx(tx *gorm.DB, order *models.Order, status constants.OrderStatus, note string) error {
	order.Status = status
	if err := tx.Model(&models.Order{}).Where("id = ?", order.ID).Update("status", status).Error; err != nil {
		return err
	}
	return tx.Create(&models.OrderHistory{OrderID: order.ID, Status: status, Notes: note}).Error
}

// completeOrderTx marks an order completed. Cash-on-delivery payments are
// collected on delivery, so they are marked paid here as well.
func completeOrderTx(tx *gorm.DB, order *models.Order, note string) error {
	if err := setOrderStatusTx(tx, order, constants.OrderStatusCompleted, note); err != nil {
		return err
	}

	now := time.Now()
	var payments []models.Payment
	if err := tx.Preload("PaymentMethod").Where("order_id = ? AND status IN ?", order.ID,
		[]string{constants.PaymentStatusPending, constants.PaymentStatusWaitingConfirmation}).Find(&payments).Error; err != nil {
		return err
	}
	for _, p := range payments {
		if p.PaymentMethod == nil || p.PaymentMethod.Type != constants.PaymentTypeCOD {
			continue
		}
		if err := tx.Model(&models.Payment{}).Where("id = ?", p.ID).
			Updates(map[string]interface{}{"status": constants.PaymentStatusPaid, "paid_at": now}).Error; err != nil {
			return err
		}
		if err := tx.Create(&models.PaymentLog{
			PaymentID: p.ID, Action: "COD_COLLECTED", OldStatus: p.Status, NewStatus: constants.PaymentStatusPaid,
			Notes: "Pembayaran COD diterima saat pesanan selesai",
		}).Error; err != nil {
			return err
		}
		if err := markInvoicePaidTx(tx, p.InvoiceID, now); err != nil {
			return err
		}
	}

	return tx.Model(&models.Shipment{}).
		Where("order_id = ? AND status NOT IN ?", order.ID, []string{"DELIVERED", "RETURNED"}).
		Updates(map[string]interface{}{"status": "DELIVERED", "delivered_at": now}).Error
}

func markInvoicePaidTx(tx *gorm.DB, invoiceID uuid.UUID, at time.Time) error {
	return tx.Model(&models.Invoice{}).Where("id = ?", invoiceID).
		Updates(map[string]interface{}{"status": constants.InvoiceStatusPaid, "paid_at": at}).Error
}

// ensureShipmentTx returns the order's shipment, creating one with the
// store's default courier when the order leaves for delivery without one.
func ensureShipmentTx(tx *gorm.DB, order *models.Order, status string, location, note string) (*models.Shipment, error) {
	var shipment models.Shipment
	err := tx.Where("order_id = ?", order.ID).First(&shipment).Error
	if err == nil {
		return &shipment, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	var method models.ShippingMethod
	if err := tx.Where("is_active = true").Order("base_rate ASC").First(&method).Error; err != nil {
		return nil, errors.New("belum ada metode pengiriman aktif; tambahkan dulu di pengaturan pengiriman")
	}

	now := time.Now()
	shipment = models.Shipment{
		TrackingNumber:   newReference("RESI"),
		OrderID:          order.ID,
		CourierID:        method.CourierID,
		ShippingMethodID: method.ID,
		Status:           status,
		ShippingFee:      order.ShippingFee,
	}
	if status == "OUT_FOR_DELIVERY" || status == "IN_TRANSIT" || status == "PICKED_UP" {
		shipment.ShippedAt = &now
	}
	if err := tx.Create(&shipment).Error; err != nil {
		return nil, err
	}
	if err := tx.Create(&models.ShipmentTracking{
		ShipmentID: shipment.ID, Status: status, Location: location, Notes: note,
	}).Error; err != nil {
		return nil, err
	}
	return &shipment, nil
}

// quoteShippingFee returns the shipping method to use and its fee. A
// city-specific flat rate overrides the method's base rate. With no method
// given, the cheapest active method is used.
func quoteShippingFee(ctx context.Context, methodRepo repositories.ShippingMethodRepository, rateRepo repositories.ShippingRateRepository, methodID *uuid.UUID, city string) (*models.ShippingMethod, float64, error) {
	var method *models.ShippingMethod
	if methodID != nil && *methodID != uuid.Nil {
		m, err := methodRepo.FindByID(ctx, *methodID)
		if err != nil {
			return nil, 0, err
		}
		if m == nil || !m.IsActive {
			return nil, 0, errors.New("metode pengiriman tidak ditemukan atau tidak aktif")
		}
		method = m
	} else {
		methods, err := methodRepo.FindAll(ctx)
		if err != nil {
			return nil, 0, err
		}
		if len(methods) == 0 {
			return nil, 0, errors.New("belum ada metode pengiriman aktif")
		}
		method = &methods[0]
		for i := range methods {
			if methods[i].BaseRate < method.BaseRate {
				method = &methods[i]
			}
		}
	}

	fee := method.BaseRate
	if city = strings.TrimSpace(city); city != "" {
		if rate, err := rateRepo.FindByMethodAndCity(ctx, method.ID, city); err == nil && rate != nil {
			fee = rate.FlatRate
		}
	}
	return method, fee, nil
}
