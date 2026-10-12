package services

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"backend-sembako/internal/constants"
	"backend-sembako/internal/dto"
	"backend-sembako/internal/models"
	"backend-sembako/internal/repositories"
	"backend-sembako/pkg/eventbus"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	ErrOrderNotFound         = errors.New("pesanan tidak ditemukan")
	ErrEmptyCartForOrder     = errors.New("keranjang kosong atau barang yang dipilih tidak ada di keranjang")
	ErrOrderCannotCancel     = errors.New("pesanan tidak bisa dibatalkan pada status ini")
	ErrUnauthorizedOrder     = errors.New("anda tidak berhak melihat atau mengubah pesanan ini")
	ErrInvalidOrderStatus    = errors.New("perubahan status pesanan tidak valid")
	ErrInvalidPaymentMethod  = errors.New("metode pembayaran tidak tersedia")
	ErrCreditNotAllowed      = errors.New("kasbon hanya untuk mitra toko (B2B) yang sudah disetujui")
	ErrIncompleteAddress     = errors.New("alamat pengiriman belum lengkap (nama penerima, nomor HP, alamat, dan kota wajib diisi)")
	ErrInvalidRecipientPhone = errors.New("nomor HP penerima tidak valid (contoh: 081234567890)")
)

var phonePattern = regexp.MustCompile(`^(08|\+628)[0-9]{8,11}$`)

// customerCancellable are the statuses in which a customer may still cancel:
// the store has not started packing yet.
var customerCancellable = map[constants.OrderStatus]bool{
	constants.OrderStatusCreated:        true,
	constants.OrderStatusPendingPayment: true,
	constants.OrderStatusConfirmed:      true,
}

type OrderService interface {
	CreateOrder(ctx context.Context, userID uuid.UUID, userRole constants.UserRole, req dto.CreateOrderDTO) (*dto.OrderResponseDTO, error)
	GetOrders(ctx context.Context, userID uuid.UUID, userRole constants.UserRole, req dto.PaginationRequest, status string) (dto.PaginatedResponse, error)
	GetOrderByID(ctx context.Context, userID uuid.UUID, userRole constants.UserRole, id uuid.UUID) (*dto.OrderResponseDTO, error)
	CancelOrder(ctx context.Context, userID uuid.UUID, userRole constants.UserRole, id uuid.UUID, req dto.CancelOrderDTO) (*dto.OrderResponseDTO, error)
	GetOrderHistory(ctx context.Context, userID uuid.UUID, userRole constants.UserRole, id uuid.UUID) ([]dto.OrderHistoryResponseDTO, error)
	GetInvoiceByID(ctx context.Context, userID uuid.UUID, userRole constants.UserRole, invoiceID uuid.UUID) (*dto.InvoiceResponseDTO, error)
	UpdateOrderStatus(ctx context.Context, userID uuid.UUID, userRole constants.UserRole, id uuid.UUID, req dto.UpdateOrderStatusDTO) (*dto.OrderResponseDTO, error)
}

type orderService struct {
	db                *gorm.DB
	orderRepo         repositories.OrderRepository
	invoiceRepo       repositories.InvoiceRepository
	cartRepo          repositories.CartRepository
	b2bProfileRepo    repositories.B2BProfileRepository
	addressRepo       repositories.AddressRepository
	paymentRepo       repositories.PaymentRepository
	paymentMethodRepo repositories.PaymentMethodRepository
	ledgerRepo        repositories.CreditLedgerRepository
	shippingMethods   repositories.ShippingMethodRepository
	shippingRates     repositories.ShippingRateRepository
	bus               eventbus.EventBus
}

func NewOrderService(
	db *gorm.DB,
	orderRepo repositories.OrderRepository,
	invoiceRepo repositories.InvoiceRepository,
	cartRepo repositories.CartRepository,
	b2bProfileRepo repositories.B2BProfileRepository,
	addressRepo repositories.AddressRepository,
	paymentRepo repositories.PaymentRepository,
	paymentMethodRepo repositories.PaymentMethodRepository,
	ledgerRepo repositories.CreditLedgerRepository,
	shippingMethods repositories.ShippingMethodRepository,
	shippingRates repositories.ShippingRateRepository,
	bus eventbus.EventBus,
) OrderService {
	return &orderService{
		db:                db,
		orderRepo:         orderRepo,
		invoiceRepo:       invoiceRepo,
		cartRepo:          cartRepo,
		b2bProfileRepo:    b2bProfileRepo,
		addressRepo:       addressRepo,
		paymentRepo:       paymentRepo,
		paymentMethodRepo: paymentMethodRepo,
		ledgerRepo:        ledgerRepo,
		shippingMethods:   shippingMethods,
		shippingRates:     shippingRates,
		bus:               bus,
	}
}

