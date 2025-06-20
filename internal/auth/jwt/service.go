package jwt_

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	config "github.com/minilikmila/standard-auth-go/configs"
	"github.com/minilikmila/standard-auth-go/internal/auth/crypto"
	"github.com/minilikmila/standard-auth-go/internal/domain/model"
	"github.com/minilikmila/standard-auth-go/internal/domain/model/enum"
	"github.com/sirupsen/logrus"
)

// // JWTClaims represents the claims in a JWT token
// type JWTClaims struct {
// 	UserID string `json:"user_id"`
// 	Email  string `json:"email"`
// 	Role   string `json:"role"`
// 	jwt.RegisteredClaims
// }

// JWTService implements the JWT service interface
type JWTService struct {
	config *config.Config
}

// RedisClient interface for token blacklisting
type RedisClient interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Del(ctx context.Context, key string) error
}

// NewJWTService creates a new JWT service instance
func NewJWTService(config *config.Config) *JWTService {
	return &JWTService{
		config: config,
	}
}

// GenerateToken generates a new access token
func (s *JWTService) GenerateToken(user *model.User, metadata interface{}) (string, error) {
	logrus.Infoln("key : --- ", s.config.JWT.Alg)

	claims := JWTClaims{
		UserID:    user.ID.String(),
		Email:     user.Email,
		Role:      user.Role,
		TokenType: enum.AccessToken,
		Name:      user.Name,
		Metadata:  metadata,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(s.config.JWT.Exp) * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    s.config.JWT.Iss,
			Subject:   user.ID.String(),
			Audience:  jwt.ClaimStrings{s.config.JWT.Aud},
		},
	}

	token := jwt.NewWithClaims(jwt.GetSigningMethod(s.config.JWT.Alg), claims)

	switch jwt.GetSigningMethod(s.config.JWT.Alg) {
	case jwt.SigningMethodHS256, jwt.SigningMethodHS384, jwt.SigningMethodHS512:
		return token.SignedString([]byte(s.config.JWT.Secret))
	case jwt.SigningMethodRS256, jwt.SigningMethodRS384, jwt.SigningMethodRS512:
		privateKey, err := crypto.ParseRSAPrivateKeyFromPemString(s.config.JWT.GetSignKey().([]byte))
		if err != nil {
			return "", err
		}
		return token.SignedString(privateKey)
	case jwt.SigningMethodES256, jwt.SigningMethodES384, jwt.SigningMethodES512:
		privateKey, err := crypto.ParseECDSAPrivateKeyFromPemString(s.config.JWT.GetSignKey().([]byte))
		if err != nil {
			return "", err
		}
		return token.SignedString(privateKey)
	default:
		return "", errors.New("unsupported signing method")
	}
}

// GenerateRefreshToken generates a new refresh token
func (s *JWTService) GenerateRefreshToken(user *model.User) (string, error) {
	claims := JWTClaims{
		UserID:    user.ID.String(),
		TokenType: enum.RefreshToken,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour * 30)), // 30 days
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    s.config.JWT.Iss,
			Subject:   user.ID.String(),
			Audience:  jwt.ClaimStrings{s.config.JWT.Aud},
		},
	}

	token := jwt.NewWithClaims(jwt.GetSigningMethod(s.config.JWT.Alg), claims)

	switch jwt.GetSigningMethod(s.config.JWT.Alg) {
	case jwt.SigningMethodHS256, jwt.SigningMethodHS384, jwt.SigningMethodHS512:
		return token.SignedString([]byte(s.config.JWT.Secret))
	case jwt.SigningMethodRS256, jwt.SigningMethodRS384, jwt.SigningMethodRS512:
		privateKey, err := crypto.ParseRSAPrivateKeyFromPemString(s.config.JWT.GetSignKey().([]byte))
		if err != nil {
			return "", err
		}
		return token.SignedString(privateKey)
	case jwt.SigningMethodES256, jwt.SigningMethodES384, jwt.SigningMethodES512:
		privateKey, err := crypto.ParseECDSAPrivateKeyFromPemString(s.config.JWT.GetSignKey().([]byte))
		if err != nil {
			return "", err
		}
		return token.SignedString(privateKey)
	default:
		return "", errors.New("unsupported signing method")
	}
}

