package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/minilikmila/standard-auth-go/config"
	jwt_ "github.com/minilikmila/standard-auth-go/internal/jwt"
	"github.com/minilikmila/standard-auth-go/internal/model"
)

// Validate user has valid authentication signature...
func ValidateAuthentication(config *config.Config) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		token := ctx.GetHeader("Authorization")
		if token == "" {

			cookie, err := ctx.Cookie("go_auth_session")
			if err != nil {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing authentication header"})
				return
			}

			token = cookie

		}

		// Check the token if it's Bearer token ....
		tokenParts := strings.Split(token, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "malformed token"})
			return
		}

		// Verify the token
		token = tokenParts[1]
		fmt.Println("decoding")
		claims, err := jwt_.Decode(token, config)
		if err != nil {
			fmt.Println("decoding error :", err)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		// put bytes of user data in the context
		marshaled_ctx, err := json.Marshal(model.UserCtxData{
			ID:   claims.Subject,
			Role: claims.Role,
		})

		if err != nil {
			fmt.Println("decoding error :", err)
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}

		ctx.Set("user_data", marshaled_ctx)

		ctx.Next()
	}
}
