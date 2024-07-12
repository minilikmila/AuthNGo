package config

import (
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/minilikmila/goAuth/model/enum"
	"github.com/minilikmila/goAuth/utils/crypto"
	log "github.com/sirupsen/logrus"
)

type JWT struct {
	SecretKey      string
	PrivateKeyPath string
	PublicKeyPath  string
	Alg            string
	Audience       string
	Exp            time.Duration
	Iss            string
	PrivateKey     interface{}
	PublicKey      interface{}
	Secret         string
	Type           string `json:"-"`
}
type AppConfig struct {
	App struct {
		Host string
		Port string
	}
	JWT `json:"jwt"`
	DB  struct {
		Host     string
		Port     string
		User     string
		Password string
		Name     string
		SSLMode  string `default:"disabled"`
		DSN      string
	}
}

var Config AppConfig
var err error

func init() {
	// we can use exported func and deal with devMod and pass the value here and do logics
	if err := godotenv.Load(); err != nil {
		log.Fatal("error loading .env file")
	}
	Config.App.Host = os.Getenv("Host")
	Config.App.Port = os.Getenv("PORT")

	Config.JWT.SecretKey = os.Getenv("JWT_SECRET")
	Config.JWT.PrivateKeyPath = os.Getenv("PRIVATE_KEY_PATH")
	Config.JWT.PublicKeyPath = os.Getenv("PUBLIC_KEY_PATH")
	Config.JWT.Alg = os.Getenv("JWT_SIGNING_ALGORITHM")
	Config.JWT.Type = os.Getenv("TYPE")

	Config.DB.Host = os.Getenv("DB_HOST")
	Config.DB.Port = os.Getenv("DB_PORT")
	Config.DB.User = os.Getenv("DB_USER")
	Config.DB.Password = os.Getenv("DB_PASSWORD")
	Config.DB.Name = os.Getenv("DB_NAME")
	Config.DB.SSLMode = os.Getenv("DB_SSL_MODE")
	Config.DB.DSN = os.Getenv("DSN")
	// Config.DB.DSNUrl = os.Getenv("DSN_URL")

	if strings.ToLower(Config.Type) == enum.Asymmetric {
		_, ok := os.LookupEnv("PRIVATE_KEY_PATH")
		_, ok2 := os.LookupEnv("PUBLIC_KEY_PATH")
		if !ok || !ok2 {
			log.Fatalln("jwt secret must be provided")
		}

		Config.PrivateKey, Config.PublicKey, err = crypto.ReadKeys(Config.PrivateKeyPath, Config.PublicKeyPath)
		if err != nil {
			log.Fatalln("error ocurred in loading jwt keys")
		}
	} else {
		_, ok := os.LookupEnv("JWT_SECRET")
		if !ok {
			log.Fatalln("jwt secret must be provided")
		}
	}
}

func NewDefaultConfig() *AppConfig {
	return &AppConfig{}
}

func NewConfig() *AppConfig {
	return &AppConfig{}
}

func (c *AppConfig) validateKeys() error {

	return nil
}

// func (j *JWT) GetSigningKey() interface{} {
// 	if j.Type == "asymmetric" {
// 		return j.PrivateKey
// 	}
// 	return j.Secret
// }

// func (j *JWT) GetDecodingKey() interface{} {
// 	if j.Type == "asymmetric" {
// 		return j.PublicKey
// 	}
// 	return j.Secret
// }