type deliveryAddress struct {
	id                                         uuid.UUID
	recipient, phone, full, city, postal, note string
}

func (s *orderService) resolveAddress(ctx context.Context, userID uuid.UUID, req dto.CreateOrderDTO) (*deliveryAddress, error) {
	if req.AddressID != nil && *req.AddressID != uuid.Nil {
		addr, err := s.addressRepo.FindByID(ctx, *req.AddressID)
		if err != nil {
			return nil, err
		}
		if addr == nil || addr.UserID != userID {
			return nil, ErrAddressNotFound
		}
		return &deliveryAddress{
			id: addr.ID, recipient: addr.RecipientName, phone: addr.Phone,
			full: addr.FullAddress, city: addr.City, postal: addr.PostalCode, note: req.Notes,
		}, nil
	}

	a := &deliveryAddress{
		recipient: strings.TrimSpace(req.RecipientName),
		phone:     strings.TrimSpace(req.Phone),
		full:      strings.TrimSpace(req.FullAddress),
		city:      strings.TrimSpace(req.City),
		postal:    strings.TrimSpace(req.PostalCode),
		note:      strings.TrimSpace(req.Notes),
	}
	if len(a.recipient) < 3 || a.phone == "" || a.full == "" || a.city == "" {
		return nil, ErrIncompleteAddress
	}
	if !phonePattern.MatchString(a.phone) {
		return nil, ErrInvalidRecipientPhone
	}
	return a, nil
}

// pickupContact uses the customer's own name and phone for an order they
// collect at the store, so no address has to be filled in.
func (s *orderService) pickupContact(ctx context.Context, userID uuid.UUID, notes string) (*deliveryAddress, error) {
	var u models.User
	if err := s.db.WithContext(ctx).Select("full_name", "phone_number").Where("id = ?", userID).First(&u).Error; err != nil {
		return nil, err
	}
	return &deliveryAddress{
		recipient: u.FullName, phone: u.PhoneNumber,
		full: "Ambil sendiri di toko", city: "-", note: strings.TrimSpace(notes),
	}, nil
}

