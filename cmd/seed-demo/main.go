// Command seed-demo fills a development database with realistic demo data
// (catalogue, customers, staff, orders in every state, kasbon) so every
// feature can be tried. Everything is created through the real API handlers,
// so stock, invoices, kasbon and notifications stay consistent; afterwards
// past orders are moved back in time so the dashboard shows two weeks of
// history.
//
//	go run ./cmd/seed-demo
//
// It refuses to run with APP_ENV=production and does nothing when the demo
// catalogue already exists.
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math/rand"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"time"

	"backend-sembako/config"
	"backend-sembako/internal/constants"
	"backend-sembako/internal/models"
	"backend-sembako/internal/routes"
	"backend-sembako/pkg/jwt"
	"backend-sembako/pkg/logger"

	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// DemoPassword is the password of every demo account (see docs/AKUN-DEMO.md).
const DemoPassword = "LuckyDemo#2026"

type seeder struct {
	db       *gorm.DB
	handler  http.Handler
	admin    string // super admin access token
	products map[string]*product
	rnd      *rand.Rand
}

type product struct {
	ID    string
	Name  string
	Units map[string]string // unit name -> product unit id
}

type account struct {
	ID, Name, Phone, Token, AddressID string
}

type line struct {
	SKU, Unit string
	Qty       int
}

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		fail("konfigurasi: %v", err)
	}
	if cfg.App.Env == "production" {
		fail("seed-demo tidak boleh dijalankan di production")
	}
	logger.InitLogger("production") // keep request logs quiet
	// "production" here only quiets the SQL log; the refusal above still
	// guards the real environment.
	db, err := config.InitDatabase(&cfg.Database, "production")
	if err != nil {
		fail("database: %v", err)
	}
	// Expected "record not found" lookups would otherwise flood the output.
	db = db.Session(&gorm.Session{Logger: gormlogger.Default.LogMode(gormlogger.Silent)})
	if err := config.AutoMigrate(db); err != nil {
		fail("migrasi: %v", err)
	}
	if err := config.SeedDatabase(db); err != nil {
		fail("seed dasar: %v", err)
	}

	var existing int64
	db.Model(&models.Product{}).Unscoped().Where("sku = ?", "LM-BRS-001").Count(&existing)
	if existing > 0 {
		fmt.Println("Data demo sudah ada (produk LM-BRS-001 ditemukan). Tidak ada yang diubah.")
		return
	}

	var admin models.User
	if err := db.Where("role = ?", constants.RoleSuperAdmin).Order("created_at").First(&admin).Error; err != nil {
		fail("super admin tidak ditemukan: %v", err)
	}
	tokens, err := jwt.NewJWTService(cfg.JWT.AccessSecret, cfg.JWT.RefreshSecret, time.Hour, time.Hour).
		GenerateTokenPair(admin.ID, admin.Role)
	if err != nil {
		fail("token admin: %v", err)
	}

	s := &seeder{
		db:       db,
		handler:  routes.SetupRouter(cfg, db),
		admin:    tokens.AccessToken,
		products: map[string]*product{},
		rnd:      rand.New(rand.NewSource(20261012)),
	}
	s.run()
	// Notifications are written by background event listeners.
	time.Sleep(2 * time.Second)
}

