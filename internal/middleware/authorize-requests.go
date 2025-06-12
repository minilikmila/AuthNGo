package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/minilikmila/standard-auth-go/internal/model"
)

// var user_ctx UserData

func AuthorizeRequests(allowed_roles []string) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		// check if user_role and request_role(get from header ctx) are matched
		if len(allowed_roles) == 0 {
			ctx.Next()
			return
		}
		// Now only read role from context

		ctx_user, err := CheckUserContext(ctx, "user_data")
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		// marshal body
		var userData model.UserCtxData

		if err := json.Unmarshal(ctx_user.([]byte), &userData); err != nil {
			fmt.Println("ctx un-marshaling error : ", err)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		// And then check if the requested role is allowed to access this requested action
		allowed := false
		for _, role := range allowed_roles {
			if strings.EqualFold(role, userData.Role) {
				allowed = true
				break
			}
		}

		if !allowed {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		ctx.Next()
	}
}