func (s *orderService) CreateOrder(ctx context.Context, userID uuid.UUID, userRole constants.UserRole, req dto.CreateOrderDTO) (*dto.OrderResponseDTO, error) {
	if userRole.IsStaff() {
		return nil, errors.New("akun staf tidak bisa membuat pesanan; gunakan akun pelanggan")
	}

	cart, err := s.cartRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if cart == nil || len(cart.Items) == 0 {
		return nil, ErrEmptyCartForOrder
	}

	items := cart.Items
	if len(req.CartItemIDs) > 0 {
		wanted := make(map[uuid.UUID]bool, len(req.CartItemIDs))
		for _, id := range req.CartItemIDs {
			wanted[id] = true
		}
		items = items[:0:0]
		for _, it := range cart.Items {
			if wanted[it.ID] {
				items = append(items, it)
				delete(wanted, it.ID)
			}
		}
		if len(wanted) > 0 || len(items) == 0 {
			return nil, ErrEmptyCartForOrder
		}
	}

	method, err := s.paymentMethodRepo.FindByCode(ctx, string(req.PaymentMethod))
	if err != nil {
		return nil, err
	}
	if method == nil || !method.IsActive {
		return nil, ErrInvalidPaymentMethod
	}
	if method.Type == constants.PaymentTypeB2BCredit {
		profile, err := s.b2bProfileRepo.FindByUserID(ctx, userID)
		if err != nil {
			return nil, err
		}
		if userRole != constants.RoleB2BCustomer || profile == nil || profile.VerificationStatus != constants.B2BStatusApproved {
			return nil, ErrCreditNotAllowed
		}
	}

	pickup := req.DeliveryType == constants.DeliveryTypePickup
	var addr *deliveryAddress
	if pickup {
		addr, err = s.pickupContact(ctx, userID, req.Notes)
	} else {
		addr, err = s.resolveAddress(ctx, userID, req)
	}
	if err != nil {
		return nil, err
	}

	shippingFee := 0.0
	deliveryType := constants.DeliveryTypeDelivery
	if pickup {
		deliveryType = constants.DeliveryTypePickup
	} else if _, shippingFee, err = quoteShippingFee(ctx, s.shippingMethods, s.shippingRates, req.ShippingMethodID, addr.city); err != nil {
		return nil, err
	}

	orderNum := newReference("ORD")
	now := time.Now()
	var orderID uuid.UUID
	var grandTotal float64

	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var totalAmount float64
		orderItems := make([]models.OrderItem, 0, len(items))

		for _, item := range items {
			var product models.Product
			if err := tx.Where("id = ?", item.ProductID).First(&product).Error; err != nil || !product.IsActive {
				return fmt.Errorf("produk di keranjang sudah tidak tersedia, hapus dari keranjang lalu coba lagi")
			}
			var unit models.ProductUnit
			if err := tx.Where("id = ? AND product_id = ?", item.ProductUnitID, product.ID).First(&unit).Error; err != nil {
				return fmt.Errorf("satuan untuk '%s' sudah tidak tersedia", product.Name)
			}

			if _, err := applyStockChange(tx, unit.ID, -item.Quantity, constants.MovementStockOut, orderNum,
				fmt.Sprintf("Penjualan pesanan %s", orderNum)); err != nil {
				var short *InsufficientStockError
				if errors.As(err, &short) {
					return fmt.Errorf("stok '%s' tidak cukup (diminta %d, tersedia %d)", product.Name, item.Quantity, short.Available)
				}
				return err
			}

			price := calculateUnitPrice(&unit, userRole, item.Quantity)
			subtotal := float64(item.Quantity) * price
			totalAmount += subtotal

			orderItems = append(orderItems, models.OrderItem{
				ProductID:     product.ID,
				ProductUnitID: unit.ID,
				ProductName:   product.Name,
				SKU:           product.SKU,
				UnitName:      unit.UnitName,
				Quantity:      item.Quantity,
				UnitPrice:     price,
				Subtotal:      subtotal,
				Notes:         item.Notes,
			})
		}

		shippingFee = applyFreeShipping(shippingFee, totalAmount, freeShippingMin(tx))
		grandTotal = totalAmount + shippingFee

		status := constants.OrderStatusConfirmed
		historyNote := "Pesanan diterima. Bayar tunai saat barang diantar."
		if pickup {
			historyNote = "Pesanan diterima. Bayar tunai saat mengambil barang di toko."
		}
		dueDate := now.Add(24 * time.Hour)
		switch method.Type {
		case constants.PaymentTypeManualTransfer:
			status = constants.OrderStatusPendingPayment
			historyNote = "Pesanan dibuat. Menunggu pembayaran transfer dan unggah bukti transfer."
		case constants.PaymentTypeB2BCredit:
			historyNote = "Pesanan dikonfirmasi dengan kasbon tempo 30 hari."
			dueDate = now.Add(30 * 24 * time.Hour)
		}

		order := &models.Order{
			OrderNumber:       orderNum,
			UserID:            userID,
			ShippingAddressID: addr.id,
			Status:            status,
			PaymentMethod:     constants.PaymentMethod(method.Code),
			DeliveryType:      deliveryType,
			TotalAmount:       totalAmount,
			ShippingFee:       shippingFee,
			GrandTotal:        grandTotal,
			DueDate:           &dueDate,
			OrderDate:         now,
			OrderItems:        orderItems,
		}
		if err := s.orderRepo.Create(ctx, tx, order); err != nil {
			return err
		}
		orderID = order.ID

		if err := tx.Create(&models.OrderAddress{
			OrderID:       order.ID,
			RecipientName: addr.recipient,
			Phone:         addr.phone,
			FullAddress:   addr.full,
			City:          addr.city,
			PostalCode:    addr.postal,
			Notes:         addr.note,
		}).Error; err != nil {
			return err
		}

		invoice := &models.Invoice{
			InvoiceNumber: strings.Replace(orderNum, "ORD-", "INV-", 1),
			OrderID:       order.ID,
			UserID:        userID,
			TotalAmount:   grandTotal,
			Status:        constants.InvoiceStatusUnpaid,
			DueDate:       dueDate,
			IssuedAt:      now,
		}
		if err := s.invoiceRepo.Create(ctx, tx, invoice); err != nil {
			return err
		}

		payment := &models.Payment{
			PaymentNumber:   strings.Replace(orderNum, "ORD-", "PAY-", 1),
			OrderID:         order.ID,
			InvoiceID:       invoice.ID,
			UserID:          userID,
			PaymentMethodID: method.ID,
			Amount:          grandTotal,
			GrandTotal:      grandTotal,
			Status:          constants.PaymentStatusPending,
			ExpiredAt:       dueDate,
		}
		if method.Type == constants.PaymentTypeB2BCredit {
			if err := chargeCreditTx(tx, userID, order.ID, grandTotal, orderNum); err != nil {
				return err
			}
			payment.Status = constants.PaymentStatusPaid
			payment.PaidAt = &now
			payment.Notes = "Dibayar dengan kasbon"
		}
		if err := s.paymentRepo.Create(ctx, tx, payment); err != nil {
			return err
		}
		if err := tx.Create(&models.PaymentLog{
			PaymentID: payment.ID, Action: "CREATE_PAYMENT", OldStatus: "NONE", NewStatus: payment.Status,
			Notes: "Tagihan dibuat bersama pesanan " + orderNum,
		}).Error; err != nil {
			return err
		}
		if payment.Status == constants.PaymentStatusPaid {
			if err := markInvoicePaidTx(tx, invoice.ID, now); err != nil {
				return err
			}
		}

		if err := s.orderRepo.AddHistory(ctx, tx, &models.OrderHistory{OrderID: order.ID, Status: status, Notes: historyNote}); err != nil {
			return err
		}

		ids := make([]uuid.UUID, len(items))
		for i, it := range items {
			ids[i] = it.ID
		}
		return tx.Where("id IN ?", ids).Delete(&models.CartItem{}).Error
	})
	if err != nil {
		return nil, err
	}

	created, err := s.GetOrderByID(ctx, userID, userRole, orderID)
	if err != nil {
		return nil, err
	}
	customerName := "Pelanggan"
	if created.Customer != nil {
		customerName = created.Customer.FullName
	}

	publishEvent(s.bus, "ORDER_CREATED", "ORDER", orderID, map[string]interface{}{
		"user_id":       userID.String(),
		"customer_name": customerName,
		"order_number":  orderNum,
		"grand_total":   formatRupiah(grandTotal),
		"title":         "Pesanan " + orderNum + " berhasil dibuat",
		"message":       "Total " + formatRupiah(grandTotal) + ". Kami akan segera memproses pesanan Anda.",
		"notify_staff":  true,
		"staff_title":   "Pesanan baru " + orderNum,
		"staff_message": fmt.Sprintf("Pesanan baru dari %s senilai %s dengan pembayaran %s.", customerName, formatRupiah(grandTotal), method.Name),
	})

	return created, nil
}

