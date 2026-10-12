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

	// 2. Seed Default Categories
	var categoryCount int64
	db.Model(&models.Category{}).Count(&categoryCount)
	if categoryCount == 0 {
		categories := []models.Category{
			{Name: "Beras & Palawija", Slug: "beras-palawija", Description: "Beras premium, pulen, ketan, dan palawija", IsActive: true},
			{Name: "Minyak Goreng & Margarin", Slug: "minyak-goreng-margarin", Description: "Minyak goreng kelapa sawit, mentega, margarin", IsActive: true},
			{Name: "Gula & Pemanis", Slug: "gula-pemanis", Description: "Gula pasir putih, gula merah, pemanis alami", IsActive: true},
			{Name: "Telur & Hasil Tani", Slug: "telur-hasil-tani", Description: "Telur ayam negeri, telur bebek, hasil tani segar", IsActive: true},
			{Name: "Mie & Bumbu Dapur", Slug: "mie-bumbu-dapur", Description: "Mie instan, bihun, bumbu racik, rempah-rempah", IsActive: true},
		}
		for _, cat := range categories {
			c := cat
			db.Where("slug = ?", c.Slug).FirstOrCreate(&c)
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

	// 8. Seed Dummy Products
	var productCount int64
	db.Model(&models.Product{}).Where("is_active = true").Count(&productCount)
	if productCount == 0 {
		var catBeras, catMinyak, catGula, catTelur, catMie models.Category
		db.Where("name ILIKE ?", "%Beras%").First(&catBeras)
		db.Where("name ILIKE ?", "%Minyak%").First(&catMinyak)
		db.Where("name ILIKE ?", "%Gula%").First(&catGula)
		db.Where("name ILIKE ?", "%Telur%").First(&catTelur)
		db.Where("name ILIKE ?", "%Mie%").First(&catMie)

		var allCats []models.Category
		db.Find(&allCats)
		if len(allCats) > 0 {
			if catBeras.ID == (models.Category{}).ID {
				catBeras = allCats[0]
			}
			if catMinyak.ID == (models.Category{}).ID {
				catMinyak = allCats[0]
			}
			if catGula.ID == (models.Category{}).ID {
				catGula = allCats[0]
			}
			if catTelur.ID == (models.Category{}).ID {
				catTelur = allCats[0]
			}
			if catMie.ID == (models.Category{}).ID {
				catMie = allCats[0]
			}
		}

		var unitKrg, unitDus, unitKg, unitPcs, unitPck, unitBtl models.Unit
		db.Where("short_name = ?", "krg").First(&unitKrg)
		db.Where("short_name = ?", "dus").First(&unitDus)
		db.Where("short_name = ?", "kg").First(&unitKg)
		db.Where("short_name = ?", "pcs").First(&unitPcs)
		db.Where("short_name = ?", "pck").First(&unitPck)
		db.Where("short_name = ?", "btl").First(&unitBtl)

		var allUnits []models.Unit
		db.Find(&allUnits)
		if len(allUnits) > 0 {
			if unitKrg.ID == (models.Unit{}).ID {
				unitKrg = allUnits[0]
			}
			if unitDus.ID == (models.Unit{}).ID {
				unitDus = allUnits[0]
			}
			if unitKg.ID == (models.Unit{}).ID {
				unitKg = allUnits[0]
			}
			if unitPcs.ID == (models.Unit{}).ID {
				unitPcs = allUnits[0]
			}
			if unitPck.ID == (models.Unit{}).ID {
				unitPck = allUnits[0]
			}
			if unitBtl.ID == (models.Unit{}).ID {
				unitBtl = allUnits[0]
			}
		}

		dummyProducts := []struct {
			Product        models.Product
			Unit           models.Unit
			UnitName       string
			RetailPrice    float64
			WholesalePrice float64
			ImageURL       string
			Stock          int
		}{
			{
				Product: models.Product{CategoryID: catBeras.ID, SKU: "BRS-RJL-5KG", Name: "Beras Premium Rojo Lele 5kg", Description: "Beras pulen kualitas super, putih alami tanpa pemutih & pengawet.", IsActive: true},
				Unit:    unitKrg, UnitName: "Karung 5kg", RetailPrice: 72000, WholesalePrice: 69000, ImageURL: "https://images.unsplash.com/photo-1586201375761-83865001e31c?w=500&auto=format&fit=crop", Stock: 150,
			},
			{
				Product: models.Product{CategoryID: catBeras.ID, SKU: "BRS-PDW-10KG", Name: "Beras Pandan Wangi Super 10kg", Description: "Beras beraroma pandan alami, biji utuh & sangat pulen.", IsActive: true},
				Unit:    unitKrg, UnitName: "Karung 10kg", RetailPrice: 145000, WholesalePrice: 139000, ImageURL: "https://images.unsplash.com/photo-1536304993881-ff6e9eefa2a6?w=500&auto=format&fit=crop", Stock: 80,
			},
			{
				Product: models.Product{CategoryID: catMinyak.ID, SKU: "MYK-BML-2L", Name: "Minyak Goreng Bimoli Refill 2L", Description: "Minyak goreng kelapa sawit murni kualitas emulsi tinggi.", IsActive: true},
				Unit:    unitPcs, UnitName: "Pouch 2L", RetailPrice: 36500, WholesalePrice: 34500, ImageURL: "https://images.unsplash.com/photo-1474979266404-7eaacbcd87c5?w=500&auto=format&fit=crop", Stock: 200,
			},
			{
				Product: models.Product{CategoryID: catMinyak.ID, SKU: "MYK-MKT-1L", Name: "Minyak Goreng Minyakita 1L", Description: "Minyak goreng program pemerintah, bening & ekonomis.", IsActive: true},
				Unit:    unitPcs, UnitName: "Pouch 1L", RetailPrice: 15000, WholesalePrice: 14000, ImageURL: "https://images.unsplash.com/photo-1620706857370-e1b9770e8bb1?w=500&auto=format&fit=crop", Stock: 350,
			},
			{
				Product: models.Product{CategoryID: catGula.ID, SKU: "GLA-GLK-1KG", Name: "Gula Pasir Gulaku Premium 1kg", Description: "Gula pasir tebu pilihan, kristal putih bersih & manis alami.", IsActive: true},
				Unit:    unitPck, UnitName: "Pack 1kg", RetailPrice: 17500, WholesalePrice: 16500, ImageURL: "https://images.unsplash.com/photo-1622484210800-410a8b9e6b4f?w=500&auto=format&fit=crop", Stock: 250,
			},
			{
				Product: models.Product{CategoryID: catTelur.ID, SKU: "TLR-AYM-1KG", Name: "Telur Ayam Negeri Super (1kg)", Description: "Telur ayam murni segar langsung dari peternakan utama.", IsActive: true},
				Unit:    unitKg, UnitName: "Kilogram", RetailPrice: 28000, WholesalePrice: 26500, ImageURL: "https://images.unsplash.com/photo-1582722872445-44dc5f7e3c8f?w=500&auto=format&fit=crop", Stock: 100,
			},
			{
				Product: models.Product{CategoryID: catMie.ID, SKU: "MIE-IDM-GRG", Name: "Indomie Goreng Spesial (Karton / 40 Pcs)", Description: "Mie instan goreng favorit Indonesia dengan bumbu lezat.", IsActive: true},
				Unit:    unitDus, UnitName: "Dus (40 Pcs)", RetailPrice: 112000, WholesalePrice: 108000, ImageURL: "https://images.unsplash.com/photo-1612929633738-8fe44f7ec841?w=500&auto=format&fit=crop", Stock: 120,
			},
			{
				Product: models.Product{CategoryID: catMie.ID, SKU: "TPG-SGB-1KG", Name: "Tepung Terigu Segitiga Biru 1kg", Description: "Tepung terigu serbaguna untuk aneka kue, martabak & gorengan.", IsActive: true},
				Unit:    unitPck, UnitName: "Pack 1kg", RetailPrice: 12500, WholesalePrice: 11800, ImageURL: "https://images.unsplash.com/photo-1509440159596-0249088772ff?w=500&auto=format&fit=crop", Stock: 300,
			},
		}

		for _, item := range dummyProducts {
			p := item.Product
			if err := db.Create(&p).Error; err == nil {
				// Create Product Unit
				pu := models.ProductUnit{
					ProductID:        p.ID,
					UnitID:           item.Unit.ID,
					UnitName:         item.UnitName,
					ConversionFactor: 1,
					RetailPrice:      item.RetailPrice,
					WholesalePrice:   item.WholesalePrice,
					MinWholesaleQty:  5,
					StockQuantity:    item.Stock,
				}
				db.Create(&pu)

				// Create Product Image
				pi := models.ProductImage{
					ProductID: p.ID,
					ImageURL:  item.ImageURL,
					IsPrimary: true,
					SortOrder: 1,
				}
				db.Create(&pi)

				// Create Inventory Record
				if defaultWarehouse.ID != [16]byte{} {
					inv := models.Inventory{
						WarehouseID:   defaultWarehouse.ID,
						ProductID:     p.ID,
						ProductUnitID: pu.ID,
						CurrentStock:  item.Stock,
						ReservedStock: 0,
						MinStockAlert: 10,
					}
					db.Create(&inv)
				}
			}
		}
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
