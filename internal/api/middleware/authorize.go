package middleware

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/minilikmila/standard-auth-go/internal/domain/model"
	"github.com/minilikmila/standard-auth-go/pkg/utils"
)

func Authorize(
	allowedRoles []string,
	requiredPermissions []string,
	customPolicy func(user model.UserCtxData, ctx *gin.Context) bool,
) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctxUser, err := utils.CheckUserContext(ctx, "user_data")
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		var userData model.UserCtxData
		userBytes, ok := ctxUser.([]byte)
		if !ok || json.Unmarshal(userBytes, &userData) != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		// Role check
		if len(allowedRoles) > 0 {
			allowed := false
			for _, role := range allowedRoles {
				if strings.EqualFold(role, userData.Role) {
					allowed = true
					break
				}
			}
			if !allowed {
				ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden: role"})
				return
			}
		}

		// Permission check
		if len(requiredPermissions) > 0 {
			permMap := make(map[string]struct{}, len(userData.Permissions))
			for _, p := range userData.Permissions {
				permMap[p] = struct{}{}
			}
			for _, rp := range requiredPermissions {
				if _, ok := permMap[rp]; !ok {
					ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden: permission"})
					return
				}
			}
		}

		// Custom policy check
		if customPolicy != nil && !customPolicy(userData, ctx) {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden: policy"})
			return
		}

		ctx.Next()
	}
}

// UserOrAdminPolicy allows admin to access any resource, and users to access only their own resource (by id param)
func CustomAuthorizePolicy(user model.UserCtxData, ctx *gin.Context) bool {
	if strings.EqualFold(user.Role, "admin") {
		return true
	}
	idParam := ctx.Param("id")
	return strings.EqualFold(user.Role, "user") && idParam == user.ID
}
