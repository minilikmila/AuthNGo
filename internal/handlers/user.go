package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	database_ "github.com/minilikmila/standard-auth-go/internal/database"
	"github.com/minilikmila/standard-auth-go/internal/model"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func SignUp(db database_.Database) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		dbInstance := db.DB()

		var body = model.SignUpForm{}

		if err := ctx.BindJSON(&body); err != nil {
			ctx.JSON(http.StatusInternalServerError, model.Response{
				Message:    "invalid request",
				StatusCode: http.StatusInternalServerError,
				Data:       nil,
				Error:      err,
			})
			return
		}

		fmt.Println("signup form ", body)

		user := &model.User{
			Name:     &body.Name,
			Phone:    &body.Phone,
			Email:    &body.Email,
			Password: &body.Password,
		}
		fmt.Println("User form ", user)
		if err := user.Create(dbInstance); err != nil {
			log.Errorln("error : ", err)
			ctx.JSON(http.StatusInternalServerError, model.Response{
				Message:    "error encountered when sign up",
				Error:      err,
				StatusCode: http.StatusInternalServerError,
				Data:       nil,
			})
			return
		}

		// u.Create(db.DB)

		ctx.JSON(http.StatusOK, gin.H{
			"message": "successful",
		})
	}
}

func Login(db database_.Database) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "successful",
		})
	}
}

func ForgotPassword(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {

	}
}

func ResetPassword(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {}
}
