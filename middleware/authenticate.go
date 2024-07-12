package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	jwt_ "github.com/minilikmila/goAuth/utils/jwt"
)

// Validate user has valid authentication signature...
func ValidateAuthentication() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// fmt.Println("Header : ", ctx.Request.RemoteAddr)

		token := ctx.GetHeader("Authorization")
		if token == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing authentication token"})
			return
		}

		// Check the token if it's Bearer token ....
		tokenParts := strings.Split(token, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "malformed token"})
			return
		}

		// Verify the token
		token = tokenParts[1]

		claims, err := jwt_.Decode(token)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		ctx.Set("role", claims.Role) // this may be array of string ...think
		ctx.Set("user_id", claims.ID)

		ctx.Next()
	}
}
