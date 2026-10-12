package dto

import (
	"time"

	"github.com/google/uuid"
)

type CalculateShippingDTO struct {
	City             string     `json:"city" binding:"omitempty"`
	ShippingMethodID *uuid.UUID `json:"shipping_method_id" binding:"omitempty"`
	TotalWeightKg    float64    `json:"total_weight_kg" binding:"omitempty,gte=0"`
}

type ShippingRateResponseDTO struct {
	ShippingMethodID uuid.UUID `json:"shipping_method_id"`
	CourierName      string    `json:"courier_name"`
	MethodName       string    `json:"method_name"`
	EstimatedDays    string    `json:"estimated_days"`
	CalculatedFee    float64   `json:"calculated_fee"`
}

type ShippingMethodDTO struct {
	ID            uuid.UUID `json:"id"`
	Code          string    `json:"code"`
	Name          string    `json:"name"`
	Description   string    `json:"description,omitempty"`
	EstimatedDays string    `json:"estimated_days"`
	BaseRate      float64   `json:"base_rate"`
	CourierID     uuid.UUID `json:"courier_id"`
	CourierName   string    `json:"courier_name"`
}

type CreateShipmentDTO struct {
	OrderID          uuid.UUID  `json:"order_id" binding:"required"`
	CourierID        *uuid.UUID `json:"courier_id" binding:"omitempty"`
	ShippingMethodID *uuid.UUID `json:"shipping_method_id" binding:"omitempty"`
	ShippingFee      *float64   `json:"shipping_fee" binding:"omitempty,gte=0"`
	Notes            string     `json:"notes" binding:"omitempty,max=255"`
}

type UpdateShipmentStatusDTO struct {
	Status   string `json:"status" binding:"required,oneof=PENDING READY_TO_SHIP PICKED_UP IN_TRANSIT OUT_FOR_DELIVERY DELIVERED FAILED RETURNED"`
	Location string `json:"location" binding:"omitempty,max=255"`
	Notes    string `json:"notes" binding:"omitempty,max=500"`
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
	OrderNumber    string                        `json:"order_number,omitempty"`
	OrderStatus    string                        `json:"order_status,omitempty"`
	CustomerName   string                        `json:"customer_name,omitempty"`
	RecipientName  string                        `json:"recipient_name,omitempty"`
	RecipientPhone string                        `json:"recipient_phone,omitempty"`
	Address        string                        `json:"address,omitempty"`
	CourierName    string                        `json:"courier_name"`
	MethodName     string                        `json:"method_name"`
	Status         string                        `json:"status"`
	ShippingFee    float64                       `json:"shipping_fee"`
	ShippedAt      *time.Time                    `json:"shipped_at,omitempty"`
	DeliveredAt    *time.Time                    `json:"delivered_at,omitempty"`
	Trackings      []ShipmentTrackingResponseDTO `json:"trackings"`
	CreatedAt      time.Time                     `json:"created_at"`
}
