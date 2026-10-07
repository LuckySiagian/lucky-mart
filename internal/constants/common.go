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
	OrderStatusCreated          OrderStatus = "CREATED"
	OrderStatusPendingPayment   OrderStatus = "PENDING_PAYMENT"
	OrderStatusPaid             OrderStatus = "PAID"
	OrderStatusConfirmed        OrderStatus = "CONFIRMED"
	OrderStatusProcessingPacking OrderStatus = "PROCESSING_PACKING"
	OrderStatusOutForDelivery   OrderStatus = "OUT_FOR_DELIVERY"
	OrderStatusCompleted        OrderStatus = "COMPLETED"
	OrderStatusCancelled        OrderStatus = "CANCELLED"
)

// Payment Method Enums
type PaymentMethod string

const (
	PaymentMethodCashOnDelivery     PaymentMethod = "CASH_ON_DELIVERY"
	PaymentMethodManualBankTransfer PaymentMethod = "MANUAL_BANK_TRANSFER"
	PaymentMethodB2BCreditTempo     PaymentMethod = "B2B_CREDIT_TEMPO"
)

// Credit Transaction Type Enums
type CreditTransactionType string

const (
	CreditTxUsage      CreditTransactionType = "CREDIT_USAGE"
	CreditTxSettlement CreditTransactionType = "PAYMENT_SETTLEMENT"
	CreditTxAdjustment CreditTransactionType = "LIMIT_ADJUSTMENT"
)
