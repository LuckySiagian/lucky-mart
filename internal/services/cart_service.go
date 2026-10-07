package services

import (
	"context"
	"errors"
	"fmt"

	"backend-sembako/internal/constants"
	"backend-sembako/internal/dto"
	"backend-sembako/internal/models"
	"backend-sembako/internal/repositories"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrCartItemNotFound   = errors.New("cart item not found")
	ErrInsufficientStock  = errors.New("requested quantity exceeds available stock")
	ErrInactiveProduct    = errors.New("cannot add inactive or deleted product to cart")
	ErrUnauthorizedCart   = errors.New("you do not have permission to modify this cart item")
)

type CartService interface {
	GetCart(ctx context.Context, userID uuid.UUID, userRole constants.UserRole) (*dto.CartResponseDTO, error)
	AddToCart(ctx context.Context, userID uuid.UUID, userRole constants.UserRole, req dto.AddToCartDTO) (*dto.CartResponseDTO, error)
	UpdateCartItem(ctx context.Context, userID uuid.UUID, userRole constants.UserRole, itemID uuid.UUID, req dto.UpdateCartItemDTO) (*dto.CartResponseDTO, error)
	RemoveCartItem(ctx context.Context, userID uuid.UUID, userRole constants.UserRole, itemID uuid.UUID) (*dto.CartResponseDTO, error)
	ClearCart(ctx context.Context, userID uuid.UUID) error
	GetCheckoutSummary(ctx context.Context, userID uuid.UUID, userRole constants.UserRole) (*dto.CheckoutSummaryDTO, error)
}

type cartService struct {
	db           *gorm.DB
	cartRepo     repositories.CartRepository
	cartItemRepo repositories.CartItemRepository
	productRepo  repositories.ProductRepository
}

func NewCartService(
	db *gorm.DB,
	cartRepo repositories.CartRepository,
	cartItemRepo repositories.CartItemRepository,
	productRepo repositories.ProductRepository,
) CartService {
	return &cartService{
		db:           db,
		cartRepo:     cartRepo,
		cartItemRepo: cartItemRepo,
		productRepo:  productRepo,
	}
}

func (s *cartService) GetCart(ctx context.Context, userID uuid.UUID, userRole constants.UserRole) (*dto.CartResponseDTO, error) {
	cart, err := s.cartRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if cart == nil {
		return &dto.CartResponseDTO{
			UserID:      userID,
			TotalItems:  0,
			TotalAmount: 0.00,
			Items:       []dto.CartItemResponseDTO{},
		}, nil
	}

	return s.buildCartResponse(ctx, cart, userRole)
}

func (s *cartService) AddToCart(ctx context.Context, userID uuid.UUID, userRole constants.UserRole, req dto.AddToCartDTO) (*dto.CartResponseDTO, error) {
	// Validate Product & Active status
	product, err := s.productRepo.FindByID(ctx, req.ProductID)
	if err != nil || product == nil || !product.IsActive {
		return nil, ErrInactiveProduct
	}

	// Find ProductUnit
	var targetUnit *models.ProductUnit
	for _, u := range product.ProductUnits {
		if u.ID == req.ProductUnitID {
			targetUnit = &u
			break
		}
	}
	if targetUnit == nil {
		return nil, errors.New("specified product unit not found")
	}

	// Automatic Cart Creation
	cart, err := s.cartRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if cart == nil {
		cart = &models.Cart{
			UserID:      userID,
			TotalItems:  0,
			TotalAmount: 0.00,
		}
		if err := s.cartRepo.Create(ctx, cart); err != nil {
			return nil, err
		}
	}

	// Check if item already exists in cart
	existingItem, err := s.cartItemRepo.FindByCartAndUnit(ctx, cart.ID, req.ProductUnitID)
	if err != nil {
		return nil, err
	}

	targetQty := req.Quantity
	if existingItem != nil {
		targetQty += existingItem.Quantity
	}

	// Validate Available Stock
	if targetQty > targetUnit.StockQuantity {
		return nil, fmt.Errorf("%w: requested %d, available %d", ErrInsufficientStock, targetQty, targetUnit.StockQuantity)
	}

	unitPrice := calculateUnitPrice(targetUnit, userRole, targetQty)
	subtotal := float64(targetQty) * unitPrice

	if existingItem != nil {
		existingItem.Quantity = targetQty
		existingItem.UnitPrice = unitPrice
		existingItem.Subtotal = subtotal
		existingItem.Notes = req.Notes
		if err := s.cartItemRepo.Update(ctx, existingItem); err != nil {
			return nil, err
		}
	} else {
		newItem := &models.CartItem{
			CartID:        cart.ID,
			ProductID:     req.ProductID,
			ProductUnitID: req.ProductUnitID,
			Quantity:      targetQty,
			UnitPrice:     unitPrice,
			Subtotal:      subtotal,
			Notes:         req.Notes,
		}
		if err := s.cartItemRepo.Create(ctx, newItem); err != nil {
			return nil, err
		}
	}

	return s.GetCart(ctx, userID, userRole)
}