func (s *seeder) run() {
	fmt.Println("1/6 Info toko, pembayaran, ongkir")
	s.storeSettings()

	fmt.Println("2/6 Kategori dan produk")
	s.catalog()

	fmt.Println("3/6 Akun staf")
	for _, st := range []struct {
		name, phone string
		role        constants.UserRole
	}{
		{"Maria Simanjuntak", "089900000301", constants.RoleStoreManager},
		{"Andi Pratama", "089900000302", constants.RoleWarehouseStaff},
		{"Joko Susilo", "089900000303", constants.RoleDriver},
	} {
		s.must("POST", "/users", s.admin, map[string]any{
			"full_name": st.name, "phone_number": st.phone, "password": DemoPassword, "role": st.role,
		})
	}

	fmt.Println("4/6 Pelanggan dan mitra warung")
	retail := []*account{
		s.customer("Sari Wulandari", "089900000101", "Jl. Melati No. 5, RT 02/RW 03, Kel. Sukaramai", false, ""),
		s.customer("Budi Santoso", "089900000102", "Jl. Kenanga No. 17, Kel. Sei Kera Hilir", false, ""),
		s.customer("Dewi Lestari", "089900000103", "Gg. Mawar No. 3, Kel. Pahlawan", false, ""),
		s.customer("Ahmad Fauzi", "089900000104", "Jl. Pelita II No. 21, Kel. Sidorame", false, ""),
		s.customer("Rina Marlina", "089900000105", "Jl. Sutomo Ujung No. 8, Kel. Durian", false, ""),
		s.customer("Yohana Sitompul", "089900000106", "Jl. Bilal No. 44, Kel. Pulo Brayan", false, ""),
	}
	tini := s.customer("Tini Rahayu", "089900000201", "Jl. Gaharu No. 12 (Warung Bu Tini)", true, "Warung Bu Tini")
	harun := s.customer("Harun Nasution", "089900000202", "Jl. Krakatau No. 90 (Toko Sinar Jaya)", true, "Toko Sinar Jaya")
	s.customer("Rosmawati", "089900000203", "Jl. Cemara No. 2 (Kedai Mak Ros)", true, "Kedai Mak Ros") // waits for verification
	s.must("PATCH", "/users/"+tini.ID+"/b2b", s.admin, map[string]any{"status": "APPROVED", "credit_limit": 3000000})
	s.must("PATCH", "/users/"+harun.ID+"/b2b", s.admin, map[string]any{"status": "APPROVED", "credit_limit": 5000000})

	fmt.Println("5/6 Pesanan dua minggu terakhir")
	s.orders(retail, tini, harun)

	fmt.Println("6/6 Barang masuk dari supplier")
	for _, in := range []struct {
		sku, unit string
		qty       int
	}{
		{"LM-GLA-001", "Pak 1 kg", 50}, {"LM-MIE-001", "Dus (40 pcs)", 10}, {"LM-MYK-004", "Pouch 1 L", 48}, {"LM-TLR-001", "Kg", 30},
	} {
		s.must("POST", "/inventory/adjust", s.admin, map[string]any{
			"product_unit_id": s.unit(in.sku, in.unit), "adjustment_type": "STOCK_IN", "quantity_changed": in.qty,
			"reference_number": fmt.Sprintf("PO-%s-%03d", time.Now().Format("0102"), s.rnd.Intn(900)+100),
			"notes":            "Barang masuk dari supplier",
		})
	}

	fmt.Println("\nSelesai! Akun demo dan cara mencobanya ada di docs/AKUN-DEMO.md")
}

func (s *seeder) storeSettings() {
	s.must("PUT", "/store", s.admin, map[string]any{
		"name":              "Lucky Mart",
		"tagline":           "Toko sembako keluarga — eceran & grosir, bisa diantar atau ambil sendiri.",
		"address":           "Jl. Contoh Raya No. 12, Kel. Sukamaju (alamat contoh, ganti dengan alamat toko)",
		"whatsapp":          "",
		"open_hours":        "Setiap hari 07.00 – 21.00",
		"maps_url":          "",
		"receipt_note":      "Terima kasih sudah berbelanja di Lucky Mart!",
		"announcement":      "Gratis ongkir untuk belanja minimal Rp300.000!",
		"free_shipping_min": 300000,
	})

	var methods []models.PaymentMethod
	s.db.Find(&methods)
	for _, m := range methods {
		body := map[string]any{"name": m.Name, "description": m.Description, "is_active": m.IsActive}
		switch m.Code {
		case string(constants.PaymentMethodManualBankTransfer):
			body["name"] = "Transfer Bank"
			body["description"] = "BCA 1234567890 a.n. Lucky Mart (CONTOH — ganti dengan rekening asli toko)"
		case string(constants.PaymentMethodCashOnDelivery):
			body["name"] = "Bayar Tunai"
		case string(constants.PaymentMethodB2BCreditTempo):
		default:
			body["is_active"] = false // older duplicate methods would confuse customers
		}
		s.must("PUT", "/payment-methods/"+m.ID.String(), s.admin, body)
	}
}

