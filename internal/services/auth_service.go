package services

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strings"
	"time"

	"backend-sembako/helpers"
	"backend-sembako/internal/constants"
	"backend-sembako/internal/dto"
	"backend-sembako/internal/models"
	"backend-sembako/internal/repositories"
	"backend-sembako/pkg/eventbus"
	"backend-sembako/pkg/jwt"

	"github.com/google/uuid"
)

var (
	ErrPhoneAlreadyExists = errors.New("nomor HP sudah terdaftar")
	ErrEmailAlreadyExists = errors.New("email sudah terdaftar")
	ErrInvalidCredentials = errors.New("nomor HP atau kata sandi salah")
	ErrAccountLocked      = errors.New("akun dikunci sementara karena terlalu banyak percobaan masuk yang gagal; coba lagi 15 menit lagi")
	ErrUserInactive       = errors.New("akun tidak aktif; hubungi admin toko")
	ErrInvalidToken       = errors.New("sesi tidak valid atau sudah kedaluwarsa; silakan masuk lagi")
)

const (
	maxFailedLogins = 5
	lockoutDuration = 15 * time.Minute
)

type AuthService interface {
	RegisterCustomer(ctx context.Context, req dto.RegisterCustomerDTO) (*dto.AuthResponseDTO, error)
	Login(ctx context.Context, req dto.LoginDTO, ipAddress, userAgent, device string) (*dto.AuthResponseDTO, error)
	RefreshToken(ctx context.Context, req dto.RefreshTokenDTO) (*dto.AuthResponseDTO, error)
	Logout(ctx context.Context, userID uuid.UUID) error
	GetCurrentUser(ctx context.Context, userID uuid.UUID) (*dto.UserProfileResponse, error)
	UpdateProfile(ctx context.Context, userID uuid.UUID, req dto.UpdateProfileDTO) (*dto.UserProfileResponse, error)
	ChangePassword(ctx context.Context, userID uuid.UUID, req dto.ChangePasswordDTO) error
	ForgotPassword(ctx context.Context, req dto.ForgotPasswordDTO) error
	ResetPassword(ctx context.Context, req dto.ResetPasswordDTO) error
	VerifyEmail(ctx context.Context, req dto.VerifyEmailDTO) error
}

type authService struct {
	userRepo         repositories.UserRepository
	b2bProfileRepo   repositories.B2BProfileRepository
	refreshTokenRepo repositories.RefreshTokenRepository
	loginHistoryRepo repositories.LoginHistoryRepository
	jwtService       jwt.JWTService
	refreshTTL       time.Duration
	bus              eventbus.EventBus
}

func NewAuthService(
	userRepo repositories.UserRepository,
	b2bProfileRepo repositories.B2BProfileRepository,
	refreshTokenRepo repositories.RefreshTokenRepository,
	loginHistoryRepo repositories.LoginHistoryRepository,
	jwtService jwt.JWTService,
	refreshTTL time.Duration,
	bus eventbus.EventBus,
) AuthService {
	return &authService{
		userRepo:         userRepo,
		b2bProfileRepo:   b2bProfileRepo,
		refreshTokenRepo: refreshTokenRepo,
		loginHistoryRepo: loginHistoryRepo,
		jwtService:       jwtService,
		refreshTTL:       refreshTTL,
		bus:              bus,
	}
}

func (s *authService) RegisterCustomer(ctx context.Context, req dto.RegisterCustomerDTO) (*dto.AuthResponseDTO, error) {
	if err := helpers.ValidatePasswordStrength(req.Password); err != nil {
		return nil, err
	}

	existingPhone, err := s.userRepo.FindByPhoneNumber(ctx, req.PhoneNumber)
	if err != nil {
		return nil, err
	}
	if existingPhone != nil {
		return nil, ErrPhoneAlreadyExists
	}

	if req.Email != nil && strings.TrimSpace(*req.Email) != "" {
		existingEmail, err := s.userRepo.FindByEmail(ctx, *req.Email)
		if err != nil {
			return nil, err
		}
		if existingEmail != nil {
			return nil, ErrEmailAlreadyExists
		}
	} else {
		req.Email = nil
	}

	hashedPassword, err := helpers.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	role := constants.RoleB2CCustomer
	if req.IsB2BRegistration {
		role = constants.RoleB2BCustomer
	}

	verToken := uuid.New().String()
	verExpires := time.Now().Add(24 * time.Hour)

	user := &models.User{
		FullName:                   strings.TrimSpace(req.FullName),
		Email:                      req.Email,
		PhoneNumber:                req.PhoneNumber,
		PasswordHash:               hashedPassword,
		Role:                       role,
		IsActive:                   true,
		EmailVerificationToken:     &verToken,
		EmailVerificationExpiresAt: &verExpires,
	}
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	if req.IsB2BRegistration {
		if err := s.b2bProfileRepo.Create(ctx, &models.B2BProfile{
			UserID:             user.ID,
			StoreName:          strings.TrimSpace(req.StoreName),
			StoreAddress:       strings.TrimSpace(req.StoreAddress),
			TaxID:              req.TaxID,
			VerificationStatus: constants.B2BStatusPending,
		}); err != nil {
			return nil, err
		}
		publishEvent(s.bus, "B2B_REGISTERED", "USER", user.ID, map[string]interface{}{
			"notify_staff":  true,
			"staff_title":   "Pendaftaran mitra toko baru",
			"staff_message": user.FullName + " (" + req.StoreName + ") mendaftar sebagai mitra toko. Mohon diverifikasi di menu Pengguna.",
		})
	}

	return s.issueTokens(ctx, user.ID)
}