// chargeCreditTx deducts an order total from the customer's kasbon limit.
// The ledger row is locked so two simultaneous orders cannot both pass the
// limit check.
func chargeCreditTx(tx *gorm.DB, userID, orderID uuid.UUID, amount float64, reference string) error {
	var ledger models.CreditLedger
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("user_id = ?", userID).First(&ledger).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("akun kasbon belum dibuat oleh toko; hubungi admin")
		}
		return err
	}
	remaining := ledger.CreditLimit - ledger.UsedCredit
	if amount > remaining {
		return fmt.Errorf("sisa limit kasbon tidak cukup (dibutuhkan %s, sisa %s)", formatRupiah(amount), formatRupiah(remaining))
	}

	ledger.UsedCredit += amount
	ledger.OutstandingBalance += amount
	ledger.RemainingCredit = ledger.CreditLimit - ledger.UsedCredit
	if err := tx.Model(&models.CreditLedger{}).Where("id = ?", ledger.ID).Updates(map[string]interface{}{
		"used_credit":         ledger.UsedCredit,
		"outstanding_balance": ledger.OutstandingBalance,
		"remaining_credit":    ledger.RemainingCredit,
	}).Error; err != nil {
		return err
	}

	return tx.Create(&models.CreditPayment{
		LedgerID:         ledger.ID,
		OrderID:          &orderID,
		Type:             string(constants.CreditTxUsage),
		Amount:           amount,
		PaymentDate:      time.Now(),
		RemainingBalance: ledger.OutstandingBalance,
		Notes:            "Pemakaian kasbon untuk pesanan " + reference,
	}).Error
}