func (s *seeder) catalog() {
	var categories []map[string]any
	s.decode(s.must("GET", "/categories", "", nil), &categories)
	existing := map[string]string{}
	for _, c := range categories {
		existing[fmt.Sprint(c["name"])] = fmt.Sprint(c["id"])
	}

	for _, cat := range catalog {
		catID, ok := existing[cat.Name]
		if !ok {
			var created struct {
				ID string `json:"id"`
			}
			s.decode(s.must("POST", "/categories", s.admin, map[string]any{"name": cat.Name, "description": cat.Description}), &created)
			catID = created.ID
		}
		for _, p := range cat.Products {
			units := make([]map[string]any, len(p.Units))
			for i, u := range p.Units {
				units[i] = map[string]any{
					"unit_name": u.Name, "conversion_factor": u.Conversion, "retail_price": u.Retail,
					"wholesale_price": u.Wholesale, "min_wholesale_qty": u.MinWhole, "initial_stock": u.Stock,
				}
			}
			var created struct {
				ID    string `json:"id"`
				Units []struct {
					ID       string `json:"id"`
					UnitName string `json:"unit_name"`
				} `json:"units"`
			}
			s.decode(s.must("POST", "/products", s.admin, map[string]any{
				"name": p.Name, "sku": p.SKU, "category_id": catID, "description": p.Description,
				"min_stock_alert": p.MinStock, "units": units,
			}), &created)
			pr := &product{ID: created.ID, Name: p.Name, Units: map[string]string{}}
			for _, u := range created.Units {
				pr.Units[u.UnitName] = u.ID
			}
			s.products[p.SKU] = pr
		}
	}
}

func (s *seeder) customer(name, phone, street string, b2b bool, storeName string) *account {
	body := map[string]any{"full_name": name, "phone_number": phone, "password": DemoPassword}
	if b2b {
		body["is_b2b_registration"] = true
		body["store_name"] = storeName
		body["store_address"] = street
	}
	var res struct {
		AccessToken string `json:"access_token"`
		User        struct {
			ID string `json:"id"`
		} `json:"user"`
	}
	s.decode(s.must("POST", "/auth/register", "", body), &res)
	acc := &account{ID: res.User.ID, Name: name, Phone: phone, Token: res.AccessToken}

	var addr struct {
		ID string `json:"id"`
	}
	s.decode(s.must("POST", "/addresses", acc.Token, map[string]any{
		"title": map[bool]string{true: "Warung", false: "Rumah"}[b2b], "recipient_name": name, "phone": phone,
		"full_address": street, "city": "Medan", "province": "Sumatera Utara", "postal_code": "20111", "is_default": true,
	}), &addr)
	acc.AddressID = addr.ID
	return acc
}

func (s *seeder) unit(sku, unit string) string {
	p, ok := s.products[sku]
	if !ok || p.Units[unit] == "" {
		fail("produk %s satuan %q tidak ada di katalog demo", sku, unit)
	}
	return p.Units[unit]
}

type orderPlan struct {
	DaysAgo  int
	Hour     int
	Customer *account
	Items    []line
	Payment  string // COD, MANUAL_BANK_TRANSFER, B2B_CREDIT_TEMPO
	Pickup   bool
	Notes    string
	// Final state: COMPLETED, CONFIRMED, PROCESSING_PACKING, OUT_FOR_DELIVERY,
	// PROOF_SENT (transfer waiting for check), UNPAID (transfer not paid yet),
	// CANCELLED_CUSTOMER, CANCELLED_STORE.
	State string
}

