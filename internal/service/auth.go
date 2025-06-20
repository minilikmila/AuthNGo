package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	config "github.com/minilikmila/standard-auth-go/configs"
	"github.com/minilikmila/standard-auth-go/internal/domain/model"
	"github.com/minilikmila/standard-auth-go/internal/infrastructure/database"
	"github.com/minilikmila/standard-auth-go/pkg/utils"
)

// AuthServiceImpl implements the AuthService interface
type AuthServiceImpl struct {
	repo         *database.Repository
	config       *config.Config
	emailService EmailService
	smsService   SMSService
	jwtService   JWTService
}

// NewAuthService creates a new instance of AuthService
func NewAuthService(repo *database.Repository, config *config.Config, emailService EmailService, smsService SMSService, jwtService JWTService) AuthService {
	return &AuthServiceImpl{
		repo:         repo,
		config:       config,
		emailService: emailService,
		smsService:   smsService,
		jwtService:   jwtService,
	}
}

// CreateUser implements user creation
func (s *AuthServiceImpl) CreateUser(ctx context.Context, user *model.User) error {
	// Check if user already exists
	exists, err := s.repo.UserExists(ctx, *user.Email)
	if err != nil {
		return err
	}
	if exists {
		return ErrUserAlreadyExists
	}

	// Hash password
	hashedPassword, err := utils.EncryptPassword(*user.Password, 10)
	if err != nil {
		return err
	}
	user.Password = &hashedPassword

	// Create user
	return s.repo.CreateUser(ctx, user)
}

// UserExists checks if a user exists
func (s *AuthServiceImpl) UserExists(ctx context.Context, email string) (bool, error) {
	return s.repo.UserExists(ctx, email)
}

// Login implements user login
func (s *AuthServiceImpl) Login(ctx context.Context, email, password string) (*model.User, error) {
	//
	// Get user by email
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	isEmailVerified, err := s.repo.IsEmailVerified(ctx, *user.ID)
	if err != nil {
		return nil, err
	}

	if !user.IsEmailVerified || !isEmailVerified {
		return nil, ErrEmailNotVerified
	}

	// Verify password
	if !utils.ComparePassword(*user.Password, password) {
		// Increment incorrect login attempts
		if err := s.repo.IncrementLoginAttempts(ctx, *user.ID); err != nil {
			return nil, err
		}
		return nil, ErrInvalidCredentials
	}

	// Reset incorrect login attempts
	if err := s.repo.ResetLoginAttempts(ctx, *user.ID); err != nil {
		return nil, err
	}

	// Update last login time
	if err := s.repo.UpdateLastLogin(ctx, *user.ID); err != nil {
		return nil, err
	}

	return user, nil
}

// Logout implements user logout
func (s *AuthServiceImpl) Logout(ctx context.Context, token string) error {
	// Add token to blacklist and hashed it for security and compromising traits
	hashedToken := utils.HashToken(token)
	logrus.Infoln("Hashed token to logout : ", hashedToken)
	return s.repo.BlacklistToken(ctx, hashedToken)
}
func derefString(s *string) string {
	if s != nil {
		return *s
	}
	return ""
}

