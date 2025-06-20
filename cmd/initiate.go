package cmd

import (
	"embed"
	"fmt"
	"net/http"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	config "github.com/minilikmila/standard-auth-go/configs"
	"github.com/minilikmila/standard-auth-go/internal/api/routes"
	jwt_ "github.com/minilikmila/standard-auth-go/internal/auth/jwt"
	"github.com/minilikmila/standard-auth-go/internal/infrastructure/database"
	"github.com/minilikmila/standard-auth-go/internal/service"
	"github.com/sirupsen/logrus"
)

func Run(files embed.FS, args map[string]string) {
	mode, ok := args["mode"]
	if !ok {
		mode = "debug"
	}

	// Load configuration
	cfg, err := config.New("config.json")
	if err != nil {
		logrus.Fatalf("Failed to load config: %v", err)
	}

	// Initialize Database connection
	db, err := database.InitDatabase(cfg)
	if err != nil {
		logrus.Fatalf("Failed to initialize database: %v", err)
	}

	// Initialize repository
	repo := database.NewRepository(db.DB())

	// Initialize JWT service
	jwtService := jwt_.NewJWTService(cfg)

	// Initialize email service (you'll need to implement this)
	emailService := service.NewEmailService(cfg)

	// Initialize SMS service (you'll need to implement this)
	smsService := service.NewSMSService(cfg)

	// Initialize auth service
	authService := service.NewAuthService(repo, cfg, emailService, smsService, jwtService)

	// Initialize routes with all services
	routes := routes.InitRoute(repo, cfg, mode, authService)

	host := fmt.Sprintf("%s:%v", cfg.Host, cfg.Port)
	logrus.WithField("host", "http://"+host).Info("Started Go authentication server")

	logrus.Fatalln(http.ListenAndServe(host, routes))
}