// GeneratePasswordResetToken generates a token for password reset
func (s *JWTService) GeneratePasswordResetToken(user *model.User) (string, error) {
	claims := JWTClaims{
		UserID:    user.ID.String(),
		Email:     user.Email,
		Role:      user.Role,
		TokenType: enum.PasswordResetToken,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)), // 1 hour
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    s.config.JWT.Iss,
			Subject:   user.ID.String(),
			Audience:  jwt.ClaimStrings{s.config.JWT.Aud},
		},
	}

	token := jwt.NewWithClaims(jwt.GetSigningMethod(s.config.JWT.Alg), claims)

	switch jwt.GetSigningMethod(s.config.JWT.Alg) {
	case jwt.SigningMethodHS256, jwt.SigningMethodHS384, jwt.SigningMethodHS512:
		return token.SignedString([]byte(s.config.JWT.Secret))
	case jwt.SigningMethodRS256, jwt.SigningMethodRS384, jwt.SigningMethodRS512:
		privateKey, err := crypto.ParseRSAPrivateKeyFromPemString(s.config.JWT.GetSignKey().([]byte))
		if err != nil {
			return "", err
		}
		return token.SignedString(privateKey)
	case jwt.SigningMethodES256, jwt.SigningMethodES384, jwt.SigningMethodES512:
		privateKey, err := crypto.ParseECDSAPrivateKeyFromPemString(s.config.JWT.GetSignKey().([]byte))
		if err != nil {
			return "", err
		}
		return token.SignedString(privateKey)
	default:
		return "", errors.New("unsupported signing method")
	}
}

// GenerateEmailVerificationToken generates a token for email verification
func (s *JWTService) GenerateEmailVerificationToken(user *model.User) (string, error) {
	claims := JWTClaims{
		UserID:    user.ID.String(),
		Email:     user.Email,
		Phone:     user.Phone,
		TokenType: enum.EmailVerificationToken,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // 24 hours
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    s.config.JWT.Iss,
			Subject:   user.ID.String(),
			Audience:  jwt.ClaimStrings{s.config.JWT.Aud},
		},
	}

	token := jwt.NewWithClaims(jwt.GetSigningMethod(s.config.JWT.Alg), claims)

	switch jwt.GetSigningMethod(s.config.JWT.Alg) {
	case jwt.SigningMethodHS256, jwt.SigningMethodHS384, jwt.SigningMethodHS512:
		return token.SignedString([]byte(s.config.JWT.Secret))
	case jwt.SigningMethodRS256, jwt.SigningMethodRS384, jwt.SigningMethodRS512:
		privateKey, err := crypto.ParseRSAPrivateKeyFromPemString(s.config.JWT.GetSignKey().([]byte))
		if err != nil {
			return "", err
		}
		return token.SignedString(privateKey)
	case jwt.SigningMethodES256, jwt.SigningMethodES384, jwt.SigningMethodES512:
		privateKey, err := crypto.ParseECDSAPrivateKeyFromPemString(s.config.JWT.GetSignKey().([]byte))
		if err != nil {
			return "", err
		}
		return token.SignedString(privateKey)
	default:
		return "", errors.New("unsupported signing method")
	}
}

// ValidateToken validates a JWT token and returns the user ID
func (s *JWTService) ValidateToken(token string) (uuid.UUID, error) {
	var err error
	var jwtToken *jwt.Token

	signingMethod := jwt.GetSigningMethod(s.config.JWT.Alg)
	if signingMethod == nil {
		return uuid.Nil, errors.New("unsupported signing method")
	}

	claims := &JWTClaims{}

	switch signingMethod {
	case jwt.SigningMethodHS256, jwt.SigningMethodHS384, jwt.SigningMethodHS512:
		jwtToken, err = jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
			return []byte(s.config.JWT.Secret), nil
		})
	case jwt.SigningMethodRS256, jwt.SigningMethodRS384, jwt.SigningMethodRS512:
		jwtToken, err = jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
			pubKey, err := crypto.ParseRSAPublicKeyFromPemString(s.config.JWT.GetDecodeKey().([]byte))
			if err != nil {
				return uuid.Nil, err
			}
			return pubKey, nil
		})
	case jwt.SigningMethodES256, jwt.SigningMethodES384, jwt.SigningMethodES512:
		jwtToken, err = jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
			pubKey, err := crypto.ParseECDSAPublicKeyFromPemString(s.config.JWT.GetDecodeKey().([]byte))
			if err != nil {
				return nil, err
			}
			return pubKey, nil
		})
	default:
		logrus.Infoln("Signin Method unsupported : --- ", signingMethod)
		return uuid.Nil, errors.New("unsupported signing method")
	}
	if err != nil {
		return uuid.Nil, err
	}
	parsedClaims, ok := jwtToken.Claims.(*JWTClaims)
	if !ok || !jwtToken.Valid {
		return uuid.Nil, errors.New("invalid token")
	}
	return uuid.Parse(parsedClaims.UserID)
}

// ValidateRefreshToken validates a refresh token
func (s *JWTService) ValidateRefreshToken(tokenString string) (uuid.UUID, error) {
	return s.ValidateToken(tokenString)
}

// ValidatePasswordResetToken validates a password reset token
func (s *JWTService) ValidatePasswordResetToken(tokenString string) (uuid.UUID, error) {
	return s.ValidateToken(tokenString)
}

// ValidateEmailVerificationToken validates an email verification token
func (s *JWTService) ValidateEmailVerificationToken(tokenString string) (uuid.UUID, error) {
	return s.ValidateToken(tokenString)
}

// InvalidateToken adds a token to the blacklist
func (s *JWTService) InvalidateToken(ctx context.Context, tokenString string) error {
	// TODO: Implement token blacklisting logic
	return nil
}
