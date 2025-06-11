package config

import (
	"encoding/json"
	"os"
	"time"

	"github.com/minilikmila/standard-auth-go/internal/crypto"
)

type JWTConfig struct {
	Aud            string        `json:"audience"`
	Alg            string        `json:"algorithm"`
	Exp            time.Duration `json:"expiry"`
	Iss            string        `json:"issuer"`
	PrivateKeyPath string        `json:"private_key_path"`
	PublicKeyPath  string        `json:"public_key_path"`
	Secret         string        `json:"secret"`
	Type           string        `json:"-"`
	privateKey     interface{}
	publicKey      interface{}
}

type LockoutPolicy struct {
	Attempts int           `json:"attempts"`
	For      time.Duration `json:"for"`
}

type SocialAuthConfig struct {
	ClientID string `json:"client_id"`
	SecretID string `json:"secret_id"`
	Enabled  bool   `json:"enabled"`
}

type Config struct {
	Host                     string           `json:"host"`
	Port                     int              `json:"port"`
	DatabaseUri              string           `json:"database_uri"`
	InstanceUrl              string           `json:"instance_url"`
	DisableSignup            bool             `json:"disable_signup"`
	DisableEmail             bool             `json:"disable_email"`
	DisablePhone             bool             `json:"disable_phone"`
	GeneieApiKey             string           `json:"genie_api_key"`
	CloudinaryCloudName      string           `json:"cloudinary_cloud_name"`
	CloudinaryKey            string           `json:"cloudinary_key"`
	CloudinarySecret         string           `json:"cloudinary_secret"`
	AccessTokenCookieName    string           `json:"access_token_cookie_name"`
	AccessTokenCookieDomain  string           `json:"access_token_cookie_domain"`
	RefreshTokenCookieName   string           `json:"refresh_token_cookie_name"`
	RefreshTokenCookieDomain string           `json:"refresh_token_cookie_domain"`
	SessionCookieName        string           `json:"session_cookie_name"`
	SessionCookieDomain      string           `json:"session_cookie_domain"`
	JWT                      JWTConfig        `json:"jwt"`
	Google                   SocialAuthConfig `json:"google"`
	Github                   SocialAuthConfig `json:"github"`
	Linkedin                 SocialAuthConfig `json:"linkedin"`
	Facebook                 SocialAuthConfig `json:"facebook"`
	Apple                    SocialAuthConfig `json:"apple"`
	Twitter                  SocialAuthConfig `json:"twitter"`
	Slack                    SocialAuthConfig `json:"slack"`
	Discord                  SocialAuthConfig `json:"discord"`
	SocialAuthRedirectUrl    string           `json:"social_auth_redirect_url"`
	MaxConnectionPoolSize    int              `json:"max_connection_pool_size"`
	LockoutPolicy            LockoutPolicy    `json:"lockout_policy"`
	AdminRoles               []string         `json:"admin_roles"`
	ReadOnlyRoles            []string         `json:"read_only_roles"`
}

func DefaultConfig() *Config {
	return &Config{
		Host:                  "localhost",
		Port:                  3001,
		DatabaseUri:           "",
		MaxConnectionPoolSize: 10,
		DisableSignup:         false,
		DisableEmail:          false,
		DisablePhone:          true,
		SessionCookieName:     "genie_session",
		JWT: JWTConfig{
			Exp: 2000,
			Alg: "RS512",
		},
		LockoutPolicy: LockoutPolicy{
			Attempts: 10,
			For:      60,
		},
		Google:                SocialAuthConfig{Enabled: false},
		Github:                SocialAuthConfig{Enabled: false},
		Linkedin:              SocialAuthConfig{Enabled: false},
		Facebook:              SocialAuthConfig{Enabled: false},
		Apple:                 SocialAuthConfig{Enabled: false},
		Twitter:               SocialAuthConfig{Enabled: false},
		Slack:                 SocialAuthConfig{Enabled: false},
		Discord:               SocialAuthConfig{Enabled: false},
		SocialAuthRedirectUrl: "http://localhost:5173",
	}
}

func New(path string) (*Config, error) {
	config := DefaultConfig()

	content, err := os.ReadFile(path)

	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(content, &config); err != nil {
		return nil, err
	}

	if err = validateSocial(config); err != nil {
		return nil, err
	}

	if err = validateCommon(config); err != nil {
		return nil, err
	}

	if config.JWT.privateKey, config.JWT.publicKey, err = crypto.ReadKeys(config.JWT.PrivateKeyPath, config.JWT.PublicKeyPath); err != nil {
		return nil, err
	}

	return config, nil
}

func validateSocial(config *Config) error {

	if config.Google.Enabled && (config.Google.ClientID == "" || config.Google.SecretID == "") {
		return ErrGoogleConfig
	}

	if config.Linkedin.Enabled && (config.Linkedin.ClientID == "" || config.Linkedin.SecretID == "") {
		return ErrLinkedinConfig
	}

	if config.Facebook.Enabled && (config.Facebook.ClientID == "" || config.Facebook.SecretID == "") {
		return ErrFacebookConfig
	}

	if config.Apple.Enabled && (config.Apple.ClientID == "" || config.Apple.SecretID == "") {
		return ErrAppleConfig
	}
	// ...will continue
	return nil
}

func validateCommon(config *Config) error {

	if config.DatabaseUri == "" {
		return ErrDatabaseURIRequired
	}

	if config.DisableEmail && config.DisablePhone {
		return ErrPhoneEmailDisabled
	}

	return nil
}

func (j *JWTConfig) GetSignKey() interface{} {
	return j.privateKey
}

func (j *JWTConfig) GetDecodeKey() interface{} {
	return j.publicKey
}
