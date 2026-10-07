package repositories

import (
	"context"

	"backend-sembako/internal/dto"
	"backend-sembako/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	FindByPhoneNumber(ctx context.Context, phone string) (*models.User, error)
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	FindByPasswordResetToken(ctx context.Context, token string) (*models.User, error)
	FindByEmailVerificationToken(ctx context.Context, token string) (*models.User, error)
	Create(ctx context.Context, user *models.User) error
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id uuid.UUID) error
	FindAll(ctx context.Context, req dto.PaginationRequest) ([]models.User, int64, error)
}

type B2BProfileRepository interface {
	FindByUserID(ctx context.Context, userID uuid.UUID) (*models.B2BProfile, error)
	Create(ctx context.Context, profile *models.B2BProfile) error
	Update(ctx context.Context, profile *models.B2BProfile) error
}

type RefreshTokenRepository interface {
	Create(ctx context.Context, token *models.RefreshToken) error
	FindByTokenHash(ctx context.Context, hash string) (*models.RefreshToken, error)
	RevokeByUserID(ctx context.Context, userID uuid.UUID) error
	RevokeToken(ctx context.Context, id uuid.UUID) error
}

type LoginHistoryRepository interface {
	Create(ctx context.Context, history *models.LoginHistory) error
	FindByUserID(ctx context.Context, userID uuid.UUID, limit int) ([]models.LoginHistory, error)
}

type PermissionRepository interface {
	GetRolePermissions(ctx context.Context, role string) ([]string, error)
	HasPermission(ctx context.Context, role string, permissionCode string) (bool, error)
}

type CategoryRepository interface {
	FindAll(ctx context.Context) ([]models.Category, error)
	FindByID(ctx context.Context, id uuid.UUID) (*models.Category, error)
	Create(ctx context.Context, category *models.Category) error
	Update(ctx context.Context, category *models.Category) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type ProductRepository interface {
	FindAll(ctx context.Context, req dto.PaginationRequest, categoryID *uuid.UUID) ([]models.Product, int64, error)
	FindByID(ctx context.Context, id uuid.UUID) (*models.Product, error)
	FindBySKU(ctx context.Context, sku string) (*models.Product, error)
	Create(ctx context.Context, tx *gorm.DB, product *models.Product) error
	Update(ctx context.Context, tx *gorm.DB, product *models.Product) error
	Delete(ctx context.Context, id uuid.UUID) error
	AddImage(ctx context.Context, image *models.ProductImage) error
}

type InventoryRepository interface {
	FindByProductUnitID(ctx context.Context, unitID uuid.UUID) (*models.Inventory, error)
	FindByID(ctx context.Context, id uuid.UUID) (*models.Inventory, error)
	FindLowStock(ctx context.Context) ([]models.Inventory, error)
	Create(ctx context.Context, tx *gorm.DB, inv *models.Inventory) error
	UpdateStock(ctx context.Context, tx *gorm.DB, inv *models.Inventory) error
}

type StockMovementRepository interface {
	Create(ctx context.Context, tx *gorm.DB, movement *models.StockMovement) error
	FindByInventoryID(ctx context.Context, inventoryID uuid.UUID, req dto.PaginationRequest) ([]models.StockMovement, int64, error)
}

type CartRepository interface {
	FindByUserID(ctx context.Context, userID uuid.UUID) (*models.Cart, error)
	Create(ctx context.Context, cart *models.Cart) error
	Update(ctx context.Context, cart *models.Cart) error
	Clear(ctx context.Context, cartID uuid.UUID) error
}

type CartItemRepository interface {
	FindByID(ctx context.Context, itemID uuid.UUID) (*models.CartItem, error)
	FindByCartAndUnit(ctx context.Context, cartID, unitID uuid.UUID) (*models.CartItem, error)
	Create(ctx context.Context, item *models.CartItem) error
	Update(ctx context.Context, item *models.CartItem) error
	Delete(ctx context.Context, itemID uuid.UUID) error
}

type OrderRepository interface {
	Create(ctx context.Context, tx *gorm.DB, order *models.Order) error
	Update(ctx context.Context, tx *gorm.DB, order *models.Order) error
	FindByID(ctx context.Context, id uuid.UUID) (*models.Order, error)
	FindByUserID(ctx context.Context, userID uuid.UUID, req dto.PaginationRequest) ([]models.Order, int64, error)
	FindAll(ctx context.Context, req dto.PaginationRequest) ([]models.Order, int64, error)
	AddHistory(ctx context.Context, tx *gorm.DB, history *models.OrderHistory) error
}

type InvoiceRepository interface {
	Create(ctx context.Context, tx *gorm.DB, invoice *models.Invoice) error
	FindByID(ctx context.Context, id uuid.UUID) (*models.Invoice, error)
	FindByOrderID(ctx context.Context, orderID uuid.UUID) (*models.Invoice, error)
	UpdateStatus(ctx context.Context, tx *gorm.DB, invoiceID uuid.UUID, status string) error
}

type AddressRepository interface {
	Create(ctx context.Context, addr *models.Address) error
	Update(ctx context.Context, addr *models.Address) error
	Delete(ctx context.Context, id uuid.UUID) error
	FindByID(ctx context.Context, id uuid.UUID) (*models.Address, error)
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]models.Address, error)
	UnsetDefault(ctx context.Context, userID uuid.UUID) error
}

