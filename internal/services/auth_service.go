package services

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"

	"backend-sembako/helpers"
	"backend-sembako/internal/constants"
	"backend-sembako/internal/dto"
	"backend-sembako/internal/models"
	"backend-sembako/internal/repositories"
	"backend-sembako/pkg/jwt"

	"github.com/google/uuid"
)

var (
	ErrPhoneAlreadyExists = errors.New("phone number is already registered")
	ErrEmailAlreadyExists = errors.New("email address is already registered")
	ErrInvalidCredentials = errors.New("invalid phone number or password")
	ErrAccountLocked      = errors.New("account is temporarily locked due to multiple failed login attempts. Please try again later")
	ErrUserInactive       = errors.New("account is inactive. Please contact store administrator")
	ErrInvalidToken       = errors.New("token is invalid or expired")
)

type AuthService interface {
	RegisterCustomer(ctx context.Context, req dto.RegisterCustomerDTO) (*dto.AuthResponseDTO, error)
	Login(ctx context.Context, req dto.LoginDTO, ipAddress, userAgent, device string) (*dto.AuthResponseDTO, error)
	RefreshToken(ctx context.Context, req dto.RefreshTokenDTO) (*dto.AuthResponseDTO, error)
	Logout(ctx context.Context, userID uuid.UUID) error
	GetCurrentUser(ctx context.Context, userID uuid.UUID) (*dto.UserProfileResponse, error)
	ChangePassword(ctx context.Context, userID uuid.UUID, req dto.ChangePasswordDTO) error
	ForgotPassword(ctx context.Context, req dto.ForgotPasswordDTO) (string, error)
	ResetPassword(ctx context.Context, req dto.ResetPasswordDTO) error
	VerifyEmail(ctx context.Context, req dto.VerifyEmailDTO) error
}

type authService struct {
	userRepo         repositories.UserRepository
	b2bProfileRepo   repositories.B2BProfileRepository
	refreshTokenRepo repositories.RefreshTokenRepository
	loginHistoryRepo repositories.LoginHistoryRepository
	jwtService       jwt.JWTService
}

func NewAuthService(
	userRepo repositories.UserRepository,
	b2bProfileRepo repositories.B2BProfileRepository,
	refreshTokenRepo repositories.RefreshTokenRepository,
	loginHistoryRepo repositories.LoginHistoryRepository,
	jwtService jwt.JWTService,
) AuthService {
	return &authService{
		userRepo:         userRepo,
		b2bProfileRepo:   b2bProfileRepo,
		refreshTokenRepo: refreshTokenRepo,
		loginHistoryRepo: loginHistoryRepo,
		jwtService:       jwtService,
	}
}