func (s *seeder) orders(retail []*account, tini, harun *account) {
	basket := [][]line{
		{{"LM-BRS-002", "Karung 5 kg", 1}, {"LM-MYK-004", "Pouch 1 L", 2}, {"LM-TLR-001", "Kg", 1}},
		{{"LM-MIE-001", "Pcs", 10}, {"LM-TLR-001", "Kg", 1}, {"LM-BMB-005", "1/4 kg", 1}},
		{{"LM-GLA-001", "Pak 1 kg", 2}, {"LM-KTS-001", "Renceng (10 sachet)", 1}, {"LM-KTS-004", "Kaleng 370 g", 1}},
		{{"LM-BRS-001", "Karung 5 kg", 1}, {"LM-MYK-001", "Pouch 2 L", 1}, {"LM-BMB-002", "Pouch 520 ml", 1}},
		{{"LM-MNM-002", "Galon", 2}, {"LM-RMH-002", "Bungkus 770 g", 1}, {"LM-RMH-001", "Batang", 3}},
		{{"LM-TPG-001", "Pak 1 kg", 2}, {"LM-GLA-001", "Pak 1 kg", 1}, {"LM-TLR-001", "Tray (30 butir)", 1}, {"LM-MYK-003", "Pouch 2 L", 1}},
		{{"LM-MIE-003", "Pcs", 5}, {"LM-MNM-001", "Botol", 6}, {"LM-BMB-001", "Pak 250 g", 1}},
		{{"LM-BRS-001", "Karung 5 kg", 2}, {"LM-MYK-002", "Pouch 2 L", 2}, {"LM-GLA-001", "Pak 1 kg", 3}, {"LM-TLR-001", "Tray (30 butir)", 1}},
	}
	wholesale := [][]line{
		{{"LM-MIE-001", "Dus (40 pcs)", 2}, {"LM-MYK-004", "Dus (12 pouch)", 1}, {"LM-GLA-001", "Pak 1 kg", 10}},
		{{"LM-BRS-003", "Karung 25 kg", 2}, {"LM-MYK-001", "Dus (6 pouch)", 1}},
		{{"LM-MNM-001", "Dus (24 botol)", 3}, {"LM-KTS-001", "Renceng (10 sachet)", 10}, {"LM-MIE-002", "Dus (40 pcs)", 1}},
		{{"LM-TLR-001", "Kg", 10}, {"LM-TPG-001", "Pak 1 kg", 10}, {"LM-GLA-002", "Plastik 1 kg", 10}},
	}

	var plans []orderPlan
	payments := []string{"COD", "COD", "COD", "MANUAL_BANK_TRANSFER"}
	// Completed retail orders spread over the past two weeks.
	for day := 14; day >= 1; day-- {
		for n := 0; n < 1+s.rnd.Intn(3); n++ {
			plans = append(plans, orderPlan{
				DaysAgo: day, Hour: 8 + s.rnd.Intn(12), Customer: retail[s.rnd.Intn(len(retail))],
				Items: basket[s.rnd.Intn(len(basket))], Payment: payments[s.rnd.Intn(len(payments))],
				Pickup: s.rnd.Intn(3) == 0, State: "COMPLETED",
			})
		}
	}
	// Mitra warung buying wholesale on kasbon and by transfer.
	plans = append(plans,
		orderPlan{DaysAgo: 35, Hour: 9, Customer: tini, Items: wholesale[0], Payment: "B2B_CREDIT_TEMPO", State: "COMPLETED", Notes: "Antar pagi sebelum warung buka"},
		orderPlan{DaysAgo: 12, Hour: 10, Customer: harun, Items: wholesale[1], Payment: "B2B_CREDIT_TEMPO", State: "COMPLETED"},
		orderPlan{DaysAgo: 8, Hour: 14, Customer: tini, Items: wholesale[2], Payment: "B2B_CREDIT_TEMPO", State: "COMPLETED"},
		orderPlan{DaysAgo: 4, Hour: 11, Customer: harun, Items: wholesale[3], Payment: "MANUAL_BANK_TRANSFER", State: "COMPLETED"},
	)
	// Cancelled examples.
	plans = append(plans,
		orderPlan{DaysAgo: 6, Hour: 16, Customer: retail[1], Items: basket[1], Payment: "COD", State: "CANCELLED_CUSTOMER"},
		orderPlan{DaysAgo: 3, Hour: 19, Customer: retail[4], Items: basket[4], Payment: "MANUAL_BANK_TRANSFER", State: "CANCELLED_STORE"},
	)
	// Today's work for the store: one order in every step.
	plans = append(plans,
		orderPlan{DaysAgo: 0, Hour: 7, Customer: retail[0], Items: basket[7], Payment: "COD", State: "OUT_FOR_DELIVERY", Notes: "Rumah pagar hijau"},
		orderPlan{DaysAgo: 0, Hour: 7, Customer: retail[2], Items: basket[2], Payment: "COD", Pickup: true, State: "OUT_FOR_DELIVERY", Notes: "Saya ambil jam 5 sore"},
		orderPlan{DaysAgo: 0, Hour: 8, Customer: harun, Items: wholesale[0], Payment: "B2B_CREDIT_TEMPO", State: "PROCESSING_PACKING"},
		orderPlan{DaysAgo: 0, Hour: 8, Customer: retail[3], Items: basket[0], Payment: "COD", State: "CONFIRMED"},
		orderPlan{DaysAgo: 0, Hour: 8, Customer: retail[5], Items: basket[6], Payment: "COD", Pickup: true, State: "CONFIRMED"},
		orderPlan{DaysAgo: 0, Hour: 8, Customer: retail[1], Items: basket[5], Payment: "MANUAL_BANK_TRANSFER", State: "PROOF_SENT"},
		orderPlan{DaysAgo: 0, Hour: 8, Customer: retail[4], Items: basket[3], Payment: "MANUAL_BANK_TRANSFER", Pickup: true, State: "UNPAID"},
	)

	var proof string
	for i, p := range plans {
		if p.Payment == "MANUAL_BANK_TRANSFER" && proof == "" {
			proof = s.uploadProof(p.Customer.Token)
		}
		orderID, orderNumber, paymentID := s.place(p)
		switch p.State {
		case "COMPLETED":
			if p.Payment == "MANUAL_BANK_TRANSFER" {
				s.pay(p.Customer.Token, paymentID, proof)
			}
			s.advance(orderID, p.Pickup, "COMPLETED")
		case "OUT_FOR_DELIVERY", "PROCESSING_PACKING":
			s.advance(orderID, p.Pickup, p.State)
		case "PROOF_SENT":
			s.must("POST", "/payments/upload-proof", p.Customer.Token, map[string]any{"payment_id": paymentID, "proof_image_url": proof})
		case "CANCELLED_CUSTOMER":
			s.must("PUT", "/orders/"+orderID+"/cancel", p.Customer.Token, map[string]any{"reason": "Salah pilih barang, mau pesan ulang"})
		case "CANCELLED_STORE":
			s.must("PUT", "/orders/"+orderID+"/cancel", s.admin, map[string]any{"reason": "Tidak ada pembayaran setelah 1 hari"})
		}
		s.backdate(orderID, orderNumber, p.DaysAgo, p.Hour, i)
	}

	// Mitra pays part of their kasbon.
	var ledger models.CreditLedger
	if err := s.db.Where("user_id = ?", tini.ID).First(&ledger).Error; err != nil {
		fail("kasbon %s: %v", tini.Name, err)
	}
	s.must("POST", "/credit-ledgers/"+ledger.ID.String()+"/settle", s.admin, map[string]any{"amount": 500000, "notes": "Bayar tunai di toko"})
}

