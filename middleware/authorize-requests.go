package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthorizeRequests(allowed_roles []string) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		given_role := ctx.Request.Header.Get("x-user-role")
		if given_role == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		// check if user_role and request_role(get from header ctx) are matched
		user_role, exist := ctx.Get("role")
		if !exist || (given_role != user_role) {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		// And then check if the requested role is allowed to access this requested action
		allowed := false
		for _, role := range allowed_roles {
			if role == given_role {
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
