package v1

import (
        "errors"
	"net/http"

        "github.com/cin-lawrence/gosandbox/pkg/api/error"
        "github.com/cin-lawrence/gosandbox/pkg/services"

	"github.com/gin-gonic/gin"
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
