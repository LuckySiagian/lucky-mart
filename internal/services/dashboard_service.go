package services

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DailySales struct {
	Date    string  `json:"date"`
	Revenue float64 `json:"revenue"`
	Orders  int     `json:"orders"`
}

type TopProduct struct {
	ProductName string  `json:"product_name"`
	Quantity    int     `json:"quantity"`
	Revenue     float64 `json:"revenue"`
}

type RecentOrder struct {
	ID           uuid.UUID `json:"id"`
	OrderNumber  string    `json:"order_number"`
	CustomerName string    `json:"customer_name"`
	GrandTotal   float64   `json:"grand_total"`
	Status       string    `json:"status"`
	OrderDate    time.Time `json:"order_date"`
}

type PendingWork struct {
	OrdersToConfirm  int64 `json:"orders_to_confirm"`
	OrdersToPack     int64 `json:"orders_to_pack"`
	OrdersPacking    int64 `json:"orders_packing"`
	OrdersDelivering int64 `json:"orders_delivering"`
	AwaitingPayment  int64 `json:"awaiting_payment"`
	PaymentsToVerify int64 `json:"payments_to_verify"`
	B2BToVerify      int64 `json:"b2b_to_verify"`
	LowStockItems    int64 `json:"low_stock_items"`
}

type DashboardStats struct {
	TodayRevenue      float64       `json:"today_revenue"`
	TodayOrders       int64         `json:"today_orders"`
	MonthRevenue      float64       `json:"month_revenue"`
	MonthOrders       int64         `json:"month_orders"`
	ActiveProducts    int64         `json:"active_products"`
	TotalCustomers    int64         `json:"total_customers"`
	KasbonOutstanding float64       `json:"kasbon_outstanding"`
	Pending           PendingWork   `json:"pending"`
	SalesLast7Days    []DailySales  `json:"sales_last_7_days"`
	TopProducts       []TopProduct  `json:"top_products"`
	RecentOrders      []RecentOrder `json:"recent_orders"`
}

type DashboardService interface {
	GetStats(ctx context.Context) (*DashboardStats, error)
}

type dashboardService struct {
	db *gorm.DB
}

func NewDashboardService(db *gorm.DB) DashboardService {
	return &dashboardService{db: db}
}

