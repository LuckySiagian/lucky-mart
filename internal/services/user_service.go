package services

import (
	"context"
	"errors"
	"strings"

	"backend-sembako/helpers"
	"backend-sembako/internal/constants"
	"backend-sembako/internal/dto"
	"backend-sembako/internal/models"
	"backend-sembako/internal/repositories"
	"backend-sembako/pkg/eventbus"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrUserNotFound      = errors.New("pengguna tidak ditemukan")
	ErrForbiddenUserEdit = errors.New("anda tidak berhak mengubah akun ini")
	ErrLastSuperAdmin    = errors.New("tidak bisa menonaktifkan atau menghapus Super Admin terakhir")
	ErrSelfModification  = errors.New("anda tidak bisa menonaktifkan, menghapus, atau mengubah peran akun sendiri")
	ErrInvalidRole       = errors.New("peran pengguna tidak valid")
)

type UserService interface {
	GetUsers(ctx context.Context, req dto.PaginationRequest, role string) (dto.PaginatedResponse, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*dto.UserProfileResponse, error)
	CreateUser(ctx context.Context, actorRole constants.UserRole, req dto.CreateUserDTO) (*dto.UserProfileResponse, error)
	UpdateUser(ctx context.Context, actorID uuid.UUID, actorRole constants.UserRole, id uuid.UUID, req dto.UpdateUserDTO) (*dto.UserProfileResponse, error)
	ResetPassword(ctx context.Context, actorRole constants.UserRole, id uuid.UUID, newPassword string) error
	VerifyB2B(ctx context.Context, id uuid.UUID, req dto.VerifyB2BDTO) (*dto.UserProfileResponse, error)
	DeleteUser(ctx context.Context, actorID uuid.UUID, actorRole constants.UserRole, id uuid.UUID) error
}

type userService struct {
	db               *gorm.DB
	userRepo         repositories.UserRepository
	b2bProfileRepo   repositories.B2BProfileRepository
	refreshTokenRepo repositories.RefreshTokenRepository
	bus              eventbus.EventBus
}

func NewUserService(
	db *gorm.DB,
	userRepo repositories.UserRepository,
	b2bProfileRepo repositories.B2BProfileRepository,
	refreshTokenRepo repositories.RefreshTokenRepository,
	bus eventbus.EventBus,
) UserService {
	return &userService{
		db:               db,
		userRepo:         userRepo,
		b2bProfileRepo:   b2bProfileRepo,
		refreshTokenRepo: refreshTokenRepo,
		bus:              bus,
	}
}

// canManage reports whether actor may create or change an account with the
// given role. Only a Super Admin may touch admin-level accounts.
func canManage(actor, target constants.UserRole) bool {
	if actor == constants.RoleSuperAdmin {
		return true
	}
	if actor == constants.RoleStoreManager {
		return target != constants.RoleSuperAdmin && target != constants.RoleStoreManager
	}
	return false
}

func (s *userService) GetUsers(ctx context.Context, req dto.PaginationRequest, role string) (dto.PaginatedResponse, error) {
	users, total, err := s.userRepo.FindAll(ctx, req, role)
	if err != nil {
		return dto.PaginatedResponse{}, err
	}

	dtos := make([]dto.UserProfileResponse, len(users))
	for i := range users {
		dtos[i] = mapUserToProfile(&users[i])
	}

	req.SetDefaults()
	return dto.NewPaginatedResponse(dtos, req.Page, req.Limit, total), nil
}

func (s *userService) GetUserByID(ctx context.Context, id uuid.UUID) (*dto.UserProfileResponse, error) {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}
	res := mapUserToProfile(user)
	return &res, nil
}

