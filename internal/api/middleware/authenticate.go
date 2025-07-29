package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/minilikmila/standard-auth-go/internal/domain/model"
	"github.com/minilikmila/standard-auth-go/internal/service"
	"github.com/minilikmila/standard-auth-go/pkg/utils"
)

// Validate user has valid authentication signature...
func Authenticate(authService service.AuthService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var accessToken string
		header := ctx.GetHeader("Authorization")
		if header != "" {
			accessToken = utils.ExtractTokenFromHeader(header)
		} else if cookie, err := ctx.Cookie("access_token"); err == nil {
			accessToken = cookie // or utils.ExtractTokenFromCookie(cookie) if you re-add it
		}
		if accessToken == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing or malformed token"})
			return
		}

		// Validate token
		userId, err := authService.ValidateToken(ctx, accessToken)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token!"})
			return
		}

		user, err := authService.GetProfile(ctx, userId)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		// TODO: validate token type is "access_token" not refresh
		// TODO: compare the attached token in the header role with user actual token and to reject unauthorized access.
		// TODO: decode and extract token

		// put bytes of user data in the context
		marshaled_ctx, err := json.Marshal(model.UserCtxData{
			ID:   userId,
			Role: user.Role,
		})
		logrus.Infoln("User data : ", model.UserCtxData{
			ID:   userId,
			Role: user.Role,
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