func (s *seeder) place(p orderPlan) (orderID, orderNumber, paymentID string) {
	s.must("DELETE", "/cart", p.Customer.Token, nil)
	for _, it := range p.Items {
		s.must("POST", "/cart/items", p.Customer.Token, map[string]any{
			"product_id": s.products[it.SKU].ID, "product_unit_id": s.unit(it.SKU, it.Unit), "quantity": it.Qty,
		})
	}
	body := map[string]any{"payment_method": p.Payment, "notes": p.Notes, "delivery_type": "DELIVERY", "address_id": p.Customer.AddressID}
	if p.Pickup {
		body = map[string]any{"payment_method": p.Payment, "notes": p.Notes, "delivery_type": "PICKUP"}
	}
	var order struct {
		ID          string `json:"id"`
		OrderNumber string `json:"order_number"`
		Payments    []struct {
			ID string `json:"id"`
		} `json:"payments"`
	}
	s.decode(s.must("POST", "/orders", p.Customer.Token, body), &order)
	return order.ID, order.OrderNumber, order.Payments[0].ID
}

func (s *seeder) pay(token, paymentID, proof string) {
	s.must("POST", "/payments/upload-proof", token, map[string]any{"payment_id": paymentID, "proof_image_url": proof})
	s.must("PATCH", "/payments/"+paymentID+"/approve", s.admin, map[string]any{"notes": "Dana sudah masuk"})
}

