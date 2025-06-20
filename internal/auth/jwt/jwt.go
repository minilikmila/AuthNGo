package jwt_

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"

	config "github.com/minilikmila/standard-auth-go/configs"
	"github.com/minilikmila/standard-auth-go/internal/auth/crypto"
	"github.com/minilikmila/standard-auth-go/internal/domain/model"
)

type JWTClaims struct {
	jwt.RegisteredClaims
	UserID    string      `json:"user_id"`
	Name      *string     `json:"name"`
	Phone     *string     `json:"phone_number,omitempty"`
	Email     *string     `json:"email"`
	Metadata  interface{} `json:"metadata,omitempty"`
	Role      string      `json:"role"`
	TokenType string      `json:"token_type"`
	config    *config.Config
}

func New(user *model.User, metadata interface{}, config *config.Config) *JWTClaims {
	now := time.Now()
	//  j, _:=time.ParseDuration(config.Config.JWT.Exp)
	//  ex := time.Now().Add(j).uni

	return &JWTClaims{
		jwt.RegisteredClaims{
			Audience:  jwt.ClaimStrings{config.JWT.Aud},
			Issuer:    config.JWT.Iss,
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Minute * config.JWT.Exp)),
			Subject:   user.ID.String(),
		},
		user.ID.String(),
		user.Name,
		user.Email,
		user.Phone,
		metadata,
		user.Role,
		"access",
		config,
	}
}

func (j *JWTClaims) Sign() (token string, err error) {

	jwtAlg := j.config.JWT.Alg

	signMethod := jwt.GetSigningMethod(jwtAlg)
	if signMethod == nil {
		return "", errors.New("unsupported signing method")
	}

	t := jwt.NewWithClaims(signMethod, j) // we can use jwt.New(alg) and we can bind the claims later.
	if t == nil {
		return "", errors.New("unsupported signing method")
	}

	switch signMethod {
	case jwt.SigningMethodHS256, jwt.SigningMethodHS384, jwt.SigningMethodHS512:
		return t.SignedString(j.config.JWT.Secret)
	case jwt.SigningMethodRS256, jwt.SigningMethodRS384, jwt.SigningMethodRS512:
		privateKey, err := crypto.ParseRSAPrivateKeyFromPemString(j.config.JWT.GetSignKey().([]byte))
		if err != nil {
			return "", err
		}
		return t.SignedString(privateKey)
	case jwt.SigningMethodES256, jwt.SigningMethodES384, jwt.SigningMethodES512:
		privateKey, err := crypto.ParseECDSAPrivateKeyFromPemString(j.config.JWT.GetSignKey().([]byte))
		if err != nil {
			return "", err
		}
		return t.SignedString(privateKey)
	default:
		return "", errors.New("unsupported signing method")
	}
}

// Parse of verify JWT token
func Decode(token string, config *config.Config) (*JWTClaims, error) {

	var err error
	var jwtToken *jwt.Token

	jwtAlg := config.JWT.Alg

	signingMethod := jwt.GetSigningMethod(jwtAlg)
	if signingMethod == nil {
		return nil, errors.New("unsupported signing method")
	}

	var claims *JWTClaims

	switch signingMethod {
	case jwt.SigningMethodHS256, jwt.SigningMethodHS384, jwt.SigningMethodHS512:
		jwtToken, err = jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
			return []byte(config.JWT.Secret), nil
		})
	case jwt.SigningMethodRS256, jwt.SigningMethodRS384, jwt.SigningMethodRS512:
		jwtToken, err = jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
			pubKey, err := crypto.ParseRSAPublicKeyFromPemString(config.JWT.GetDecodeKey().([]byte))
			if err != nil {
				return nil, err
			}
			return pubKey, nil
		})
	case jwt.SigningMethodES256, jwt.SigningMethodES384, jwt.SigningMethodES512:
		jwtToken, err = jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
			pubKey, err := crypto.ParseECDSAPublicKeyFromPemString(config.JWT.GetDecodeKey().([]byte))
			if err != nil {
				return nil, err
			}
			return pubKey, nil
		})
	default:
		return nil, errors.New("unsupported signing method")
	}
	if err != nil {
		return nil, err
	}
	claims, ok := jwtToken.Claims.(*JWTClaims)
	if !ok || !jwtToken.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
