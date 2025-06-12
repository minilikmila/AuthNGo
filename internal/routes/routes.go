package routes

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/minilikmila/standard-auth-go/config"
	database_ "github.com/minilikmila/standard-auth-go/internal/database"
	"github.com/minilikmila/standard-auth-go/internal/handlers"
	"github.com/minilikmila/standard-auth-go/internal/middleware"
	"github.com/minilikmila/standard-auth-go/internal/model"
)

func InitRoute(db database_.Database, config *config.Config, envMode string) *gin.Engine {
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
	v2.Use(middleware.ValidateAuthentication(config))

	// Use dependency injection
	v1.POST("/sign-up", handlers.SignUp(db))
	v1.POST("/login", handlers.Login(db))

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
