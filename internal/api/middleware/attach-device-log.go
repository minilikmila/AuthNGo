package middleware

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/minilikmila/standard-auth-go/internal/domain/model"

	"github.com/sirupsen/logrus"
)

func AttachDeviceLog() gin.HandlerFunc {
	return func(c *gin.Context) {

		// this one take care of x-real-ip, x-forwarded-for and RemoteAddress ... gin
		// ctx.Request = ctx.Request.WithContext(context.WithValue(context.Background(), "provider", provider))
		marshaledData, err := json.Marshal(model.LogData{
			ClientIp:          c.ClientIP(),
			UserAgent:         c.Request.UserAgent(),
			RequestedResource: c.Request.URL.Path,
			At:                time.Now(),
		})
		if err != nil {
			logrus.Error("error: attach device log error :", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err})
			return
		}
		// Log info
		logrus.Infoln("Device fingerprint : ", model.LogData{
			ClientIp:          c.ClientIP(),
			UserAgent:         c.Request.UserAgent(),
			RequestedResource: c.Request.URL.Path,
			At:                time.Now(),
		})
		c.Set("log_data", marshaledData)

		c.Next()
		// c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Forbidden!"})

	}
}
