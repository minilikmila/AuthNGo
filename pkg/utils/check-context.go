package utils

import (
	"errors"

	"github.com/gin-gonic/gin"
)

func CheckUserContext(ctx *gin.Context, name string) (any, error) {
	ctx_user, exist := ctx.Get("user_data")

	if !exist {
		return nil, errors.New("error: context not found")
	}
	return ctx_user, nil
}
