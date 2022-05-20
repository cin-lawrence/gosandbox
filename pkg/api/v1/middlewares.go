package v1

import (
	"errors"
	"net/http"

	"github.com/cin-lawrence/gosandbox/pkg/api/error"
	"github.com/cin-lawrence/gosandbox/pkg/services"

	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

var authService *services.AuthService = services.NewAuthService()

func ValidateToken(ctx *gin.Context) {
	am, err := authService.ExtractAccessMeta(ctx.Request)
	if err != nil {
		v1.SendError(ctx, http.StatusUnauthorized, errors.New("Invalid access token"))
		return
	}

	userID, err := authService.FetchAuth(am)
	if err != nil {
		v1.SendError(ctx, http.StatusUnauthorized, errors.New("Can't fetch auth"))
		return
	}
	ctx.Set("userID", userID)
}

func Authorize() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ValidateToken(ctx)
		ctx.Next()
	}

}

func CORS() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost")
		ctx.Writer.Header().Set("Access-Control-Max-Age", "86400")
		ctx.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		ctx.Writer.Header().Set("Access-Control-Allow-Headers", "X-Requested-With, Content-Type, Origin, Authorization, Accept, Client-Security-Token, Accept-Encoding, x-access-token")
		ctx.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
		ctx.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if ctx.Request.Method == "OPTIONS" {
			ctx.AbortWithStatus(200)
		} else {
			ctx.Next()
		}
	}
}

func GenerateRequestID() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestID := uuid.Must(uuid.NewV4())
		ctx.Writer.Header().Set("X-Request-Id", requestID.String())
		ctx.Next()
	}
}
