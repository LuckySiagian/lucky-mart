package config

import (
	"fmt"

	"backend-sembako/internal/models"
	"backend-sembako/pkg/logger"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

func InitDatabase(cfg *DatabaseConfig, env string) (*gorm.DB, error) {
	// Managed Postgres providers (Neon, Supabase, Render, Railway) hand out a
	// single connection URL; use it when present.
	dsn := cfg.URL
	if dsn == "" {
		dsn = fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=Asia/Jakarta",
			cfg.Host, cfg.User, cfg.Password, cfg.Name, cfg.Port, cfg.SSLMode,
		)
	}

	var logLevel gormlogger.LogLevel
	switch env {
	case "production":
		logLevel = gormlogger.Error
	case "staging":
		logLevel = gormlogger.Warn
	default:
		logLevel = gormlogger.Info
	}

	gormCfg := &gorm.Config{
		Logger: gormlogger.Default.LogMode(logLevel),
	}

	db, err := gorm.Open(postgres.Open(dsn), gormCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to PostgreSQL database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB instance: %w", err)
	}

	// Configure Connection Pool
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	logger.Log.Info().
		Int("max_open_conns", cfg.MaxOpenConns).
		Int("max_idle_conns", cfg.MaxIdleConns).
		Msg("PostgreSQL Connection Pool initialized successfully")

	return db, nil
}

func AutoMigrate(db *gorm.DB) error {
	logger.Log.Info().Msg("Starting PostgreSQL Auto-Migration for all domain entities...")

	err := db.AutoMigrate(
		&models.User{},
		&models.RefreshToken{},
		&models.LoginHistory{},
		&models.B2BProfile{},
		&models.Permission{},
		&models.Category{},
		&models.Brand{},
		&models.Unit{},
		&models.Supplier{},
		&models.Warehouse{},
		&models.Product{},
		&models.ProductUnit{},
		&models.ProductImage{},
		&models.Inventory{},
		&models.StockMovement{},
		&models.Cart{},
		&models.CartItem{},
		&models.Order{},
		&models.OrderItem{},
		&models.OrderAddress{},
		&models.OrderHistory{},
		&models.Invoice{},
		&models.Address{},
		&models.Courier{},
		&models.ShippingMethod{},
		&models.ShippingRate{},
		&models.Shipment{},
		&models.ShipmentTracking{},
		&models.Payment{},
		&models.PaymentMethod{},
		&models.PaymentTransaction{},
		&models.PaymentLog{},
		&models.CreditLedger{},
		&models.CreditPayment{},
		&models.Notification{},
		&models.NotificationTemplate{},
		&models.NotificationLog{},
		&models.EventLog{},
		&models.StoreSetting{},
	)
	if err != nil {
		return err
	}

	// AutoMigrate never relaxes NOT NULL, so columns that became optional are
	// fixed here for databases created by older versions. Each statement is a
	// no-op when already applied.
	for _, stmt := range []string{
		// Kasbon usage and reversals are not tied to a payment row.
		`ALTER TABLE credit_payments ALTER COLUMN payment_id DROP NOT NULL`,
	} {
		if err := db.Exec(stmt).Error; err != nil {
			return err
		}
	}
	return nil
}