func (s *authService) RegisterCustomer(ctx context.Context, req dto.RegisterCustomerDTO) (*dto.AuthResponseDTO, error) {
	// Validate password strength
	if err := helpers.ValidatePasswordStrength(req.Password); err != nil {
		return nil, err
	}

	// Check existing phone number
	existingPhone, err := s.userRepo.FindByPhoneNumber(ctx, req.PhoneNumber)
	if err != nil {
		return nil, err
	}
	if existingPhone != nil {
		return nil, ErrPhoneAlreadyExists
	}

	// Check existing email
	if req.Email != nil && *req.Email != "" {
		existingEmail, err := s.userRepo.FindByEmail(ctx, *req.Email)
		if err != nil {
			return nil, err
		}
		if existingEmail != nil {
			return nil, ErrEmailAlreadyExists
		}
	}

	// Hash password
	hashedPassword, err := helpers.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	role := constants.RoleB2CCustomer
	if req.IsB2BRegistration {
		role = constants.RoleB2BCustomer
	}

	// Email Verification Token
	verToken := uuid.New().String()
	verExpires := time.Now().Add(24 * time.Hour)

	user := &models.User{
		FullName:                   req.FullName,
		Email:                      req.Email,
		PhoneNumber:                req.PhoneNumber,
		PasswordHash:               hashedPassword,
		Role:                       role,
		IsActive:                   true,
		IsEmailVerified:            false,
		EmailVerificationToken:     &verToken,
		EmailVerificationExpiresAt: &verExpires,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	// Create B2B Profile if applicable
	var b2bDto *dto.B2BProfileDTO
	if req.IsB2BRegistration {
		b2bProfile := &models.B2BProfile{
			UserID:              user.ID,
			StoreName:           req.StoreName,
			StoreAddress:        req.StoreAddress,
			TaxID:               req.TaxID,
			VerificationStatus:  constants.B2BStatusPending,
			ApprovedCreditLimit: 0,
		}
		if err := s.b2bProfileRepo.Create(ctx, b2bProfile); err != nil {
			return nil, err
		}
		b2bDto = &dto.B2BProfileDTO{
			ID:                  b2bProfile.ID,
			UserID:              b2bProfile.UserID,
			StoreName:           b2bProfile.StoreName,
			StoreAddress:        b2bProfile.StoreAddress,
			TaxID:               b2bProfile.TaxID,
			VerificationStatus:  b2bProfile.VerificationStatus,
			ApprovedCreditLimit: b2bProfile.ApprovedCreditLimit,
			CreatedAt:           b2bProfile.CreatedAt,
		}
	}

	// Issue JWT Tokens
	tokenPair, err := s.jwtService.GenerateTokenPair(user.ID, user.Role)
	if err != nil {
		return nil, err
	}

	// Store Refresh Token Hash
	s.saveRefreshTokenHash(ctx, user.ID, tokenPair.RefreshToken)

	return &dto.AuthResponseDTO{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		TokenType:    tokenPair.TokenType,
		ExpiresIn:    tokenPair.ExpiresIn,
		User: dto.UserProfileResponse{
			ID:              user.ID,
			FullName:        user.FullName,
			Email:           user.Email,
			PhoneNumber:     user.PhoneNumber,
			Role:            user.Role,
			IsActive:        user.IsActive,
			IsEmailVerified: user.IsEmailVerified,
			B2BProfile:      b2bDto,
		},
	}, nil
}

func (s *authService) Login(ctx context.Context, req dto.LoginDTO, ipAddress, userAgent, device string) (*dto.AuthResponseDTO, error) {
	user, err := s.userRepo.FindByPhoneNumber(ctx, req.PhoneNumber)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrInvalidCredentials
	}

	if !user.IsActive {
		return nil, ErrUserInactive
	}

	// Check Account Lockout
	if user.IsLockedOut() {
		_ = s.loginHistoryRepo.Create(ctx, &models.LoginHistory{
			UserID:    user.ID,
			IPAddress: ipAddress,
			UserAgent: userAgent,
			Device:    device,
			Status:    "LOCKED",
		})
		return nil, ErrAccountLocked
	}

	// Validate Password
	if !helpers.CheckPasswordHash(req.Password, user.PasswordHash) {
		user.FailedLoginAttempts++
		if user.FailedLoginAttempts >= 5 {
			lockUntil := time.Now().Add(15 * time.Minute)
			user.LockoutUntil = &lockUntil
		}
		_ = s.userRepo.Update(ctx, user)

		_ = s.loginHistoryRepo.Create(ctx, &models.LoginHistory{
			UserID:    user.ID,
			IPAddress: ipAddress,
			UserAgent: userAgent,
			Device:    device,
			Status:    "FAILED",
		})
		return nil, ErrInvalidCredentials
	}

	// Login Success: Reset Lockout Counters & Update LastLogin
	now := time.Now()
	user.FailedLoginAttempts = 0
	user.LockoutUntil = nil
	user.LastLoginAt = &now
	_ = s.userRepo.Update(ctx, user)

	_ = s.loginHistoryRepo.Create(ctx, &models.LoginHistory{
		UserID:    user.ID,
		IPAddress: ipAddress,
		UserAgent: userAgent,
		Device:    device,
		Status:    "SUCCESS",
	})

	// Issue Tokens
	tokenPair, err := s.jwtService.GenerateTokenPair(user.ID, user.Role)
	if err != nil {
		return nil, err
	}

	s.saveRefreshTokenHash(ctx, user.ID, tokenPair.RefreshToken)

	var b2bDto *dto.B2BProfileDTO
	if user.B2BProfile != nil {
		b2bDto = &dto.B2BProfileDTO{
			ID:                  user.B2BProfile.ID,
			UserID:              user.B2BProfile.UserID,
			StoreName:           user.B2BProfile.StoreName,
			StoreAddress:        user.B2BProfile.StoreAddress,
			TaxID:               user.B2BProfile.TaxID,
			VerificationStatus:  user.B2BProfile.VerificationStatus,
			ApprovedCreditLimit: user.B2BProfile.ApprovedCreditLimit,
			CreatedAt:           user.B2BProfile.CreatedAt,
		}
	}

	return &dto.AuthResponseDTO{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		TokenType:    tokenPair.TokenType,
		ExpiresIn:    tokenPair.ExpiresIn,
		User: dto.UserProfileResponse{
			ID:              user.ID,
			FullName:        user.FullName,
			Email:           user.Email,
			PhoneNumber:     user.PhoneNumber,
			Role:            user.Role,
			IsActive:        user.IsActive,
			IsEmailVerified: user.IsEmailVerified,
			AvatarURL:       user.AvatarURL,
			LastLoginAt:     user.LastLoginAt,
			B2BProfile:      b2bDto,
		},
	}, nil
}