func (s *cartService) UpdateCartItem(ctx context.Context, userID uuid.UUID, userRole constants.UserRole, itemID uuid.UUID, req dto.UpdateCartItemDTO) (*dto.CartResponseDTO, error) {
	item, err := s.cartItemRepo.FindByID(ctx, itemID)
	if err != nil || item == nil {
		return nil, ErrCartItemNotFound
	}

	cart, err := s.cartRepo.FindByUserID(ctx, userID)
	if err != nil || cart == nil || item.CartID != cart.ID {
		return nil, ErrUnauthorizedCart
	}

	product, err := s.productRepo.FindByID(ctx, item.ProductID)
	if err != nil || product == nil || !product.IsActive {
		return nil, ErrInactiveProduct
	}

	var targetUnit *models.ProductUnit
	for _, u := range product.ProductUnits {
		if u.ID == item.ProductUnitID {
			targetUnit = &u
			break
		}
	}
	if targetUnit == nil {
		return nil, errors.New("product unit not found")
	}

	if req.Quantity > targetUnit.StockQuantity {
		return nil, fmt.Errorf("%w: requested %d, available %d", ErrInsufficientStock, req.Quantity, targetUnit.StockQuantity)
	}

	unitPrice := calculateUnitPrice(targetUnit, userRole, req.Quantity)
	item.Quantity = req.Quantity
	item.UnitPrice = unitPrice
	item.Subtotal = float64(req.Quantity) * unitPrice
	item.Notes = req.Notes

	if err := s.cartItemRepo.Update(ctx, item); err != nil {
		return nil, err
	}

	return s.GetCart(ctx, userID, userRole)
}

func (s *cartService) RemoveCartItem(ctx context.Context, userID uuid.UUID, userRole constants.UserRole, itemID uuid.UUID) (*dto.CartResponseDTO, error) {
	item, err := s.cartItemRepo.FindByID(ctx, itemID)
	if err != nil || item == nil {
		return nil, ErrCartItemNotFound
	}

	cart, err := s.cartRepo.FindByUserID(ctx, userID)
	if err != nil || cart == nil || item.CartID != cart.ID {
		return nil, ErrUnauthorizedCart
	}

	if err := s.cartItemRepo.Delete(ctx, itemID); err != nil {
		return nil, err
	}

	return s.GetCart(ctx, userID, userRole)
}

func (s *cartService) ClearCart(ctx context.Context, userID uuid.UUID) error {
	cart, err := s.cartRepo.FindByUserID(ctx, userID)
	if err != nil || cart == nil {
		return nil
	}
	return s.cartRepo.Clear(ctx, cart.ID)
}