func (s *userService) CreateUser(ctx context.Context, actorRole constants.UserRole, req dto.CreateUserDTO) (*dto.UserProfileResponse, error) {
	if !req.Role.IsValid() {
		return nil, ErrInvalidRole
	}
	if !canManage(actorRole, req.Role) {
		return nil, ErrForbiddenUserEdit
	}
	if err := helpers.ValidatePasswordStrength(req.Password); err != nil {
		return nil, err
	}

	existing, err := s.userRepo.FindByPhoneNumber(ctx, req.PhoneNumber)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, ErrPhoneAlreadyExists
	}
	if req.Email != nil && *req.Email != "" {
		if other, err := s.userRepo.FindByEmail(ctx, *req.Email); err != nil {
			return nil, err
		} else if other != nil {
			return nil, ErrEmailAlreadyExists
		}
	} else {
		req.Email = nil
	}

	hashed, err := helpers.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		FullName:     strings.TrimSpace(req.FullName),
		Email:        req.Email,
		PhoneNumber:  req.PhoneNumber,
		PasswordHash: hashed,
		Role:         req.Role,
		IsActive:     true,
	}
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return s.GetUserByID(ctx, user.ID)
}

func (s *userService) UpdateUser(ctx context.Context, actorID uuid.UUID, actorRole constants.UserRole, id uuid.UUID, req dto.UpdateUserDTO) (*dto.UserProfileResponse, error) {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}
	if !canManage(actorRole, user.Role) {
		return nil, ErrForbiddenUserEdit
	}

	if req.Role != nil && *req.Role != user.Role {
		if !req.Role.IsValid() {
			return nil, ErrInvalidRole
		}
		if !canManage(actorRole, *req.Role) {
			return nil, ErrForbiddenUserEdit
		}
		if id == actorID {
			return nil, ErrSelfModification
		}
		if user.Role == constants.RoleSuperAdmin {
			if err := s.ensureAnotherSuperAdmin(ctx); err != nil {
				return nil, err
			}
		}
		user.Role = *req.Role
	}
	if req.IsActive != nil && *req.IsActive != user.IsActive {
		if id == actorID {
			return nil, ErrSelfModification
		}
		if !*req.IsActive && user.Role == constants.RoleSuperAdmin {
			if err := s.ensureAnotherSuperAdmin(ctx); err != nil {
				return nil, err
			}
		}
		user.IsActive = *req.IsActive
	}
	if req.FullName != nil {
		user.FullName = strings.TrimSpace(*req.FullName)
	}
	if req.Email != nil {
		email := strings.TrimSpace(*req.Email)
		if email == "" {
			user.Email = nil
		} else {
			if other, err := s.userRepo.FindByEmail(ctx, email); err != nil {
				return nil, err
			} else if other != nil && other.ID != user.ID {
				return nil, ErrEmailAlreadyExists
			}
			user.Email = &email
		}
	}

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}
	if !user.IsActive || req.Role != nil {
		// Force a fresh login so the new role / deactivation takes effect.
		_ = s.refreshTokenRepo.RevokeByUserID(ctx, user.ID)
	}

	return s.GetUserByID(ctx, id)
}

func (s *userService) ensureAnotherSuperAdmin(ctx context.Context) error {
	count, err := s.userRepo.CountByRole(ctx, string(constants.RoleSuperAdmin))
	if err != nil {
		return err
	}
	if count <= 1 {
		return ErrLastSuperAdmin
	}
	return nil
}

func (s *userService) ResetPassword(ctx context.Context, actorRole constants.UserRole, id uuid.UUID, newPassword string) error {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if user == nil {
		return ErrUserNotFound
	}
	if !canManage(actorRole, user.Role) {
		return ErrForbiddenUserEdit
	}
	if err := helpers.ValidatePasswordStrength(newPassword); err != nil {
		return err
	}
	hash, err := helpers.HashPassword(newPassword)
	if err != nil {
		return err
	}
	user.PasswordHash = hash
	user.FailedLoginAttempts = 0
	user.LockoutUntil = nil
	user.PasswordResetToken = nil
	user.PasswordResetExpiresAt = nil
	if err := s.userRepo.Update(ctx, user); err != nil {
		return err
	}
	return s.refreshTokenRepo.RevokeByUserID(ctx, user.ID)
}