func (s *authService) RefreshToken(ctx context.Context, req dto.RefreshTokenDTO) (*dto.AuthResponseDTO, error) {
	claims, err := s.jwtService.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		return nil, ErrInvalidToken
	}

	hash := hashTokenString(req.RefreshToken)
	storedToken, err := s.refreshTokenRepo.FindByTokenHash(ctx, hash)
	if err != nil || storedToken == nil {
		return nil, ErrInvalidToken
	}

	// Refresh Token Rotation: Revoke current refresh token
	_ = s.refreshTokenRepo.RevokeToken(ctx, storedToken.ID)

	user, err := s.userRepo.FindByID(ctx, claims.UserID)
	if err != nil || user == nil || !user.IsActive {
		return nil, ErrUserInactive
	}

	newTokenPair, err := s.jwtService.GenerateTokenPair(user.ID, user.Role)
	if err != nil {
		return nil, err
	}

	s.saveRefreshTokenHash(ctx, user.ID, newTokenPair.RefreshToken)

	return &dto.AuthResponseDTO{
		AccessToken:  newTokenPair.AccessToken,
		RefreshToken: newTokenPair.RefreshToken,
		TokenType:    newTokenPair.TokenType,
		ExpiresIn:    newTokenPair.ExpiresIn,
		User: dto.UserProfileResponse{
			ID:              user.ID,
			FullName:        user.FullName,
			Email:           user.Email,
			PhoneNumber:     user.PhoneNumber,
			Role:            user.Role,
			IsActive:        user.IsActive,
			IsEmailVerified: user.IsEmailVerified,
		},
	}, nil
}

func (s *authService) Logout(ctx context.Context, userID uuid.UUID) error {
	return s.refreshTokenRepo.RevokeByUserID(ctx, userID)
}

func (s *authService) GetCurrentUser(ctx context.Context, userID uuid.UUID) (*dto.UserProfileResponse, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil || user == nil {
		return nil, errors.New("user profile not found")
	}

	var b2bDto *dto.B2BProfileDTO
	if user.B2BProfile != nil {
		b2bDto = &dto.B2BProfileDTO{
			ID:                  user.B2BProfile.ID,
			UserID:              user.B2BProfile.UserID,
			StoreName:           user.B2BProfile.StoreName,
			StoreAddress:        user.B2BProfile.StoreAddress,
			TaxID:               user.B2BProfile.TaxID,
			VerificationStatus:  user.B2BProfile.VerificationStatus,
			ApprovedCreditLimit: user.B2BProfile.ApprovedCreditLimit,
			CreatedAt:           user.B2BProfile.CreatedAt,
		}
	}

	return &dto.UserProfileResponse{
		ID:              user.ID,
		FullName:        user.FullName,
		Email:           user.Email,
		PhoneNumber:     user.PhoneNumber,
		Role:            user.Role,
		IsActive:        user.IsActive,
		IsEmailVerified: user.IsEmailVerified,
		AvatarURL:       user.AvatarURL,
		LastLoginAt:     user.LastLoginAt,
		B2BProfile:      b2bDto,
	}, nil
}

func (s *authService) ChangePassword(ctx context.Context, userID uuid.UUID, req dto.ChangePasswordDTO) error {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil || user == nil {
		return errors.New("user not found")
	}

	if !helpers.CheckPasswordHash(req.OldPassword, user.PasswordHash) {
		return errors.New("incorrect old password")
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

func (s *authService) ForgotPassword(ctx context.Context, req dto.ForgotPasswordDTO) (string, error) {
	user, err := s.userRepo.FindByPhoneNumber(ctx, req.PhoneNumber)
	if err != nil || user == nil {
		return "", errors.New("user with specified phone number not found")
	}

	resetToken := uuid.New().String()
	expiresAt := time.Now().Add(1 * time.Hour)

	user.PasswordResetToken = &resetToken
	user.PasswordResetExpiresAt = &expiresAt

	if err := s.userRepo.Update(ctx, user); err != nil {
		return "", err
	}

	return resetToken, nil
}

func (s *authService) ResetPassword(ctx context.Context, req dto.ResetPasswordDTO) error {
	user, err := s.userRepo.FindByPasswordResetToken(ctx, req.ResetToken)
	if err != nil || user == nil {
		return errors.New("invalid or expired password reset token")
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

	return s.userRepo.Update(ctx, user)
}

func (s *authService) VerifyEmail(ctx context.Context, req dto.VerifyEmailDTO) error {
	user, err := s.userRepo.FindByEmailVerificationToken(ctx, req.Token)
	if err != nil || user == nil {
		return errors.New("invalid or expired email verification token")
	}

	user.IsEmailVerified = true
	user.EmailVerificationToken = nil
	user.EmailVerificationExpiresAt = nil

	return s.userRepo.Update(ctx, user)
}

func (s *authService) saveRefreshTokenHash(ctx context.Context, userID uuid.UUID, refreshTokenStr string) {
	hash := hashTokenString(refreshTokenStr)
	tokenRecord := &models.RefreshToken{
		UserID:    userID,
		TokenHash: hash,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
		IsRevoked: false,
	}
	_ = s.refreshTokenRepo.Create(ctx, tokenRecord)
}

func hashTokenString(token string) string {
	h := sha256.Sum256([]byte(token))
	return hex.EncodeToString(h[:])
}