// reverseCreditTx gives kasbon back when a kasbon order is cancelled.
func reverseCreditTx(tx *gorm.DB, userID, orderID uuid.UUID, amount float64, reference string) error {
	var ledger models.CreditLedger
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("user_id = ?", userID).First(&ledger).Error; err != nil {
		return err
	}
	ledger.UsedCredit -= amount
	if ledger.UsedCredit < 0 {
		ledger.UsedCredit = 0
	}
	ledger.OutstandingBalance -= amount
	if ledger.OutstandingBalance < 0 {
		ledger.OutstandingBalance = 0
	}
	ledger.RemainingCredit = ledger.CreditLimit - ledger.UsedCredit
	if ledger.RemainingCredit < 0 {
		ledger.RemainingCredit = 0
	}
	if err := tx.Model(&models.CreditLedger{}).Where("id = ?", ledger.ID).Updates(map[string]interface{}{
		"used_credit":         ledger.UsedCredit,
		"outstanding_balance": ledger.OutstandingBalance,
		"remaining_credit":    ledger.RemainingCredit,
	}).Error; err != nil {
		return err
	}
	return tx.Create(&models.CreditPayment{
		LedgerID:         ledger.ID,
		OrderID:          &orderID,
		Type:             string(constants.CreditTxReversal),
		Amount:           amount,
		PaymentDate:      time.Now(),
		RemainingBalance: ledger.OutstandingBalance,
		Notes:            "Pengembalian kasbon karena pesanan " + reference + " dibatalkan",
	}).Error
}

func (s *orderService) GetOrders(ctx context.Context, userID uuid.UUID, userRole constants.UserRole, req dto.PaginationRequest, status string) (dto.PaginatedResponse, error) {
	var orders []models.Order
	var total int64
	var err error

	if userRole.IsStaff() {
		orders, total, err = s.orderRepo.FindAll(ctx, req, status)
	} else {
		orders, total, err = s.orderRepo.FindByUserID(ctx, userID, req, status)
	}
	if err != nil {
		return dto.PaginatedResponse{}, err
	}

	dtos := make([]dto.OrderResponseDTO, len(orders))
	for i := range orders {
		dtos[i] = mapOrderToDTO(&orders[i], userRole)
	}

	req.SetDefaults()
	return dto.NewPaginatedResponse(dtos, req.Page, req.Limit, total), nil
}

func (s *orderService) loadOrderFor(ctx context.Context, userID uuid.UUID, userRole constants.UserRole, id uuid.UUID) (*models.Order, error) {
	order, err := s.orderRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, ErrOrderNotFound
	}
	if !userRole.IsStaff() && order.UserID != userID {
		return nil, ErrOrderNotFound
	}
	return order, nil
}

func (s *orderService) GetOrderByID(ctx context.Context, userID uuid.UUID, userRole constants.UserRole, id uuid.UUID) (*dto.OrderResponseDTO, error) {
	order, err := s.loadOrderFor(ctx, userID, userRole, id)
	if err != nil {
		return nil, err
	}

	res := mapOrderToDTO(order, userRole)
	payments, err := s.paymentRepo.FindByOrderID(ctx, order.ID)
	if err != nil {
		return nil, err
	}
	res.Payments = make([]dto.PaymentResponseDTO, len(payments))
	for i := range payments {
		res.Payments[i] = mapPaymentToDTO(&payments[i])
	}
	return &res, nil
}