// GetStats compares dates in the store's timezone (Asia/Jakarta) so "today"
// matches the shop's calendar whatever timezone the database server uses.
func (s *dashboardService) GetStats(ctx context.Context) (*DashboardStats, error) {
	db := s.db.WithContext(ctx)
	stats := &DashboardStats{}

	type sumCount struct {
		Total float64
		Count int64
	}
	var today, month sumCount
	if err := db.Raw(`
		SELECT COALESCE(SUM(grand_total), 0) AS total, COUNT(*) AS count FROM orders
		WHERE deleted_at IS NULL AND status <> 'CANCELLED'
		  AND (order_date AT TIME ZONE 'Asia/Jakarta')::date = (NOW() AT TIME ZONE 'Asia/Jakarta')::date`).
		Scan(&today).Error; err != nil {
		return nil, err
	}
	if err := db.Raw(`
		SELECT COALESCE(SUM(grand_total), 0) AS total, COUNT(*) AS count FROM orders
		WHERE deleted_at IS NULL AND status <> 'CANCELLED'
		  AND date_trunc('month', order_date AT TIME ZONE 'Asia/Jakarta') = date_trunc('month', NOW() AT TIME ZONE 'Asia/Jakarta')`).
		Scan(&month).Error; err != nil {
		return nil, err
	}
	stats.TodayRevenue, stats.TodayOrders = today.Total, today.Count
	stats.MonthRevenue, stats.MonthOrders = month.Total, month.Count

	counts := []struct {
		dest  *int64
		query string
	}{
		{&stats.ActiveProducts, `SELECT COUNT(*) FROM products WHERE deleted_at IS NULL AND is_active = true`},
		{&stats.TotalCustomers, `SELECT COUNT(*) FROM users WHERE deleted_at IS NULL AND role IN ('B2C_CUSTOMER','B2B_CUSTOMER')`},
		{&stats.Pending.OrdersToConfirm, `SELECT COUNT(*) FROM orders WHERE deleted_at IS NULL AND status = 'CREATED'`},
		{&stats.Pending.OrdersToPack, `SELECT COUNT(*) FROM orders WHERE deleted_at IS NULL AND status IN ('CONFIRMED','PAID')`},
		{&stats.Pending.OrdersPacking, `SELECT COUNT(*) FROM orders WHERE deleted_at IS NULL AND status = 'PROCESSING_PACKING'`},
		{&stats.Pending.OrdersDelivering, `SELECT COUNT(*) FROM orders WHERE deleted_at IS NULL AND status = 'OUT_FOR_DELIVERY'`},
		{&stats.Pending.AwaitingPayment, `SELECT COUNT(*) FROM orders WHERE deleted_at IS NULL AND status = 'PENDING_PAYMENT'`},
		{&stats.Pending.PaymentsToVerify, `SELECT COUNT(*) FROM payments WHERE deleted_at IS NULL AND status = 'WAITING_CONFIRMATION'`},
		{&stats.Pending.B2BToVerify, `SELECT COUNT(*) FROM b2b_profiles WHERE deleted_at IS NULL AND verification_status = 'PENDING'`},
		{&stats.Pending.LowStockItems, `
			SELECT COUNT(*) FROM inventories i JOIN products p ON p.id = i.product_id
			WHERE i.deleted_at IS NULL AND p.deleted_at IS NULL AND p.is_active = true AND i.current_stock <= i.min_stock_alert`},
	}
	for _, c := range counts {
		if err := db.Raw(c.query).Scan(c.dest).Error; err != nil {
			return nil, err
		}
	}

	if err := db.Raw(`SELECT COALESCE(SUM(outstanding_balance), 0) FROM credit_ledgers WHERE deleted_at IS NULL`).
		Scan(&stats.KasbonOutstanding).Error; err != nil {
		return nil, err
	}

	stats.SalesLast7Days = []DailySales{}
	if err := db.Raw(`
		WITH days AS (
			SELECT generate_series((NOW() AT TIME ZONE 'Asia/Jakarta')::date - 6, (NOW() AT TIME ZONE 'Asia/Jakarta')::date, '1 day')::date AS day
		)
		SELECT to_char(d.day, 'YYYY-MM-DD') AS date,
		       COALESCE(SUM(o.grand_total), 0) AS revenue,
		       COUNT(o.id) AS orders
		FROM days d
		LEFT JOIN orders o ON (o.order_date AT TIME ZONE 'Asia/Jakarta')::date = d.day
		      AND o.deleted_at IS NULL AND o.status <> 'CANCELLED'
		GROUP BY d.day ORDER BY d.day`).
		Scan(&stats.SalesLast7Days).Error; err != nil {
		return nil, err
	}

	stats.TopProducts = []TopProduct{}
	if err := db.Raw(`
		SELECT COALESCE(NULLIF(oi.product_name, ''), p.name, 'Produk') AS product_name,
		       SUM(oi.quantity) AS quantity, SUM(oi.subtotal) AS revenue
		FROM order_items oi
		JOIN orders o ON o.id = oi.order_id AND o.deleted_at IS NULL AND o.status <> 'CANCELLED'
		LEFT JOIN products p ON p.id = oi.product_id
		WHERE o.order_date >= NOW() - INTERVAL '30 days'
		GROUP BY 1 ORDER BY quantity DESC LIMIT 5`).
		Scan(&stats.TopProducts).Error; err != nil {
		return nil, err
	}

	stats.RecentOrders = []RecentOrder{}
	if err := db.Raw(`
		SELECT o.id, o.order_number, COALESCE(u.full_name, '-') AS customer_name, o.grand_total, o.status, o.order_date
		FROM orders o LEFT JOIN users u ON u.id = o.user_id
		WHERE o.deleted_at IS NULL
		ORDER BY o.order_date DESC LIMIT 6`).
		Scan(&stats.RecentOrders).Error; err != nil {
		return nil, err
	}

	return stats, nil
}
