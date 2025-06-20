package model

// METHODS... Including Database queries and mutations...
import (
	"time"

	uuid "github.com/google/uuid"
)

type User struct {
	// gorm.Model
	ID                          *uuid.UUID `gorm:"primary_key;type:uuid;default:uuid_generate_v4()" json:"id"`
	Name                        *string    `gorm:"type:varchar(25);not_null" json:"name"`
	Email                       *string    `gorm:"type:varchar(30);not_null;unique" json:"email"`
	Password                    *string    `gorm:"not_null" json:"-"`
	ProfilePicture              *string    `json:"profile_picture,omitempty"`
	Phone                       *string    `json:"phone,omitempty"`
	SignUpMethod                string     `json:"signup_method" gorm:"type:varchar(20);default:'regular'"`
	IsEmailVerified             bool       `json:"is_email_verified" gorm:"default:false"`
	EmailVerifiedAt             *time.Time `json:"email_verified_at,omitempty"`
	IsPhoneVerified             bool       `json:"is_phone_verified" gorm:"default:false"`
	PhoneVerifiedAt             *time.Time `json:"phone_verified_at,omitempty"`
	NewEmail                    *string    `json:"new_email,omitempty"`
	NewPhone                    *string    `json:"new_phone,omitempty"`
	PhoneChangedAt              *time.Time `json:"phone_changed_at,omitempty"`
	EmailChangedAt              *time.Time `json:"email_changed_at,omitempty"`
	PasswordChangedAt           *time.Time `json:"password_changed_at,omitempty"`
	IncorrectLoginAttempts      int        `json:"incorrect_login_attempts,omitempty"`
	LastIncorrectLoginAttemptAt *time.Time `json:"last_incorrect_login_attempt_at,omitempty"`
	LastLoginAt                 *time.Time `json:"last_login_at,omitempty"`
	CreatedAt                   time.Time  `json:"created_at,omitempty"`
	UpdatedAt                   time.Time  `json:"updated_at,omitempty"`
	Role                        string     `json:"role" default:"user"`

	// Relationships
}

// IsVerified checks if the user has verified their email and phone
func (u *User) IsVerified() bool {
	// This would be implemented in the service layer
	// where we can check the verification status from the verification tables
	return u.IsEmailVerified
}

// - TableName method overridden
func (User) TableName() string {
	return "users"
}
