package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/minilikmila/standard-auth-go/internal/domain/model"
)

// Common errors
var (
	ErrUserNotFound       = errors.New("user not found")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidToken       = errors.New("invalid token")
	ErrInvalidCode        = errors.New("invalid verification code")
	ErrTokenExpired       = errors.New("token expired")
	ErrTooManyAttempts    = errors.New("too many login attempts")
	ErrEmailNotVerified   = errors.New("email not verified")
	ErrPhoneNotVerified   = errors.New("phone not verified")
)

// AuthService defines the interface for authentication operations
type AuthService interface {
	// User Management
	CreateUser(ctx context.Context, user *model.User) error
	UserExists(ctx context.Context, email string) (bool, error)
	Login(ctx context.Context, email, password string) (*model.User, error)
	Logout(ctx context.Context, token string) error
	GenerateTokens(ctx context.Context, user *model.User) (string, string, error)

	// Email Verification
	SendVerificationEmail(ctx context.Context, user *model.User) error
	VerifyEmail(ctx context.Context, token string) error

	// Phone Verification
	SendVerificationSMS(ctx context.Context, user *model.User) error
	VerifyPhone(ctx context.Context, code string, phone string) error

	// Password Management
	ForgotPassword(ctx context.Context, email string) error
	ResetPassword(ctx context.Context, token, newPassword string) error
	UpdatePassword(ctx context.Context, userID string, currentPassword string, newPassword string) error

	// Token Management
	RefreshToken(ctx context.Context, refreshToken string) (string, string, error)
	ValidateToken(ctx context.Context, token string) (string, error)

	// Profile Management
	UpdateProfile(ctx context.Context, userID string, updates map[string]interface{}) error
	GetProfile(ctx context.Context, userID string) (*model.User, error)
}

// Database defines the interface for database operations
type Database interface {
	// User operations
	CreateUser(ctx context.Context, user *model.User) error
	GetUserByID(ctx context.Context, id uuid.UUID) (*model.User, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	GetUserByPhone(ctx context.Context, phone string) (*model.User, error)
	UserExists(ctx context.Context, email string) (bool, error)
	UpdateUser(ctx context.Context, id uuid.UUID, updates map[string]interface{}) error
	DeleteUser(ctx context.Context, userID uuid.UUID) error

	// Password operations
	UpdatePassword(ctx context.Context, id uuid.UUID, hashedPassword string) error
	SavePasswordResetToken(ctx context.Context, id uuid.UUID, token string) error
	InvalidatePasswordResetToken(ctx context.Context, id uuid.UUID) error

	// Email verification
	SaveEmailVerificationToken(ctx context.Context, id uuid.UUID, token string) error
	InvalidateEmailVerificationToken(ctx context.Context, id uuid.UUID) error
	VerifyEmail(ctx context.Context, id uuid.UUID) error

	// Phone verification
	SavePhoneVerificationCode(ctx context.Context, id uuid.UUID, code string) error
	InvalidatePhoneVerificationCode(ctx context.Context, id uuid.UUID) error
	ValidatePhoneVerificationCode(ctx context.Context, code string) (uuid.UUID, error)
	VerifyPhone(ctx context.Context, id uuid.UUID) error

	// Login operations
	IncrementLoginAttempts(ctx context.Context, id uuid.UUID) error
	ResetLoginAttempts(ctx context.Context, id uuid.UUID) error
	UpdateLastLogin(ctx context.Context, id uuid.UUID) error

	// Token operations
	BlacklistToken(ctx context.Context, token string) error
	IsTokenBlacklisted(ctx context.Context, token string) (bool, error)
}

// JWTService interface for JWT operations
type JWTService interface {
	GenerateToken(user *model.User, metadata interface{}) (string, error)
	GenerateRefreshToken(user *model.User) (string, error)
	GeneratePasswordResetToken(user *model.User) (string, error)
	GenerateEmailVerificationToken(user *model.User) (string, error)
	ValidateToken(token string) (uuid.UUID, error)
	ValidateRefreshToken(token string) (uuid.UUID, error)
	ValidatePasswordResetToken(token string) (uuid.UUID, error)
	ValidateEmailVerificationToken(token string) (uuid.UUID, error)
	InvalidateToken(ctx context.Context, token string) error
}

// EmailService defines the interface for email operations
type EmailService interface {
	SendVerificationEmail(ctx context.Context, email, token, receiverName string) error
	SendPasswordResetEmail(ctx context.Context, email, token string) error
	SendWelcomeEmail(ctx context.Context, email string, name string) error
}

// SMSService defines the interface for SMS operations
type SMSService interface {
	SendVerificationSMS(ctx context.Context, phone, code string) error
	SendPasswordResetSMS(ctx context.Context, phone string, code string) error
}