// VerifyB2B approves or rejects a store-partner (B2B) registration. On
// approval with a credit limit, the customer's kasbon account is created or
// its limit updated.
func (s *userService) VerifyB2B(ctx context.Context, id uuid.UUID, req dto.VerifyB2BDTO) (*dto.UserProfileResponse, error) {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}
	if user.Role != constants.RoleB2BCustomer || user.B2BProfile == nil {
		return nil, errors.New("pengguna ini bukan pendaftar mitra toko (B2B)")
	}

	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		updates := map[string]interface{}{"verification_status": req.Status}
		if req.CreditLimit != nil {
			updates["approved_credit_limit"] = *req.CreditLimit
		}
		if err := tx.Model(&models.B2BProfile{}).Where("id = ?", user.B2BProfile.ID).Updates(updates).Error; err != nil {
			return err
		}
		if req.Status != constants.B2BStatusApproved || req.CreditLimit == nil {
			return nil
		}

		var ledger models.CreditLedger
		err := tx.Where("user_id = ?", user.ID).First(&ledger).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return tx.Create(&models.CreditLedger{
				UserID:          user.ID,
				CreditLimit:     *req.CreditLimit,
				RemainingCredit: *req.CreditLimit,
			}).Error
		}
		if err != nil {
			return err
		}
		remaining := *req.CreditLimit - ledger.UsedCredit
		if remaining < 0 {
			remaining = 0
		}
		return tx.Model(&models.CreditLedger{}).Where("id = ?", ledger.ID).Updates(map[string]interface{}{
			"credit_limit": *req.CreditLimit, "remaining_credit": remaining,
		}).Error
	})
	if err != nil {
		return nil, err
	}

	title, message := "Pendaftaran mitra toko diperbarui", "Status pendaftaran mitra toko Anda: "+string(req.Status)
	switch req.Status {
	case constants.B2BStatusApproved:
		title = "Selamat! Akun mitra toko disetujui"
		message = "Anda kini mendapat harga grosir"
		if req.CreditLimit != nil && *req.CreditLimit > 0 {
			message += " dan limit kasbon " + formatRupiah(*req.CreditLimit)
		}
		message += "."
	case constants.B2BStatusRejected:
		title = "Pendaftaran mitra toko belum disetujui"
		message = "Silakan hubungi toko untuk informasi lebih lanjut."
	}
	publishEvent(s.bus, "B2B_VERIFIED", "USER", user.ID, map[string]interface{}{
		"user_id": user.ID.String(), "title": title, "message": message,
	})

	return s.GetUserByID(ctx, id)
}

func (s *userService) DeleteUser(ctx context.Context, actorID uuid.UUID, actorRole constants.UserRole, id uuid.UUID) error {
	if id == actorID {
		return ErrSelfModification
	}
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if user == nil {
		return ErrUserNotFound
	}
	if !canManage(actorRole, user.Role) {
		return ErrForbiddenUserEdit
	}
	if user.Role == constants.RoleSuperAdmin {
		if err := s.ensureAnotherSuperAdmin(ctx); err != nil {
			return err
		}
	}
	_ = s.refreshTokenRepo.RevokeByUserID(ctx, id)
	return s.userRepo.Delete(ctx, id)
}

func mapUserToProfile(u *models.User) dto.UserProfileResponse {
	res := dto.UserProfileResponse{
		ID:              u.ID,
		FullName:        u.FullName,
		Email:           u.Email,
		PhoneNumber:     u.PhoneNumber,
		Role:            u.Role,
		IsActive:        u.IsActive,
		IsEmailVerified: u.IsEmailVerified,
		AvatarURL:       u.AvatarURL,
		LastLoginAt:     u.LastLoginAt,
		CreatedAt:       u.CreatedAt,
	}
	if p := u.B2BProfile; p != nil {
		res.B2BProfile = &dto.B2BProfileDTO{
			ID:                  p.ID,
			UserID:              p.UserID,
			StoreName:           p.StoreName,
			StoreAddress:        p.StoreAddress,
			TaxID:               p.TaxID,
			VerificationStatus:  p.VerificationStatus,
			ApprovedCreditLimit: p.ApprovedCreditLimit,
			CreatedAt:           p.CreatedAt,
		}
	}
	return res
}
