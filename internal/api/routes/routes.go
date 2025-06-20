package routes

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	config "github.com/minilikmila/standard-auth-go/configs"
	"github.com/minilikmila/standard-auth-go/internal/api/handlers"
	"github.com/minilikmila/standard-auth-go/internal/api/middleware"
	"github.com/minilikmila/standard-auth-go/internal/domain/model"
	database_ "github.com/minilikmila/standard-auth-go/internal/infrastructure/database"
	"github.com/minilikmila/standard-auth-go/internal/service"
)

func InitRoute(db *database_.Repository, config *config.Config, envMode string, authService service.AuthService) *gin.Engine {
	gin.SetMode(envMode)
	fmt.Println("Gin mode: ", gin.Mode())
	route := gin.New()
	route.Use(gin.Recovery())
	route.Use(middleware.AttachDeviceLog())
	route.Use(middleware.CorsMiddleware())

	// GROUPING ROUTES
	v1 := route.Group("/api/v1/auth")
	v2 := route.Group("/api/v1/")

	// MIDDLEWARE
	v2.Use(middleware.Authenticate(authService))

	// Auth routes
	v1.POST("/sign-up", handlers.SignUp(authService))
	v1.POST("/login", handlers.Login(authService))
	v1.POST("/logout", handlers.Logout(authService))
	// v1.POST("/refresh-token", handlers.RefreshToken(authService))
	v1.POST("/verify-email", handlers.VerifyEmail(authService))
	v1.POST("/verify-phone", handlers.VerifyPhone(authService))
	v1.POST("/forgot-password", handlers.ForgotPassword(authService))
	v1.POST("/reset-password", handlers.ResetPassword(authService))

	// Profile routes
	v2.GET("/profile/:id", middleware.Authorize([]string{"admin", "user"},
		nil,
		middleware.CustomAuthorizePolicy), handlers.GetProfile(authService))
	v2.PUT("/profile/:id", middleware.Authorize([]string{"admin", "user"},
		nil,
		middleware.CustomAuthorizePolicy), handlers.UpdateProfile(authService))
	// v2.PUT("/change-password", handlers.ChangePassword(authService))

	v1.GET("/health", checkHealth)

	return route
}

func checkHealth(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Healthy",
		"success": true,
	})
}

// This is example of wrapper function ... it can be use for used as a middleware but if your logic like authorization logic is highly coupled with the handler function. else middleware is recommended and always executed before the handler start its execution, we this enable us to manage the context.
func OnlyAdmin(fn gin.HandlerFunc) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user := getUserFromDB()
		if user.Role != "admin" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		fn(ctx)

	}
}
func getUserFromDB() model.User {
	return model.User{
		Role: "user",
	}
}