// advance moves an order through the store's steps up to target.
func (s *seeder) advance(orderID string, pickup bool, target string) {
	for _, st := range []string{"PROCESSING_PACKING", "OUT_FOR_DELIVERY", "COMPLETED"} {
		s.must("PATCH", "/orders/"+orderID+"/status", s.admin, map[string]any{"status": st})
		if st == target {
			return
		}
	}
}

// backdate shifts an order and everything recorded for it to daysAgo at the
// given hour, keeping the gaps between its steps.
func (s *seeder) backdate(orderID, orderNumber string, daysAgo, hour, seq int) {
	now := time.Now()
	target := time.Date(now.Year(), now.Month(), now.Day(), hour, (seq*7)%60, 0, 0, now.Location()).AddDate(0, 0, -daysAgo)
	if target.After(now) {
		target = now.Add(-time.Duration(30-seq%30) * time.Minute)
	}
	var created time.Time
	s.db.Model(&models.Order{}).Select("created_at").Where("id = ?", orderID).Scan(&created)
	shift := fmt.Sprintf("%d seconds", int(created.Sub(target).Seconds()))
	if created.Sub(target) <= 0 {
		return
	}
	exec := func(q string, args ...any) {
		if err := s.db.Exec(q, append([]any{shift}, args...)...).Error; err != nil {
			fail("backdate: %v", err)
		}
	}
	exec(`UPDATE orders SET order_date = order_date - ?::interval, created_at = created_at - ?::interval, updated_at = updated_at - ?::interval, due_date = due_date - ?::interval WHERE id = ?`, shift, shift, shift, orderID)
	exec(`UPDATE order_histories SET created_at = created_at - ?::interval WHERE order_id = ?`, orderID)
	exec(`UPDATE invoices SET issued_at = issued_at - ?::interval, created_at = created_at - ?::interval, due_date = due_date - ?::interval, paid_at = paid_at - ?::interval WHERE order_id = ?`, shift, shift, shift, orderID)
	exec(`UPDATE payments SET created_at = created_at - ?::interval, paid_at = paid_at - ?::interval, expired_at = expired_at - ?::interval WHERE order_id = ?`, shift, shift, orderID)
	exec(`UPDATE payment_logs SET created_at = created_at - ?::interval WHERE payment_id IN (SELECT id FROM payments WHERE order_id = ?)`, orderID)
	exec(`UPDATE shipments SET created_at = created_at - ?::interval, shipped_at = shipped_at - ?::interval, delivered_at = delivered_at - ?::interval WHERE order_id = ?`, shift, shift, orderID)
	exec(`UPDATE shipment_trackings SET created_at = created_at - ?::interval WHERE shipment_id IN (SELECT id FROM shipments WHERE order_id = ?)`, orderID)
	exec(`UPDATE stock_movements SET created_at = created_at - ?::interval WHERE reference_number IN (?, ?)`, orderNumber, "CANCEL-"+orderNumber)
	exec(`UPDATE credit_payments SET created_at = created_at - ?::interval, payment_date = payment_date - ?::interval WHERE order_id = ?`, shift, orderID)
	exec(`UPDATE notifications SET created_at = created_at - ?::interval WHERE message LIKE ? OR title LIKE ?`, "%"+orderNumber+"%", "%"+orderNumber+"%")

	// Numbers carry the date (ORD-YYYYMMDD-xxxxxx); make them match the new one.
	oldDate, newDate := created.Format("20060102"), target.Format("20060102")
	if oldDate == newDate {
		return
	}
	newNumber := strings.Replace(orderNumber, oldDate, newDate, 1)
	rename := func(q string, args ...any) {
		if err := s.db.Exec(q, args...).Error; err != nil {
			fail("nomor pesanan: %v", err)
		}
	}
	rename(`UPDATE orders SET order_number = replace(order_number, ?, ?) WHERE id = ?`, oldDate, newDate, orderID)
	rename(`UPDATE invoices SET invoice_number = replace(invoice_number, ?, ?) WHERE order_id = ?`, oldDate, newDate, orderID)
	rename(`UPDATE payments SET payment_number = replace(payment_number, ?, ?) WHERE order_id = ?`, oldDate, newDate, orderID)
	rename(`UPDATE stock_movements SET reference_number = replace(reference_number, ?, ?) WHERE reference_number IN (?, ?)`, oldDate, newDate, orderNumber, "CANCEL-"+orderNumber)
	rename(`UPDATE notifications SET title = replace(title, ?, ?), message = replace(message, ?, ?) WHERE message LIKE ? OR title LIKE ?`,
		orderNumber, newNumber, orderNumber, newNumber, "%"+orderNumber+"%", "%"+orderNumber+"%")
}