func (s *orderService) CancelOrder(ctx context.Context, userID uuid.UUID, userRole constants.UserRole, id uuid.UUID, req dto.CancelOrderDTO) (*dto.OrderResponseDTO, error) {
	order, err := s.loadOrderFor(ctx, userID, userRole, id)
	if err != nil {
		return nil, err
	}

	allowed := customerCancellable[order.Status]
	if userRole.IsStaff() {
		allowed = order.Status.CanTransitionTo(constants.OrderStatusCancelled)
	}
	if !allowed {
		return nil, ErrOrderCannotCancel
	}

	by := "pelanggan"
	if userRole.IsStaff() {
		by = "toko"
	}
	reason := strings.TrimSpace(req.Reason)
	refundNote := ""

	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, item := range order.OrderItems {
			if _, err := applyStockChange(tx, item.ProductUnitID, item.Quantity, constants.MovementStockIn,
				"CANCEL-"+order.OrderNumber, "Pengembalian stok pembatalan pesanan "+order.OrderNumber); err != nil {
				return err
			}
		}

		var payments []models.Payment
		if err := tx.Preload("PaymentMethod").Where("order_id = ?", order.ID).Find(&payments).Error; err != nil {
			return err
		}
		for _, p := range payments {
			newStatus := ""
			switch {
			case p.Status == constants.PaymentStatusPending || p.Status == constants.PaymentStatusWaitingConfirmation:
				newStatus = constants.PaymentStatusCancelled
			case p.Status == constants.PaymentStatusPaid && p.PaymentMethod != nil && p.PaymentMethod.Type == constants.PaymentTypeB2BCredit:
				if err := reverseCreditTx(tx, order.UserID, order.ID, p.GrandTotal, order.OrderNumber); err != nil {
					return err
				}
				newStatus = constants.PaymentStatusRefunded
			case p.Status == constants.PaymentStatusPaid:
				refundNote = " Pembayaran sudah diterima, lakukan refund ke pelanggan dari menu Pembayaran."
			}
			if newStatus == "" {
				continue
			}
			if err := tx.Model(&models.Payment{}).Where("id = ?", p.ID).Update("status", newStatus).Error; err != nil {
				return err
			}
			if err := tx.Create(&models.PaymentLog{
				PaymentID: p.ID, Action: "ORDER_CANCELLED", OldStatus: p.Status, NewStatus: newStatus,
				Notes: "Pesanan dibatalkan: " + reason,
			}).Error; err != nil {
				return err
			}
		}

		if order.Invoice != nil && order.Invoice.Status != constants.InvoiceStatusPaid {
			if err := s.invoiceRepo.UpdateStatus(ctx, tx, order.Invoice.ID, constants.InvoiceStatusCancelled); err != nil {
				return err
			}
		}

		if err := tx.Model(&models.Shipment{}).
			Where("order_id = ? AND status NOT IN ?", order.ID, []string{"DELIVERED", "RETURNED"}).
			Update("status", "RETURNED").Error; err != nil {
			return err
		}

		return setOrderStatusTx(tx, order, constants.OrderStatusCancelled,
			fmt.Sprintf("Pesanan dibatalkan oleh %s: %s.%s", by, reason, refundNote))
	})
	if err != nil {
		return nil, fmt.Errorf("gagal membatalkan pesanan: %w", err)
	}

	publishEvent(s.bus, "ORDER_CANCELLED", "ORDER", order.ID, map[string]interface{}{
		"user_id":       order.UserID.String(),
		"order_number":  order.OrderNumber,
		"title":         "Pesanan " + order.OrderNumber + " dibatalkan",
		"message":       "Alasan: " + reason,
		"notify_staff":  !userRole.IsStaff(),
		"staff_title":   "Pesanan " + order.OrderNumber + " dibatalkan pelanggan",
		"staff_message": "Alasan: " + reason,
	})

	return s.GetOrderByID(ctx, userID, userRole, id)
}

func (s *orderService) GetOrderHistory(ctx context.Context, userID uuid.UUID, userRole constants.UserRole, id uuid.UUID) ([]dto.OrderHistoryResponseDTO, error) {
	order, err := s.GetOrderByID(ctx, userID, userRole, id)
	if err != nil {
		return nil, err
	}
	return order.Histories, nil
}