func (s *cartService) GetCheckoutSummary(ctx context.Context, userID uuid.UUID, userRole constants.UserRole) (*dto.CheckoutSummaryDTO, error) {
	cartRes, err := s.GetCart(ctx, userID, userRole)
	if err != nil {
		return nil, err
	}

	var validationErrors []string
	canProceed := true

	if len(cartRes.Items) == 0 {
		canProceed = false
		validationErrors = append(validationErrors, "Your shopping cart is empty")
	}

	for _, item := range cartRes.Items {
		if !item.IsProductActive {
			canProceed = false
			validationErrors = append(validationErrors, fmt.Sprintf("Product '%s' is no longer available", item.ProductName))
		}
		if !item.IsStockSufficient {
			canProceed = false
			validationErrors = append(validationErrors, fmt.Sprintf("Item '%s' exceeds available stock (In cart: %d, Available: %d)", item.ProductName, item.Quantity, item.AvailableStock))
		}
	}

	estimatedShipping := 15000.00 // Default flat rate shipping for MVP
	grandTotal := cartRes.TotalAmount + estimatedShipping

	return &dto.CheckoutSummaryDTO{
		CartID:               cartRes.ID,
		UserID:               userID,
		TotalItems:           cartRes.TotalItems,
		SubtotalAmount:       cartRes.TotalAmount,
		EstimatedShippingFee: estimatedShipping,
		GrandTotalAmount:     grandTotal,
		Items:                cartRes.Items,
		CanProceedToOrder:    canProceed,
		ValidationErrors:     validationErrors,
	}, nil
}

func (s *cartService) buildCartResponse(ctx context.Context, cart *models.Cart, userRole constants.UserRole) (*dto.CartResponseDTO, error) {
	var itemDTOs []dto.CartItemResponseDTO
	var totalAmount float64
	var totalItems int

	for _, item := range cart.Items {
		product, _ := s.productRepo.FindByID(ctx, item.ProductID)
		isActive := product != nil && product.IsActive

		var targetUnit *models.ProductUnit
		if product != nil {
			for _, u := range product.ProductUnits {
				if u.ID == item.ProductUnitID {
					targetUnit = &u
					break
				}
			}
		}

		availableStock := 0
		if targetUnit != nil {
			availableStock = targetUnit.StockQuantity
		}

		unitPrice := item.UnitPrice
		if targetUnit != nil {
			unitPrice = calculateUnitPrice(targetUnit, userRole, item.Quantity)
		}
		subtotal := float64(item.Quantity) * unitPrice

		isStockSufficient := item.Quantity <= availableStock

		itemDTOs = append(itemDTOs, dto.CartItemResponseDTO{
			ID:                 item.ID,
			ProductID:          item.ProductID,
			ProductName:        getItemProductName(product),
			SKU:                getItemProductSKU(product),
			ProductUnitID:      item.ProductUnitID,
			UnitName:           getItemUnitName(targetUnit),
			Quantity:           item.Quantity,
			UnitPrice:          unitPrice,
			Subtotal:           subtotal,
			Notes:              item.Notes,
			AvailableStock:     availableStock,
			IsStockSufficient:  isStockSufficient,
			IsProductActive:    isActive,
		})

		if isActive && isStockSufficient {
			totalAmount += subtotal
			totalItems += item.Quantity
		}
	}

	return &dto.CartResponseDTO{
		ID:          cart.ID,
		UserID:      cart.UserID,
		TotalItems:  totalItems,
		TotalAmount: totalAmount,
		Items:       itemDTOs,
		CreatedAt:   cart.CreatedAt,
	}, nil
}

func calculateUnitPrice(unit *models.ProductUnit, role constants.UserRole, qty int) float64 {
	if role == constants.RoleB2BCustomer && qty >= unit.MinWholesaleQty {
		return unit.WholesalePrice
	}
	return unit.RetailPrice
}

func getItemProductName(p *models.Product) string {
	if p != nil {
		return p.Name
	}
	return "Unavailable Product"
}

func getItemProductSKU(p *models.Product) string {
	if p != nil {
		return p.SKU
	}
	return ""
}

func getItemUnitName(u *models.ProductUnit) string {
	if u != nil {
		return u.UnitName
	}
	return "Unit"
}
