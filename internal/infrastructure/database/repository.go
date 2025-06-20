package database

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/minilikmila/standard-auth-go/internal/domain/model"
	"gorm.io/gorm"
)

// Repository implements the Database interface
type Repository struct {
	db *gorm.DB
}

// NewRepository creates a new repository instance
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// User Operations
func (r *Repository) CreateUser(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *Repository) GetUserByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) GetUserByPhone(ctx context.Context, phone string) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).Where("phone = ?", phone).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) UpdateUser(ctx context.Context, id uuid.UUID, updates map[string]interface{}) error {
	return r.db.WithContext(ctx).Model(&model.User{}).Where("id = ?", id).Updates(updates).Error
}

func (r *Repository) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	return r.db.WithContext(ctx).Where("id = ?", userID).Delete(&model.User{}).Error
}

func (r *Repository) UserExists(ctx context.Context, email string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.User{}).Where("email = ?", email).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// Login Operations
func (r *Repository) IncrementLoginAttempts(ctx context.Context, userID uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&model.User{}).Where("id = ?", userID).
		Updates(map[string]interface{}{
			"incorrect_login_attempts":        gorm.Expr("incorrect_login_attempts + 1"),
			"last_incorrect_login_attempt_at": time.Now(),
		}).Error
}

func (r *Repository) ResetLoginAttempts(ctx context.Context, userID uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&model.User{}).Where("id = ?", userID).Update("incorrect_login_attempts", 0).Error
}

func (r *Repository) UpdateLastLogin(ctx context.Context, userID uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&model.User{}).Where("id = ?", userID).Update("last_login_at", time.Now()).Error
}

// Token Blacklist Operations
func (r *Repository) BlacklistToken(ctx context.Context, token string) error {
	blacklistedToken := &model.BlacklistedToken{
		Token:     token,
		ExpiresAt: time.Now().Add(24 * time.Hour),
		CreatedAt: time.Now(),
	}
	return r.db.WithContext(ctx).Create(blacklistedToken).Error
}

func (r *Repository) IsTokenBlacklisted(ctx context.Context, token string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&model.BlacklistedToken{}).
		Where("token = ? AND expires_at > ?", token, time.Now()).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// Verification Status Operations
func (r *Repository) IsEmailVerified(ctx context.Context, userID uuid.UUID) (bool, error) {
	return r.IsVerified(ctx, userID, model.VerificationTypeEmail)
}

func (r *Repository) IsPhoneVerified(ctx context.Context, userID uuid.UUID) (bool, error) {
	return r.IsVerified(ctx, userID, model.VerificationTypePhone)
}

// CreateVerification creates a new verification record
func (r *Repository) CreateVerification(ctx context.Context, verification *model.Verification) error {
	return r.db.WithContext(ctx).Create(verification).Error
}

// GetVerificationByToken retrieves a verification by token and type
func (r *Repository) GetVerificationByToken(ctx context.Context, token string, vType model.VerificationType) (*model.Verification, error) {
	var verification model.Verification
	err := r.db.WithContext(ctx).
		Where("token = ? AND type = ? AND status = ? AND expires_at > ?", token, vType, model.VerificationStatusPending, time.Now()).
		First(&verification).Error
	if err != nil {
		return nil, err
	}
	return &verification, nil
}

// GetVerificationByCode retrieves a verification by code and type
func (r *Repository) GetVerificationByCode(ctx context.Context, code string, vType model.VerificationType) (*model.Verification, error) {
	var verification model.Verification
	err := r.db.WithContext(ctx).
		Where("code = ? AND type = ? AND status = ? AND expires_at > ?", code, vType, model.VerificationStatusPending, time.Now()).
		First(&verification).Error
	if err != nil {
		return nil, err
	}
	return &verification, nil
}

// UpdateVerificationStatus updates the status of a verification
func (r *Repository) UpdateVerificationStatus(ctx context.Context, id uuid.UUID, status model.VerificationStatus) error {
	return r.db.WithContext(ctx).
		Model(&model.Verification{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":     status,
			"updated_at": time.Now(),
			"verified_at": func() *time.Time {
				if status == model.VerificationStatusVerified {
					now := time.Now()
					return &now
				}
				return nil
			}(),
		}).Error
}

// IsVerified checks if a user has a verified record for a given verification type
func (r *Repository) IsVerified(ctx context.Context, userID uuid.UUID, vType model.VerificationType) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&model.Verification{}).
		Where("user_id = ? AND type = ? AND status = ?", userID, vType, model.VerificationStatusVerified).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// Transaction executes a function within a database transaction
func (r *Repository) Transaction(ctx context.Context, fn func(tx *gorm.DB) error) error {
	return r.db.WithContext(ctx).Transaction(fn)
}

// UpdatePassword updates a user's password
func (r *Repository) UpdatePassword(ctx context.Context, userID uuid.UUID, hashedPassword string) error {
	return r.db.WithContext(ctx).Model(&model.User{}).Where("id = ?", userID).Updates(map[string]interface{}{
		"password":            hashedPassword,
		"password_changed_at": time.Now(),
	}).Error
}
