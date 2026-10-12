package config

import (
	"os"

	"backend-sembako/helpers"
	"backend-sembako/internal/constants"
	"backend-sembako/internal/models"
	"backend-sembako/pkg/logger"

	"gorm.io/gorm"
)

func SeedDatabase(db *gorm.DB) error {
	logger.Log.Info().Msg("Executing LuckyMart Initial Database Seeder...")

	// 1. Seed Super Admin User
	var adminCount int64
	db.Model(&models.User{}).Where("role = ?", constants.RoleSuperAdmin).Count(&adminCount)
	if adminCount == 0 {
		// Set SEED_ADMIN_PASSWORD outside local development; the fallback is a
		// development-only default and must be changed after the first login.
		adminPassword := os.Getenv("SEED_ADMIN_PASSWORD")
		if adminPassword == "" {
			adminPassword = "SuperAdmin123!"
			logger.Log.Warn().Msg("SEED_ADMIN_PASSWORD not set — seeding the development default admin password; change it after first login")
		}
		hashedPass, _ := helpers.HashPassword(adminPassword)
		adminEmail := "admin@sembako.com"
		adminUser := &models.User{
			FullName:        "Super Admin Sembako",
			PhoneNumber:     "081234567890",
			Email:           &adminEmail,
			PasswordHash:    hashedPass,
			Role:            constants.RoleSuperAdmin,
			IsEmailVerified: true,
		}
		if err := db.Create(adminUser).Error; err == nil {
			logger.Log.Info().Msg("Default Super Admin created with phone 081234567890 (password from SEED_ADMIN_PASSWORD)")
		}
	}

	// 3. Seed Default Payment Methods (added by code, so new methods also
	// reach existing databases; admins may later edit or deactivate them).
	methods := []models.PaymentMethod{
		{Code: "COD", Name: "Bayar di Tempat (COD)", Type: "COD", IsActive: true, Description: "Bayar tunai ke kurir toko saat pesanan tiba"},
		{Code: "MANUAL_BANK_TRANSFER", Name: "Transfer Bank", Type: "MANUAL_TRANSFER", IsActive: true, Description: "Transfer ke rekening toko lalu unggah bukti transfer"},
		{Code: "B2B_CREDIT_TEMPO", Name: "Kasbon Tempo 30 Hari", Type: "B2B_CREDIT", IsActive: true, Description: "Khusus mitra toko yang sudah disetujui"},
	}
	for _, m := range methods {
		method := m
		db.Where("code = ?", method.Code).FirstOrCreate(&method)
	}

	// 4. Seed Default Couriers & Shipping Methods
	var courierCount int64
	db.Model(&models.Courier{}).Count(&courierCount)
	if courierCount == 0 {
		courier := &models.Courier{
			Code:     "KURIR_TOKO",
			Name:     "Kurir Toko Internal",
			IsActive: true,
		}
		if err := db.Create(courier).Error; err == nil {
			method := &models.ShippingMethod{
				CourierID:     courier.ID,
				Code:          "REGULAR_SAME_DAY",
				Name:          "Pengiriman Reguler Same-Day",
				Description:   "Pengiriman sembako langsung di hari yang sama",
				EstimatedDays: "1 Hari",
				BaseRate:      15000,
				IsActive:      true,
			}
			db.Create(method)
		}
	}

	// 5. Seed Default Notification Templates
	var tmplCount int64
	db.Model(&models.NotificationTemplate{}).Count(&tmplCount)
	if tmplCount == 0 {
		templates := []models.NotificationTemplate{
			{
				Name:      "Template Order Berhasil",
				EventType: "ORDER_CREATED",
				Channel:   "IN_APP",
				Subject:   "Pesanan {{order_number}} Berhasil Dibuat",
				Body:      "Halo {{customer_name}}, pesanan Anda #{{order_number}} senilai {{grand_total}} telah berhasil dibuat dan dikunci.",
				IsActive:  true,
			},
			{
				Name:      "Template Pembayaran Diverifikasi",
				EventType: "PAYMENT_PAID",
				Channel:   "IN_APP",
				Subject:   "Pembayaran Pesanan {{order_number}} Lunas",
				Body:      "Pembayaran sebesar {{grand_total}} telah diterima. Pesanan Anda kini disiapkan oleh tim gudang.",
				IsActive:  true,
			},
		}
		db.Create(&templates)
	}

	// 6. Seed Units
	var unitCount int64
	db.Model(&models.Unit{}).Count(&unitCount)
	if unitCount == 0 {
		units := []models.Unit{
			{Name: "Karung", ShortName: "krg"},
			{Name: "Karton", ShortName: "dus"},
			{Name: "Kilogram", ShortName: "kg"},
			{Name: "Pcs", ShortName: "pcs"},
			{Name: "Pack", ShortName: "pck"},
			{Name: "Botol", ShortName: "btl"},
		}
		db.Create(&units)
	}

	// 7. Seed Default Warehouse
	var whCount int64
	db.Model(&models.Warehouse{}).Count(&whCount)
	var defaultWarehouse models.Warehouse
	if whCount == 0 {
		defaultWarehouse = models.Warehouse{
			Code:     "WH-MAIN-01",
			Name:     "Gudang Utama LuckyMart",
			Address:  "Jl. Industri LuckyMart No. 88, Medan",
			IsActive: true,
		}
		db.Create(&defaultWarehouse)
	} else {
		db.First(&defaultWarehouse)
	}

	repairData(db, defaultWarehouse)

	logger.Log.Info().Msg("Database Seeding Completed Successfully.")
	return nil
}

// repairData fixes rows written by earlier versions of the app. Every step
// is idempotent, so it is safe to run on each start.
func repairData(db *gorm.DB, warehouse models.Warehouse) {
	// Product units created without an inventory row (older product form).
	if warehouse.ID != [16]byte{} {
		res := db.Exec(`
			INSERT INTO inventories (id, warehouse_id, product_id, product_unit_id, current_stock, reserved_stock, min_stock_alert, created_at, updated_at)
			SELECT gen_random_uuid(), ?, pu.product_id, pu.id, pu.stock_quantity, 0, 10, NOW(), NOW()
			FROM product_units pu
			WHERE pu.deleted_at IS NULL
			  AND NOT EXISTS (SELECT 1 FROM inventories i WHERE i.product_unit_id = pu.id AND i.deleted_at IS NULL)`, warehouse.ID)
		if res.Error == nil && res.RowsAffected > 0 {
			logger.Log.Info().Int64("rows", res.RowsAffected).Msg("Created missing inventory rows")
		}
	}

	// Inventory must mirror product_units.stock_quantity, which checkout uses.
	db.Exec(`
		UPDATE inventories i SET current_stock = pu.stock_quantity, updated_at = NOW()
		FROM product_units pu
		WHERE i.product_unit_id = pu.id AND i.current_stock <> pu.stock_quantity`)

	// Order items from before name snapshots were stored.
	db.Exec(`
		UPDATE order_items oi SET product_name = p.name, sku = p.sku
		FROM products p
		WHERE oi.product_id = p.id AND oi.product_name = ''`)
	db.Exec(`
		UPDATE order_items oi SET unit_name = pu.unit_name
		FROM product_units pu
		WHERE oi.product_unit_id = pu.id AND oi.unit_name = ''`)

	// Existing kasbon history entries are credit usage.
	db.Exec(`UPDATE credit_payments SET type = 'CREDIT_USAGE' WHERE type IS NULL OR type = ''`)
}
