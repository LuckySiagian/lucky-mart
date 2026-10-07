package dto

import (
	"time"

	"github.com/google/uuid"
)

type CalculateShippingDTO struct {
	City             string    `json:"city" binding:"required"`
	ShippingMethodID uuid.UUID `json:"shipping_method_id" binding:"required"`
	TotalWeightKg    float64   `json:"total_weight_kg" binding:"required,gt=0"`
}

type ShippingRateResponseDTO struct {
	ShippingMethodID uuid.UUID `json:"shipping_method_id"`
	CourierName      string    `json:"courier_name"`
	MethodName       string    `json:"method_name"`
	EstimatedDays    string    `json:"estimated_days"`
	CalculatedFee    float64   `json:"calculated_fee"`
}

type CreateShipmentDTO struct {
	OrderID          uuid.UUID `json:"order_id" binding:"required"`
	CourierID        uuid.UUID `json:"courier_id" binding:"required"`
	ShippingMethodID uuid.UUID `json:"shipping_method_id" binding:"required"`
	ShippingFee      float64   `json:"shipping_fee" binding:"required,gte=0"`
}

type UpdateShipmentStatusDTO struct {
	Status   string `json:"status" binding:"required,oneof=PENDING READY_TO_SHIP PICKED_UP IN_TRANSIT OUT_FOR_DELIVERY DELIVERED FAILED RETURNED"`
	Location string `json:"location" binding:"required"`
	Notes    string `json:"notes" binding:"omitempty"`
}

type ShipmentTrackingResponseDTO struct {
	ID        uuid.UUID `json:"id"`
	Status    string    `json:"status"`
	Location  string    `json:"location"`
	Notes     string    `json:"notes,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type ShipmentResponseDTO struct {
	ID             uuid.UUID                     `json:"id"`
	TrackingNumber string                        `json:"tracking_number"`
	OrderID        uuid.UUID                     `json:"order_id"`
	CourierName    string                        `json:"courier_name"`
	MethodName     string                        `json:"method_name"`
	Status         string                        `json:"status"`
	ShippingFee    float64                       `json:"shipping_fee"`
	ShippedAt      *time.Time                    `json:"shipped_at,omitempty"`
	DeliveredAt    *time.Time                    `json:"delivered_at,omitempty"`
	Trackings      []ShipmentTrackingResponseDTO `json:"trackings,omitempty"`
	CreatedAt      time.Time                     `json:"created_at"`
}