type CourierRepository interface {
	FindAll(ctx context.Context) ([]models.Courier, error)
	FindByID(ctx context.Context, id uuid.UUID) (*models.Courier, error)
}

type ShippingMethodRepository interface {
	FindAll(ctx context.Context) ([]models.ShippingMethod, error)
	FindByID(ctx context.Context, id uuid.UUID) (*models.ShippingMethod, error)
}

type ShippingRateRepository interface {
	FindByMethodAndCity(ctx context.Context, methodID uuid.UUID, city string) (*models.ShippingRate, error)
}

type ShipmentRepository interface {
	Create(ctx context.Context, tx *gorm.DB, shipment *models.Shipment) error
	Update(ctx context.Context, tx *gorm.DB, shipment *models.Shipment) error
	FindByID(ctx context.Context, id uuid.UUID) (*models.Shipment, error)
	FindByOrderID(ctx context.Context, orderID uuid.UUID) (*models.Shipment, error)
	AddTracking(ctx context.Context, tx *gorm.DB, tracking *models.ShipmentTracking) error
}

type PaymentRepository interface {
	Create(ctx context.Context, tx *gorm.DB, payment *models.Payment) error
	Update(ctx context.Context, tx *gorm.DB, payment *models.Payment) error
	FindByID(ctx context.Context, id uuid.UUID) (*models.Payment, error)
	FindByInvoiceID(ctx context.Context, invoiceID uuid.UUID) (*models.Payment, error)
	FindByUserID(ctx context.Context, userID uuid.UUID, req dto.PaginationRequest) ([]models.Payment, int64, error)
	FindAll(ctx context.Context, req dto.PaginationRequest) ([]models.Payment, int64, error)
	FindWaitingConfirmation(ctx context.Context, req dto.PaginationRequest) ([]models.Payment, int64, error)
}

type PaymentMethodRepository interface {
	FindAll(ctx context.Context) ([]models.PaymentMethod, error)
	FindByID(ctx context.Context, id uuid.UUID) (*models.PaymentMethod, error)
	FindByCode(ctx context.Context, code string) (*models.PaymentMethod, error)
	Create(ctx context.Context, method *models.PaymentMethod) error
	Update(ctx context.Context, method *models.PaymentMethod) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type PaymentTransactionRepository interface {
	Create(ctx context.Context, tx *gorm.DB, transaction *models.PaymentTransaction) error
}

type PaymentLogRepository interface {
	Create(ctx context.Context, tx *gorm.DB, log *models.PaymentLog) error
}

type CreditLedgerRepository interface {
	FindByUserID(ctx context.Context, userID uuid.UUID) (*models.CreditLedger, error)
	FindByID(ctx context.Context, id uuid.UUID) (*models.CreditLedger, error)
	FindAll(ctx context.Context, req dto.PaginationRequest) ([]models.CreditLedger, int64, error)
	Create(ctx context.Context, tx *gorm.DB, ledger *models.CreditLedger) error
	Update(ctx context.Context, tx *gorm.DB, ledger *models.CreditLedger) error
}

type CreditPaymentRepository interface {
	Create(ctx context.Context, tx *gorm.DB, payment *models.CreditPayment) error
	FindByLedgerID(ctx context.Context, ledgerID uuid.UUID, req dto.PaginationRequest) ([]models.CreditPayment, int64, error)
}

type NotificationRepository interface {
	Create(ctx context.Context, tx *gorm.DB, notification *models.Notification) error
	Update(ctx context.Context, tx *gorm.DB, notification *models.Notification) error
	FindByID(ctx context.Context, id uuid.UUID) (*models.Notification, error)
	FindByUserID(ctx context.Context, userID uuid.UUID, req dto.PaginationRequest) ([]models.Notification, int64, error)
	FindAll(ctx context.Context, req dto.PaginationRequest) ([]models.Notification, int64, error)
	MarkAllAsRead(ctx context.Context, userID uuid.UUID) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type NotificationTemplateRepository interface {
	FindAll(ctx context.Context) ([]models.NotificationTemplate, error)
	FindByID(ctx context.Context, id uuid.UUID) (*models.NotificationTemplate, error)
	FindByEventTypeAndChannel(ctx context.Context, eventType, channel string) (*models.NotificationTemplate, error)
	Create(ctx context.Context, tmpl *models.NotificationTemplate) error
	Update(ctx context.Context, tmpl *models.NotificationTemplate) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type NotificationLogRepository interface {
	Create(ctx context.Context, tx *gorm.DB, log *models.NotificationLog) error
}

type EventLogRepository interface {
	Create(ctx context.Context, tx *gorm.DB, event *models.EventLog) error
	FindByID(ctx context.Context, id uuid.UUID) (*models.EventLog, error)
	FindAll(ctx context.Context, req dto.PaginationRequest) ([]models.EventLog, int64, error)
}
