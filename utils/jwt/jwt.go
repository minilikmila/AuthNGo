package jwt_

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/minilikmila/goAuth/config"
	"github.com/minilikmila/goAuth/model"
	"github.com/minilikmila/goAuth/utils/crypto"
)

type JWTClaims struct {
	jwt.RegisteredClaims
	Name     *string      `json:"name"`
	Phone    *string      `json:"phone_number,omitempty"`
	Email    *string      `json:"email"`
	Metadata *interface{} `json:"metadata,omitempty"`
	Role     string       `json:"role"`
}

func New(user *model.User, metadata *interface{}) *JWTClaims {
	now := time.Now()
	//  j, _:=time.ParseDuration(config.Config.JWT.Exp)
	//  ex := time.Now().Add(j).uni

	return &JWTClaims{
		jwt.RegisteredClaims{
			Audience:  jwt.ClaimStrings{config.Config.JWT.Audience},
			Issuer:    config.Config.JWT.Iss,
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Minute * config.Config.JWT.Exp)),
			Subject:   user.ID.String(),
		},
		user.Name,
		user.Email,
		user.Phone,
		metadata,
		user.Role,
	}
}

func (j *JWTClaims) Sign() (token string, err error) {

	jwtAlg := config.Config.JWT.Alg

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
		return t.SignedString(config.Config.JWT.Secret)
	case jwt.SigningMethodRS256, jwt.SigningMethodRS384, jwt.SigningMethodRS512:
		privateKey, err := crypto.ParseRSAPrivateKeyFromPemString(config.Config.PrivateKey.([]byte))
		if err != nil {
			return "", err
		}
		return t.SignedString(privateKey)
	case jwt.SigningMethodES256, jwt.SigningMethodES384, jwt.SigningMethodES512:
		privateKey, err := crypto.ParseECDSAPrivateKeyFromPemString(config.Config.PrivateKey.([]byte))
		if err != nil {
			return "", err
		}
		return t.SignedString(privateKey)
	default:
		return "", errors.New("unsupported signing method")
	}
}

// Parse of verify JWT token
func Decode(token string) (*JWTClaims, error) {

	var err error
	var jwtToken *jwt.Token

	jwtAlg := config.Config.JWT.Alg

	signingMethod := jwt.GetSigningMethod(jwtAlg)
	if signingMethod == nil {
		return nil, errors.New("unsupported signing method")
	}

	var claims *JWTClaims

	switch signingMethod {
	case jwt.SigningMethodHS256, jwt.SigningMethodHS384, jwt.SigningMethodHS512:
		jwtToken, err = jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
			return []byte(config.Config.Secret), nil
		})
	case jwt.SigningMethodRS256, jwt.SigningMethodRS384, jwt.SigningMethodRS512:
		jwtToken, err = jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
			pubKey, err := crypto.ParseRSAPublicKeyFromPemString(config.Config.PublicKey.([]byte))
			if err != nil {
				return nil, err
			}
			return pubKey, nil
		})
	case jwt.SigningMethodES256, jwt.SigningMethodES384, jwt.SigningMethodES512:
		jwtToken, err = jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
			pubKey, err := crypto.ParseECDSAPublicKeyFromPemString(config.Config.PublicKey.([]byte))
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
