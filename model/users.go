package model

// METHODS... Including Database queries and mutations...
import (
	"fmt"
	"time"

	uuid "github.com/google/uuid"
	"github.com/minilikmila/goAuth/db"
	"github.com/minilikmila/goAuth/utils"
	"gorm.io/gorm"
)

type User struct {
	// gorm.Model
	ID                           *uuid.UUID `gorm:"primary_key;unique" json:"id"`
	Name                         *string    `gorm:"type:varchar(25);not_null" json:"name"`
	Email                        *string    `gorm:"type:varchar(30);not_null;unique" json:"email"`
	Password                     *string    `gorm:"not_null" json:"-"`
	ProfilePicture               *string    `json:"profile_picture,omitempty"`
	Phone                        *string    `json:"phone,omitempty"`
	EmailConfirmed               bool       `json:"email_confirmed,omitempty"`
	EmailConfirmationTokenSentAt *time.Time `json:"email_confirmation_token_sent_at,omitempty"`
	EmailConfirmedAt             *time.Time `json:"email_confirmed_at,omitempty"`
	PhoneConfirmed               bool       `json:"phone_confirmed,omitempty"`
	PhoneConfirmationTokenSentAt *time.Time `json:"phone_confirmation_token_sent_at,omitempty"`
	PhoneConfirmedAt             *time.Time `json:"phone_confirmed_at,omitempty"`
	RecoveryTokenSentAt          *time.Time `json:"recovery_token_sent_at,omitempty"`
	EmailChangeTokenSentAt       *time.Time `json:"email_change_token_sent_at,omitempty"`
	PhoneChangeTokenSentAt       *time.Time `json:"phone_change_token_sent_at,omitempty"`
	// NewEmail                     *string    `json:"new_email,omitempty"`
	// NewPhone                     *string    `json:"new_phone,omitempty"`
	// PhoneChangedAt              *time.Time `json:"phone_changed_at,omitempty"`
	// EmailChangedAt              *time.Time `json:"email_changed_at,omitempty"`
	// PasswordChangedAt           *time.Time `json:"password_changed_at,omitempty"`
	IncorrectLoginAttempts      int        `json:"incorrect_login_attempts,omitempty"`
	LastIncorrectLoginAttemptAt *time.Time `json:"last_incorrect_login_attempt_at,omitempty"`
	EmailConfirmationToken      *string    `json:"-"`
	PhoneConfirmationToken      *string    `json:"-"`
	RecoveryToken               *string    `json:"-"`
	EmailChangeToken            *string    `json:"-"`
	PhoneChangeToken            *string    `json:"-"`
	LastLoginAt                 *time.Time `json:"last_login_at,omitempty"`
	CreatedAt                   time.Time  `json:"created_at,omitempty"`
	UpdatedAt                   time.Time  `json:"updated_at,omitempty"`
	Role                        string     `json:"role" default:"user"`
}

// create new user
func (u *User) Create(db *gorm.DB) error {
	u.CreatedAt = time.Now()
	u.ID = utils.GenerateUUID()
	hashedPassword, err := utils.EncryptPassword(*u.Password, 10)
	if err != nil {
		return err
	}
	u.Password = &hashedPassword
	fmt.Println("User ID : ", *u.ID)

	return db.Create(u).Error
}

// Update user data....
func (u *User) Save(db db.Database) error {

	return db.DB().Save(u).Error

}

// Invite user by email
func (u *User) InviteByEmail(db db.Database, email string, name string) error {
	return nil
}

// Invite user by phone
func (u *User) InviteByPhone(db db.Database, phone string, name string) error {
	return nil
}

// Compare password if not matched or if password is nil return err
func (u *User) VerifyPassword(password string) error {
	return nil
}

// Generate password and set u.Password and return error if found
func (u *User) SetPassword(password string, cost int) error {
	return nil
}

// Here we count user wrong attempts ( increment user.wrong_attempt++ field)
func (u *User) IncorrectAttempt(db gorm.DB, log Log) error {
	return nil
}
