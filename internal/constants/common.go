package constants

// Application Constants
const (
	DefaultPageSize = 20
	MaxPageSize     = 100
	DefaultSortBy   = "created_at"
	DefaultSortDir  = "DESC"
)

// User Role Enums
type UserRole string

const (
	RoleSuperAdmin     UserRole = "SUPER_ADMIN"
	RoleStoreManager   UserRole = "STORE_MANAGER"
	RoleWarehouseStaff UserRole = "WAREHOUSE_STAFF"
	RoleDriver         UserRole = "DRIVER"
	RoleB2BCustomer    UserRole = "B2B_CUSTOMER"
	RoleB2CCustomer    UserRole = "B2C_CUSTOMER"
)

func (r UserRole) IsValid() bool {
	switch r {
	case RoleSuperAdmin, RoleStoreManager, RoleWarehouseStaff, RoleDriver, RoleB2BCustomer, RoleB2CCustomer:
		return true
	}
	return false
}

// IsStaff reports whether the role belongs to store staff (not a customer).
func (r UserRole) IsStaff() bool {
	switch r {
	case RoleSuperAdmin, RoleStoreManager, RoleWarehouseStaff, RoleDriver:
		return true
	}
	return false
}

// IsManager reports whether the role may manage money, users and catalog.
func (r UserRole) IsManager() bool {
	return r == RoleSuperAdmin || r == RoleStoreManager
}

// B2B Verification Status Enums
type B2BVerificationStatus string

const (
	B2BStatusPending  B2BVerificationStatus = "PENDING"
	B2BStatusApproved B2BVerificationStatus = "APPROVED"
	B2BStatusRejected B2BVerificationStatus = "REJECTED"
)

// Order Status Enums
type OrderStatus string

const (
	OrderStatusCreated           OrderStatus = "CREATED"
	OrderStatusPendingPayment    OrderStatus = "PENDING_PAYMENT"
	OrderStatusPaid              OrderStatus = "PAID"
	OrderStatusConfirmed         OrderStatus = "CONFIRMED"
	OrderStatusProcessingPacking OrderStatus = "PROCESSING_PACKING"
	OrderStatusOutForDelivery    OrderStatus = "OUT_FOR_DELIVERY"
	OrderStatusCompleted         OrderStatus = "COMPLETED"
	OrderStatusCancelled         OrderStatus = "CANCELLED"
)

// orderTransitions lists the statuses staff may move an order to from each
// status, the usual next step first. A small store often skips steps (a
// neighbour picks up straight away), so any later step may be chosen
// directly. Cancelling always goes through the cancel flow so stock,
// payments and kasbon are restored; goods already out can still be
// cancelled when the customer never takes them. PENDING_PAYMENT only
// leaves via payment approval or cancellation.
var orderTransitions = map[OrderStatus][]OrderStatus{
	OrderStatusCreated:           {OrderStatusConfirmed, OrderStatusCancelled},
	OrderStatusPendingPayment:    {OrderStatusCancelled},
	OrderStatusPaid:              {OrderStatusProcessingPacking, OrderStatusOutForDelivery, OrderStatusCompleted, OrderStatusCancelled},
	OrderStatusConfirmed:         {OrderStatusProcessingPacking, OrderStatusOutForDelivery, OrderStatusCompleted, OrderStatusCancelled},
	OrderStatusProcessingPacking: {OrderStatusOutForDelivery, OrderStatusCompleted, OrderStatusCancelled},
	OrderStatusOutForDelivery:    {OrderStatusCompleted, OrderStatusCancelled},
}

// DeliveryType says how the customer receives an order.
const (
	DeliveryTypeDelivery = "DELIVERY" // diantar ke alamat pelanggan
	DeliveryTypePickup   = "PICKUP"   // diambil sendiri di toko, tanpa ongkir
)

func (s OrderStatus) CanTransitionTo(next OrderStatus) bool {
	for _, allowed := range orderTransitions[s] {
		if allowed == next {
			return true
		}
	}
	return false
}

func (s OrderStatus) NextStatuses() []OrderStatus {
	return orderTransitions[s]
}

// Payment Method Codes (match payment_methods.code)
type PaymentMethod string

const (
	PaymentMethodCashOnDelivery     PaymentMethod = "COD"
	PaymentMethodManualBankTransfer PaymentMethod = "MANUAL_BANK_TRANSFER"
	PaymentMethodB2BCreditTempo     PaymentMethod = "B2B_CREDIT_TEMPO"
)

// Payment Method Types (payment_methods.type)
const (
	PaymentTypeCOD            = "COD"
	PaymentTypeManualTransfer = "MANUAL_TRANSFER"
	PaymentTypeB2BCredit      = "B2B_CREDIT"
	PaymentTypeCash           = "CASH"
)

// Payment Status values
const (
	PaymentStatusPending             = "PENDING"
	PaymentStatusWaitingConfirmation = "WAITING_CONFIRMATION"
	PaymentStatusPaid                = "PAID"
	PaymentStatusFailed              = "FAILED"
	PaymentStatusCancelled           = "CANCELLED"
	PaymentStatusRefunded            = "REFUNDED"
	PaymentStatusPartiallyRefunded   = "PARTIALLY_REFUNDED"
)

// Invoice Status values
const (
	InvoiceStatusUnpaid    = "UNPAID"
	InvoiceStatusPaid      = "PAID"
	InvoiceStatusCancelled = "CANCELLED"
)

// Credit Transaction Type Enums
type CreditTransactionType string

const (
	CreditTxUsage      CreditTransactionType = "CREDIT_USAGE"
	CreditTxSettlement CreditTransactionType = "PAYMENT_SETTLEMENT"
	CreditTxAdjustment CreditTransactionType = "LIMIT_ADJUSTMENT"
	CreditTxReversal   CreditTransactionType = "REVERSAL"
)

// Stock Movement Types
const (
	MovementInitialStock = "INITIAL_STOCK"
	MovementStockIn      = "STOCK_IN"
	MovementStockOut     = "STOCK_OUT"
	MovementAdjustment   = "ADJUSTMENT"
)
