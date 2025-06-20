package model

import (
	"time"

	"github.com/google/uuid"
)

type VerificationType string

const (
	VerificationTypeEmail VerificationType = "email"
	VerificationTypePhone VerificationType = "phone"
	VerificationTypeReset VerificationType = "reset"
	VerificationType2FA   VerificationType = "2fa"
)

// VerificationStatus represents the status of a verification
type VerificationStatus string

const (
	VerificationStatusPending  VerificationStatus = "pending"
	VerificationStatusVerified VerificationStatus = "verified"
	VerificationStatusExpired  VerificationStatus = "expired"
	VerificationStatusFailed   VerificationStatus = "failed"
)

// This model is used for revocation and one-time use purpose for stateless token ( use combined approach, which is secure and recommended approach)

type Verification struct {
	ID         uuid.UUID          `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	UserID     uuid.UUID          `json:"user_id" gorm:"type:uuid;not null"`
	Type       VerificationType   `json:"type" gorm:"type:varchar(20);not null"`
	Token      *string            `json:"token,omitempty" gorm:"type:varchar(255);uniqueIndex"`
	Code       *string            `json:"code,omitempty" gorm:"type:varchar(20);uniqueIndex"`
	Status     VerificationStatus `json:"status" gorm:"type:varchar(20);not null;default:'pending'"`
	SentAt     time.Time          `json:"sent_at" gorm:"not null"`
	VerifiedAt *time.Time         `json:"verified_at,omitempty"`
	ExpiresAt  time.Time          `json:"expires_at" gorm:"not null"`
	CreatedAt  time.Time          `json:"created_at" gorm:"not null"`
	UpdatedAt  time.Time          `json:"updated_at" gorm:"not null"`
}

// BlacklistedToken represents a token that has been invalidated
type BlacklistedToken struct {
	Token     string    `gorm:"primary_key" json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at" gorm:"not null"`
}

// Override TableName() method in Gorm with our own. then this TableName specifies the table name for each verification type
func (Verification) TableName() string {
	return "verifications"
}

func (BlacklistedToken) TableName() string {
	return "blacklisted_tokens"
}