func (s *orderService) GetInvoiceByID(ctx context.Context, userID uuid.UUID, userRole constants.UserRole, invoiceID uuid.UUID) (*dto.InvoiceResponseDTO, error) {
	inv, err := s.invoiceRepo.FindByID(ctx, invoiceID)
	if err != nil {
		return nil, err
	}
	if inv == nil || (!userRole.IsStaff() && inv.UserID != userID) {
		return nil, errors.New("invoice tidak ditemukan")
	}

	return &dto.InvoiceResponseDTO{
		ID:            inv.ID,
		InvoiceNumber: inv.InvoiceNumber,
		OrderID:       inv.OrderID,
		TotalAmount:   inv.TotalAmount,
		Status:        inv.Status,
		DueDate:       inv.DueDate,
		IssuedAt:      inv.IssuedAt,
		PaidAt:        inv.PaidAt,
	}, nil
}

var orderStatusMessages = map[constants.OrderStatus]string{
	constants.OrderStatusProcessingPacking: "Pesanan Anda sedang disiapkan.",
	constants.OrderStatusOutForDelivery:    "Pesanan Anda sedang diantar ke alamat Anda.",
	constants.OrderStatusCompleted:         "Pesanan sudah sampai. Terima kasih sudah berbelanja!",
}

var pickupStatusMessages = map[constants.OrderStatus]string{
	constants.OrderStatusProcessingPacking: "Pesanan Anda sedang disiapkan.",
	constants.OrderStatusOutForDelivery:    "Pesanan Anda sudah siap. Silakan ambil di toko.",
	constants.OrderStatusCompleted:         "Pesanan sudah diambil. Terima kasih sudah berbelanja!",
}

func statusMessage(o *models.Order, status constants.OrderStatus) string {
	if o.DeliveryType == constants.DeliveryTypePickup {
		return pickupStatusMessages[status]
	}
	return orderStatusMessages[status]
}

func (s *orderService) UpdateOrderStatus(ctx context.Context, userID uuid.UUID, userRole constants.UserRole, id uuid.UUID, req dto.UpdateOrderStatusDTO) (*dto.OrderResponseDTO, error) {
	if !userRole.IsStaff() {
		return nil, ErrUnauthorizedOrder
	}
	if req.Status == constants.OrderStatusCancelled {
		reason := strings.TrimSpace(req.Notes)
		if len(reason) < 3 {
			reason = "Dibatalkan oleh toko"
		}
		return s.CancelOrder(ctx, userID, userRole, id, dto.CancelOrderDTO{Reason: reason})
	}

	order, err := s.loadOrderFor(ctx, userID, userRole, id)
	if err != nil {
		return nil, err
	}
	if !order.Status.CanTransitionTo(req.Status) {
		return nil, fmt.Errorf("%w: dari %s ke %s", ErrInvalidOrderStatus, order.Status, req.Status)
	}
	if userRole == constants.RoleDriver && req.Status != constants.OrderStatusCompleted {
		return nil, errors.New("kurir hanya bisa menandai pesanan selesai")
	}

	notes := strings.TrimSpace(req.Notes)
	if notes == "" {
		notes = statusMessage(order, req.Status)
	}

	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		switch req.Status {
		case constants.OrderStatusCompleted:
			return completeOrderTx(tx, order, notes)
		case constants.OrderStatusOutForDelivery:
			if err := setOrderStatusTx(tx, order, req.Status, notes); err != nil {
				return err
			}
			if order.DeliveryType == constants.DeliveryTypePickup {
				return nil // ready at the store; no courier trip to track
			}
			shipment, err := ensureShipmentTx(tx, order, "OUT_FOR_DELIVERY", "Toko", "Pesanan diantar kurir toko")
			if err != nil {
				return err
			}
			if shipment.Status != "OUT_FOR_DELIVERY" {
				now := time.Now()
				if err := tx.Model(&models.Shipment{}).Where("id = ?", shipment.ID).
					Updates(map[string]interface{}{"status": "OUT_FOR_DELIVERY", "shipped_at": now}).Error; err != nil {
					return err
				}
				return tx.Create(&models.ShipmentTracking{
					ShipmentID: shipment.ID, Status: "OUT_FOR_DELIVERY", Location: "Toko", Notes: "Pesanan diantar kurir",
				}).Error
			}
			return nil
		default:
			return setOrderStatusTx(tx, order, req.Status, notes)
		}
	})
	if err != nil {
		return nil, fmt.Errorf("gagal memperbarui status pesanan: %w", err)
	}

	publishEvent(s.bus, "ORDER_STATUS_CHANGED", "ORDER", order.ID, map[string]interface{}{
		"user_id":      order.UserID.String(),
		"order_number": order.OrderNumber,
		"status":       string(req.Status),
		"title":        "Update pesanan " + order.OrderNumber,
		"message":      notes,
	})

	return s.GetOrderByID(ctx, userID, userRole, id)
}