// GenerateTokens generates access and refresh tokens
func (s *AuthServiceImpl) GenerateTokens(ctx context.Context, user *model.User) (string, string, error) {
	// Check if email is verified
	if !user.IsEmailVerified {
		return "", "", ErrEmailNotVerified
	}

	metadata := map[string]string{
		"userId":          user.ID.String(),
		"email":           derefString(user.Email),
		"phone_number":    derefString(user.Phone),
		"profile_picture": derefString(user.ProfilePicture),
		"signup_method":   user.SignUpMethod,
		"role":            user.Role,
		"ip":              "",
	}
	// Generate access token
	accessToken, err := s.jwtService.GenerateToken(user, metadata)
	if err != nil {
		return "", "", err
	}

	// Generate refresh token
	refreshToken, err := s.jwtService.GenerateRefreshToken(user)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

// SendVerificationEmail sends email verification
func (s *AuthServiceImpl) SendVerificationEmail(ctx context.Context, user *model.User) error {
	// Generate verification token
	token, err := s.jwtService.GenerateEmailVerificationToken(user)
	if err != nil {
		return err
	}

	hashedToken := utils.HashToken(token) // hash a token before save to db, SECURITY may be due to compromising db security
	logrus.Infoln("Hashed verification token : ", hashedToken)

	// Save verification to database
	verification := &model.Verification{
		UserID:    *user.ID,
		Type:      model.VerificationTypeEmail,
		Token:     &hashedToken,
		Status:    model.VerificationStatusPending,
		SentAt:    time.Now(),
		ExpiresAt: time.Now().Add(24 * time.Hour),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := s.repo.CreateVerification(ctx, verification); err != nil {
		return err
	}

	// Send verification email
	return s.emailService.SendVerificationEmail(ctx, *user.Email, token, *user.Name)
}

// VerifyEmail verifies email token and updates user in a transaction
func (s *AuthServiceImpl) VerifyEmail(ctx context.Context, token string) error {
	// 1. Validate the JWT token (checks signature and exp)
	userID, err := s.jwtService.ValidateEmailVerificationToken(token)
	if err != nil {
		return ErrInvalidToken
	}

	logrus.Infoln("User ID : ", userID)

	// 2. Check the database for a pending verification record
	hashedToken := utils.HashToken(token) // hashed before query since it's hashed when stored

	verification, err := s.repo.GetVerificationByToken(ctx, hashedToken, model.VerificationTypeEmail)
	if err != nil || verification.UserID != userID {
		return ErrInvalidToken
	}

	// 3. Use a transaction to mark as verified and update user
	return s.repo.Transaction(ctx, func(tx *gorm.DB) error {
		if err := tx.Model(&model.Verification{}).
			Where("id = ?", verification.ID).
			Updates(map[string]interface{}{
				"status":      model.VerificationStatusVerified,
				"updated_at":  time.Now(),
				"verified_at": time.Now(),
			}).Error; err != nil {
			return err
		}
		if err := tx.Model(&model.User{}).
			Where("id = ?", verification.UserID).
			Updates(map[string]interface{}{
				"is_email_verified": true,
				"email_verified_at": time.Now(),
			}).Error; err != nil {
			return err
		}
		return nil
	})
}

// SendVerificationSMS sends SMS verification
func (s *AuthServiceImpl) SendVerificationSMS(ctx context.Context, user *model.User) error {
	code := utils.GenerateVerificationCode()
	verification := &model.Verification{
		UserID:    *user.ID,
		Type:      model.VerificationTypePhone,
		Code:      &code,
		Status:    model.VerificationStatusPending,
		SentAt:    time.Now(),
		ExpiresAt: time.Now().Add(15 * time.Minute),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := s.repo.CreateVerification(ctx, verification); err != nil {
		return err
	}
	return s.smsService.SendVerificationSMS(ctx, *user.Phone, code)
}

// VerifyPhone verifies phone code and updates user in a transaction
func (s *AuthServiceImpl) VerifyPhone(ctx context.Context, code, phone string) error {
	verification, err := s.repo.GetVerificationByCode(ctx, code, model.VerificationTypePhone)
	if err != nil {
		return ErrInvalidCode
	}

	// Fetch the user
	user, err := s.repo.GetUserByID(ctx, verification.UserID)
	if err != nil {
		return ErrUserNotFound
	}

	// Compare the provided phone with the user's phone
	if user.Phone == nil || *user.Phone != phone {
		return ErrInvalidCode // or ErrPhoneMismatch
	}

	// Transaction: mark verification as verified and update user
	return s.repo.Transaction(ctx, func(tx *gorm.DB) error {
		if err := tx.Model(&model.Verification{}).
			Where("id = ?", verification.ID).
			Updates(map[string]interface{}{
				"status":      model.VerificationStatusVerified,
				"updated_at":  time.Now(),
				"verified_at": time.Now(),
			}).Error; err != nil {
			return err
		}
		if err := tx.Model(&model.User{}).
			Where("id = ?", verification.UserID).
			Updates(map[string]interface{}{
				"is_phone_verified": true,
				"phone_verified_at": time.Now(),
			}).Error; err != nil {
			return err
		}
		return nil
	})
}

// ForgotPassword implements password reset request
func (s *AuthServiceImpl) ForgotPassword(ctx context.Context, email string) error {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return ErrUserNotFound
	}
	token, err := s.jwtService.GeneratePasswordResetToken(user)
	if err != nil {
		return err
	}
	hashedToken := utils.HashToken(token)
	verification := &model.Verification{
		UserID:    *user.ID,
		Type:      model.VerificationTypeReset,
		Token:     &hashedToken,
		Status:    model.VerificationStatusPending,
		SentAt:    time.Now(),
		ExpiresAt: time.Now().Add(1 * time.Hour),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := s.repo.CreateVerification(ctx, verification); err != nil {
		return err
	}
	return s.emailService.SendPasswordResetEmail(ctx, *user.Email, token)
}

// ResetPassword implements password reset .. used transaction here since we do've > 1 database mutation
func (s *AuthServiceImpl) ResetPassword(ctx context.Context, token, newPassword string) error {
	// 1. Validate the JWT token (checks signature and exp)
	userID, err := s.jwtService.ValidateEmailVerificationToken(token)
	if err != nil {
		return ErrInvalidToken
	}

	hashedToken := utils.HashToken(token) // hashed before query since it's hashed when stored
	verification, err := s.repo.GetVerificationByToken(ctx, hashedToken, model.VerificationTypeReset)
	if err != nil || userID != verification.UserID {
		return ErrInvalidToken
	}

	hashedPassword, err := utils.EncryptPassword(newPassword, 10)
	if err != nil {
		return err
	}

	// Use a transaction for both updates
	err = s.repo.Transaction(ctx, func(tx *gorm.DB) error {
		// Update user's password
		if err := tx.Model(&model.User{}).
			Where("id = ?", verification.UserID).
			Updates(map[string]interface{}{
				"password":            hashedPassword,
				"password_changed_at": time.Now(),
			}).Error; err != nil {
			return err
		}

		// Mark the verification as used/verified
		if err := tx.Model(&model.Verification{}).
			Where("id = ?", verification.ID).
			Updates(map[string]interface{}{
				"status":      model.VerificationStatusVerified,
				"updated_at":  time.Now(),
				"verified_at": time.Now(),
			}).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}
	return nil
}

func (s *AuthServiceImpl) UpdatePassword(ctx context.Context, userID string, currentPassword, newPassword string) error {
	id, err := uuid.Parse(userID)
	if err != nil {
		return err
	}

	// Get user
	user, err := s.repo.GetUserByID(ctx, id)
	if err != nil {
		return ErrUserNotFound
	}

	// Verify current password
	if !utils.ComparePassword(*user.Password, currentPassword) {
		return ErrInvalidCredentials
	}

	// Hash new password
	hashedPassword, err := utils.EncryptPassword(newPassword, 10)
	if err != nil {
		return err
	}

	// Update password
	return s.repo.UpdatePassword(ctx, id, hashedPassword)
}

// RefreshToken implements token refresh
func (s *AuthServiceImpl) RefreshToken(ctx context.Context, refreshToken string) (string, string, error) {
	// Validate refresh token
	userID, err := s.jwtService.ValidateRefreshToken(refreshToken)
	if err != nil {
		return "", "", ErrInvalidToken
	}

	// Get user
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return "", "", ErrUserNotFound
	}

	// Generate new tokens
	return s.GenerateTokens(ctx, user)
}

// ValidateToken validates access token
func (s *AuthServiceImpl) ValidateToken(ctx context.Context, token string) (string, error) {

	// compare hashed value
	hashedToken := utils.HashToken(token)
	// Check if token is blacklisted
	isBlacklisted, err := s.repo.IsTokenBlacklisted(ctx, hashedToken)
	if err != nil {
		return "", err
	}
	if isBlacklisted {
		return "", ErrInvalidToken
	}

	// Validate token
	userID, err := s.jwtService.ValidateToken(token)
	if err != nil {
		return "", err
	}

	return userID.String(), nil
}

// UpdateProfile implements profile update
func (s *AuthServiceImpl) UpdateProfile(ctx context.Context, userID string, updates map[string]interface{}) error {
	id, err := uuid.Parse(userID)
	if err != nil {
		return err
	}
	return s.repo.UpdateUser(ctx, id, updates)
}

// GetProfile implements profile retrieval
func (s *AuthServiceImpl) GetProfile(ctx context.Context, userID string) (*model.User, error) {
	id, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}
	return s.repo.GetUserByID(ctx, id)
}

// ChangeEmail changes the user's email after verification TODO: needs verification implementation
func (s *AuthServiceImpl) ChangeEmail(ctx context.Context, userID uuid.UUID, newEmail, token string) error {
	verification, err := s.repo.GetVerificationByToken(ctx, token, model.VerificationTypeEmail)
	if err != nil {
		return ErrInvalidToken
	}
	return s.repo.Transaction(ctx, func(tx *gorm.DB) error {
		// Mark verification as verified
		if err := tx.Model(&model.Verification{}).
			Where("id = ?", verification.ID).
			Updates(map[string]interface{}{
				"status":      model.VerificationStatusVerified,
				"updated_at":  time.Now(),
				"verified_at": time.Now(),
			}).Error; err != nil {
			return err
		}
		// Update user's email
		if err := tx.Model(&model.User{}).
			Where("id = ?", userID).
			Update("email", newEmail).Error; err != nil {
			return err
		}
		return nil
	})
}

// ChangePhone changes the user's phone after verification TODO: needs verification implementation
func (s *AuthServiceImpl) ChangePhone(ctx context.Context, userID uuid.UUID, newPhone, code string) error {
	verification, err := s.repo.GetVerificationByCode(ctx, code, model.VerificationTypePhone)
	if err != nil {
		return ErrInvalidCode
	}
	return s.repo.Transaction(ctx, func(tx *gorm.DB) error {
		// Mark verification as verified
		if err := tx.Model(&model.Verification{}).
			Where("id = ?", verification.ID).
			Updates(map[string]interface{}{
				"status":      model.VerificationStatusVerified,
				"updated_at":  time.Now(),
				"verified_at": time.Now(),
			}).Error; err != nil {
			return err
		}
		// Update user's phone
		if err := tx.Model(&model.User{}).
			Where("id = ?", userID).
			Update("phone", newPhone).Error; err != nil {
			return err
		}
		return nil
	})
}

// // IsVerified checks if a user has completed all required verifications
// func (s *AuthServiceImpl) IsVerified(ctx context.Context, userID uuid.UUID) (bool, error) {
// 	// TODO: Check email verification : mostly common as default, but it can change based on issuer config setup
// 	emailVerified, err := s.repo.IsEmailVerified(ctx, userID)
// 	if err != nil {
// 		return false, fmt.Errorf("failed to check email verification: %w", err)
// 	}

// 	// TODO: Check phone verification : here needs to check a phone auth option and change implementation based on it
// 	phoneVerified, err := s.repo.IsPhoneVerified(ctx, userID)
// 	if err != nil {
// 		return false, fmt.Errorf("failed to check phone verification: %w", err)
// 	}

// 	// User is considered verified if both email and phone are verified
// 	return emailVerified && phoneVerified, nil
// }

// // ChangePassword implements password change
// func (s *AuthServiceImpl) ChangePassword(ctx context.Context, userID uuid.UUID, currentPassword, newPassword string) error {
// 	// Get user
// 	user, err := s.repo.GetUserByID(ctx, userID)
// 	if err != nil {
// 		return ErrUserNotFound
// 	}

// 	// Verify current password
// 	if err := user.VerifyPassword(currentPassword); err != nil {
// 		return ErrInvalidCredentials
// 	}

// 	// Hash new password
// 	hashedPassword, err := utils.EncryptPassword(newPassword, 10)
// 	if err != nil {
// 		return err
// 	}

// 	// Update password
// 	return s.repo.UpdatePassword(ctx, *user.ID, hashedPassword)
// }

// // SendVerificationPhone implements phone verification request
// func (s *AuthServiceImpl) SendVerificationPhone(ctx context.Context, phone string) error {
// 	user, err := s.repo.GetUserByPhone(ctx, phone)
// 	if err != nil {
// 		return ErrUserNotFound
// 	}

// 	// Generate verification code
// 	code := utils.GenerateVerificationCode()

// 	// Save code to database
// 	if err := s.repo.SavePhoneVerificationCode(ctx, *user.ID, code); err != nil {
// 		return err
// 	}

// 	// Send verification SMS
// 	return s.smsService.SendVerificationSMS(ctx, phone, code)
// }

// UpdatePassword implements password update
