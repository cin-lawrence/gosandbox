package v1

import (
	"github.com/gin-gonic/gin"
)

type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func SendError(ctx *gin.Context, status int, err error) {
	apiError := APIError{
		Code:    status,
		Message: err.Error(),
	}
	ctx.AbortWithStatusJSON(status, apiError)
}