func (s *authService) Login(ctx context.Context, req dto.LoginDTO, ipAddress, userAgent, device string) (*dto.AuthResponseDTO, error) {
	user, err := s.userRepo.FindByPhoneNumber(ctx, req.PhoneNumber)
	if err != nil {
		return nil, err
	}
	if user == nil {
		// Hash anyway so response time does not reveal whether the number exists.
		_ = helpers.CheckPasswordHash(req.Password, dummyPasswordHash)
		return nil, ErrInvalidCredentials
	}

	history := &models.LoginHistory{UserID: user.ID, IPAddress: ipAddress, UserAgent: userAgent, Device: device}

	if user.IsLockedOut() {
		history.Status = "LOCKED"
		_ = s.loginHistoryRepo.Create(ctx, history)
		return nil, ErrAccountLocked
	}

	if !helpers.CheckPasswordHash(req.Password, user.PasswordHash) {
		user.FailedLoginAttempts++
		if user.FailedLoginAttempts >= maxFailedLogins {
			lockUntil := time.Now().Add(lockoutDuration)
			user.LockoutUntil = &lockUntil
			user.FailedLoginAttempts = 0
		}
		_ = s.userRepo.Update(ctx, user)

		history.Status = "FAILED"
		_ = s.loginHistoryRepo.Create(ctx, history)
		return nil, ErrInvalidCredentials
	}

	if !user.IsActive {
		return nil, ErrUserInactive
	}

	now := time.Now()
	user.FailedLoginAttempts = 0
	user.LockoutUntil = nil
	user.LastLoginAt = &now
	_ = s.userRepo.Update(ctx, user)

	history.Status = "SUCCESS"
	_ = s.loginHistoryRepo.Create(ctx, history)

	return s.issueTokens(ctx, user.ID)
}

func (s *authService) RefreshToken(ctx context.Context, req dto.RefreshTokenDTO) (*dto.AuthResponseDTO, error) {
	claims, err := s.jwtService.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		return nil, ErrInvalidToken
	}

	storedToken, err := s.refreshTokenRepo.FindByTokenHash(ctx, hashTokenString(req.RefreshToken))
	if err != nil || storedToken == nil || storedToken.UserID != claims.UserID {
		return nil, ErrInvalidToken
	}

	// Refresh token rotation: each refresh token can be used only once.
	_ = s.refreshTokenRepo.RevokeToken(ctx, storedToken.ID)

	user, err := s.userRepo.FindByID(ctx, claims.UserID)
	if err != nil || user == nil || !user.IsActive {
		return nil, ErrUserInactive
	}

	return s.issueTokens(ctx, user.ID)
}

func (s *authService) issueTokens(ctx context.Context, userID uuid.UUID) (*dto.AuthResponseDTO, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil || user == nil {
		return nil, ErrUserNotFound
	}

	pair, err := s.jwtService.GenerateTokenPair(user.ID, user.Role)
	if err != nil {
		return nil, err
	}

	if err := s.refreshTokenRepo.Create(ctx, &models.RefreshToken{
		UserID:    user.ID,
		TokenHash: hashTokenString(pair.RefreshToken),
		ExpiresAt: time.Now().Add(s.refreshTTL),
	}); err != nil {
		return nil, err
	}

	return &dto.AuthResponseDTO{
		AccessToken:  pair.AccessToken,
		RefreshToken: pair.RefreshToken,
		TokenType:    pair.TokenType,
		ExpiresIn:    pair.ExpiresIn,
		User:         mapUserToProfile(user),
	}, nil
}

func (s *authService) Logout(ctx context.Context, userID uuid.UUID) error {
	return s.refreshTokenRepo.RevokeByUserID(ctx, userID)
}