func mapOrderToDTO(o *models.Order, viewer constants.UserRole) dto.OrderResponseDTO {
	res := dto.OrderResponseDTO{
		ID:            o.ID,
		OrderNumber:   o.OrderNumber,
		UserID:        o.UserID,
		Status:        o.Status,
		NextStatuses:  []constants.OrderStatus{},
		PaymentMethod: o.PaymentMethod,
		DeliveryType:  o.DeliveryType,
		TotalAmount:   o.TotalAmount,
		ShippingFee:   o.ShippingFee,
		GrandTotal:    o.GrandTotal,
		DueDate:       o.DueDate,
		OrderDate:     o.OrderDate,
	}

	if viewer.IsStaff() {
		if next := o.Status.NextStatuses(); next != nil {
			res.NextStatuses = next
		}
		res.CanCancel = o.Status.CanTransitionTo(constants.OrderStatusCancelled)
	} else {
		res.CanCancel = customerCancellable[o.Status]
	}

	if o.User != nil {
		res.Customer = &dto.OrderCustomerDTO{
			ID: o.User.ID, FullName: o.User.FullName, PhoneNumber: o.User.PhoneNumber, Role: o.User.Role,
		}
	}

	if o.OrderAddress != nil {
		res.Address = &dto.OrderAddressResponseDTO{
			RecipientName: o.OrderAddress.RecipientName,
			Phone:         o.OrderAddress.Phone,
			FullAddress:   o.OrderAddress.FullAddress,
			City:          o.OrderAddress.City,
			PostalCode:    o.OrderAddress.PostalCode,
			Notes:         o.OrderAddress.Notes,
		}
	}

	if o.Invoice != nil {
		res.Invoice = &dto.InvoiceResponseDTO{
			ID:            o.Invoice.ID,
			InvoiceNumber: o.Invoice.InvoiceNumber,
			OrderID:       o.Invoice.OrderID,
			TotalAmount:   o.Invoice.TotalAmount,
			Status:        o.Invoice.Status,
			DueDate:       o.Invoice.DueDate,
			IssuedAt:      o.Invoice.IssuedAt,
			PaidAt:        o.Invoice.PaidAt,
		}
	}

	res.Items = make([]dto.OrderItemResponseDTO, len(o.OrderItems))
	for i, item := range o.OrderItems {
		name, sku, unitName, image := item.ProductName, item.SKU, item.UnitName, ""
		if item.Product != nil {
			if name == "" {
				name = item.Product.Name
			}
			if sku == "" {
				sku = item.Product.SKU
			}
		}
		if unitName == "" && item.ProductUnit != nil {
			unitName = item.ProductUnit.UnitName
		}
		if name == "" {
			name = "Produk"
		}

		res.Items[i] = dto.OrderItemResponseDTO{
			ID:            item.ID,
			ProductID:     item.ProductID,
			ProductName:   name,
			SKU:           sku,
			ImageURL:      image,
			ProductUnitID: item.ProductUnitID,
			UnitName:      unitName,
			Quantity:      item.Quantity,
			UnitPrice:     item.UnitPrice,
			Subtotal:      item.Subtotal,
			Notes:         item.Notes,
		}
		res.TotalItems += item.Quantity
	}

	res.Histories = make([]dto.OrderHistoryResponseDTO, len(o.Histories))
	for i, h := range o.Histories {
		res.Histories[i] = dto.OrderHistoryResponseDTO{ID: h.ID, Status: h.Status, Notes: h.Notes, CreatedAt: h.CreatedAt}
	}

	return res
}