// uploadProof uploads a generated picture of a bank transfer receipt.
func (s *seeder) uploadProof(token string) string {
	img := image.NewRGBA(image.Rect(0, 0, 360, 560))
	fill := func(x0, y0, x1, y1 int, c color.Color) {
		for y := y0; y < y1; y++ {
			for x := x0; x < x1; x++ {
				img.Set(x, y, c)
			}
		}
	}
	fill(0, 0, 360, 560, color.White)
	fill(0, 0, 360, 90, color.RGBA{0, 82, 155, 255})       // bank header
	fill(130, 130, 230, 230, color.RGBA{22, 163, 74, 255}) // success mark
	for i, w := range []int{260, 200, 240, 180, 220, 160} {
		fill(50, 270+i*40, 50+w, 284+i*40, color.RGBA{203, 213, 225, 255})
	}
	var buf bytes.Buffer
	_ = png.Encode(&buf, img)

	var body bytes.Buffer
	w := multipart.NewWriter(&body)
	part, _ := w.CreateFormFile("file", "bukti-transfer.png")
	part.Write(buf.Bytes())
	w.Close()
	req := httptest.NewRequest("POST", "/api/v1/uploads/images", &body)
	req.Header.Set("Content-Type", w.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+token)
	var res struct {
		URL string `json:"url"`
	}
	s.decode(s.serve(req, "POST /uploads/images"), &res)
	return res.URL
}

// must calls the API in-process and stops the program on any error.
func (s *seeder) must(method, path, token string, body any) json.RawMessage {
	var reader *bytes.Reader
	if body != nil {
		b, _ := json.Marshal(body)
		reader = bytes.NewReader(b)
	} else {
		reader = bytes.NewReader(nil)
	}
	req := httptest.NewRequest(method, "/api/v1"+path, reader)
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	return s.serve(req, method+" "+path)
}

func (s *seeder) serve(req *http.Request, label string) json.RawMessage {
	rec := httptest.NewRecorder()
	s.handler.ServeHTTP(rec, req)
	var env struct {
		Success bool            `json:"success"`
		Message string          `json:"message"`
		Data    json.RawMessage `json:"data"`
		Errors  json.RawMessage `json:"errors"`
	}
	_ = json.Unmarshal(rec.Body.Bytes(), &env)
	if rec.Code >= 300 {
		fail("%s gagal (%d): %s %s", label, rec.Code, env.Message, string(env.Errors))
	}
	return env.Data
}

func (s *seeder) decode(raw json.RawMessage, v any) {
	if err := json.Unmarshal(raw, v); err != nil {
		fail("respons tidak terbaca: %v", err)
	}
}

func fail(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "seed-demo: "+format+"\n", args...)
	os.Exit(1)
}