func (s *authService) GetCurrentUser(ctx context.Context, userID uuid.UUID) (*dto.UserProfileResponse, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}
	res := mapUserToProfile(user)
	return &res, nil
}

func (s *authService) UpdateProfile(ctx context.Context, userID uuid.UUID, req dto.UpdateProfileDTO) (*dto.UserProfileResponse, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	if req.FullName != nil {
		user.FullName = strings.TrimSpace(*req.FullName)
	}
	if req.Email != nil {
		email := strings.TrimSpace(*req.Email)
		if email == "" {
			user.Email = nil
		} else if user.Email == nil || *user.Email != email {
			if other, err := s.userRepo.FindByEmail(ctx, email); err != nil {
				return nil, err
			} else if other != nil && other.ID != user.ID {
				return nil, ErrEmailAlreadyExists
			}
			user.Email = &email
			user.IsEmailVerified = false
		}
	}

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}
	return s.GetCurrentUser(ctx, userID)
}

func (s *authService) ChangePassword(ctx context.Context, userID uuid.UUID, req dto.ChangePasswordDTO) error {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}
	if user == nil {
		return ErrUserNotFound
	}

	if !helpers.CheckPasswordHash(req.OldPassword, user.PasswordHash) {
		return errors.New("kata sandi lama salah")
	}
	if err := helpers.ValidatePasswordStrength(req.NewPassword); err != nil {
		return err
	}

	newHash, err := helpers.HashPassword(req.NewPassword)
	if err != nil {
		return err
	}

	user.PasswordHash = newHash
	return s.userRepo.Update(ctx, user)
}

// ForgotPassword creates a reset token for the account. The token is never
// returned to the caller: it must reach the customer through a channel only
// they control (WhatsApp/SMS/email), which is not connected yet. Until then
// the store resets passwords from the admin dashboard. The response is the
// same whether or not the number exists, so it cannot be used to discover
// registered numbers.
func (s *authService) ForgotPassword(ctx context.Context, req dto.ForgotPasswordDTO) error {
	user, err := s.userRepo.FindByPhoneNumber(ctx, req.PhoneNumber)
	if err != nil || user == nil {
		return nil
	}

	resetToken := uuid.New().String()
	expiresAt := time.Now().Add(1 * time.Hour)
	user.PasswordResetToken = &resetToken
	user.PasswordResetExpiresAt = &expiresAt
	if err := s.userRepo.Update(ctx, user); err != nil {
		return err
	}

	publishEvent(s.bus, "PASSWORD_RESET_REQUESTED", "USER", user.ID, map[string]interface{}{
		"notify_staff":  true,
		"staff_title":   "Permintaan reset kata sandi",
		"staff_message": user.FullName + " (" + user.PhoneNumber + ") meminta reset kata sandi. Verifikasi identitasnya lalu reset dari menu Pengguna.",
	})
	return nil
}

func (s *authService) ResetPassword(ctx context.Context, req dto.ResetPasswordDTO) error {
	user, err := s.userRepo.FindByPasswordResetToken(ctx, req.ResetToken)
	if err != nil || user == nil {
		return errors.New("token reset tidak valid atau sudah kedaluwarsa")
	}
	if err := helpers.ValidatePasswordStrength(req.NewPassword); err != nil {
		return err
	}

	newHash, err := helpers.HashPassword(req.NewPassword)
	if err != nil {
		return err
	}

	user.PasswordHash = newHash
	user.PasswordResetToken = nil
	user.PasswordResetExpiresAt = nil
	user.FailedLoginAttempts = 0
	user.LockoutUntil = nil
	if err := s.userRepo.Update(ctx, user); err != nil {
		return err
	}
	return s.refreshTokenRepo.RevokeByUserID(ctx, user.ID)
}

func (s *authService) VerifyEmail(ctx context.Context, req dto.VerifyEmailDTO) error {
	user, err := s.userRepo.FindByEmailVerificationToken(ctx, req.Token)
	if err != nil || user == nil {
		return errors.New("token verifikasi email tidak valid atau sudah kedaluwarsa")
	}

	user.IsEmailVerified = true
	user.EmailVerificationToken = nil
	user.EmailVerificationExpiresAt = nil
	return s.userRepo.Update(ctx, user)
}

func hashTokenString(token string) string {
	h := sha256.Sum256([]byte(token))
	return hex.EncodeToString(h[:])
}

// dummyPasswordHash is compared against when a phone number is unknown so a
// failed login takes the same time whether or not the account exists.
var dummyPasswordHash, _ = helpers.HashPassword("lucky-mart-timing-equalizer")
