package v1

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/cin-lawrence/gosandbox/pkg/api/error"
	"github.com/cin-lawrence/gosandbox/pkg/models"
	"github.com/cin-lawrence/gosandbox/pkg/services"
	"github.com/cin-lawrence/gosandbox/pkg/validator"

	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v4"
)

type V1AuthAPI struct {
        AuthService services.AuthService
        UserService services.UserService
        Validator *validator.UserValidator
}

var v1AuthAPI *V1AuthAPI

func NewV1AuthGroup(rg *gin.RouterGroup) *gin.RouterGroup {
	v1AuthAPI = &V1AuthAPI{
                AuthService: *services.NewAuthService(),
		UserService: *services.NewUserService(),
                Validator: new(validator.UserValidator),
	}

	v1AuthGroup := rg.Group("/auth")
	v1AuthGroup.POST("/login", v1AuthAPI.Login)
	return v1AuthGroup
}

func (api *V1AuthAPI) RefreshToken(ctx *gin.Context) {
        var refreshTokenInput models.RefreshTokenInput

        if ctx.ShouldBindJSON(&refreshTokenInput) != nil {
                v1.SendError(ctx, http.StatusUnprocessableEntity, errors.New("Invalid form"))
                return
        }

        token, err := jwt.Parse(refreshTokenInput.RefreshToken, api.AuthService.GetSecretIfValid)
        if err != nil {
                v1.SendError(ctx, http.StatusUnauthorized, errors.New("Invalid authorization, please login again"))
                return
        }

        claims, ok := token.Claims.(jwt.MapClaims)
        if !(ok && token.Valid) {
                v1.SendError(ctx, http.StatusUnauthorized, errors.New("Invalid authorization, please login again"))
                return
        }

        refreshUUID, ok := claims["refresh_uuid"].(string)
        if !ok {
                v1.SendError(ctx, http.StatusUnauthorized, errors.New("Invalid refresh UUID"))
                return
        }

        userID64, err := strconv.ParseUint(fmt.Sprintf("%.f", claims["user_id"]), 10, 64)
        if err != nil {
                v1.SendError(ctx, http.StatusUnauthorized, errors.New("Invalid user ID"))
                return
        }
        userID := uint(userID64)

        deleted, err := api.AuthService.DeleteAuth(refreshUUID)
        if err != nil || deleted == 0 {
                v1.SendError(ctx, http.StatusUnauthorized, errors.New("Refresh token doesn't exist"))
                return
        }

        tokenMeta, err := api.AuthService.CreateToken(userID)
        if err != nil {
                v1.SendError(ctx, http.StatusForbidden, errors.New("Can't create new token"))
                return
        }

        err = api.AuthService.CreateAuth(userID, tokenMeta)
        if err != nil {
                v1.SendError(ctx, http.StatusForbidden, errors.New("Can't save new Auth session"))
                return
        }

        payload := models.Tokens{
                AccessToken: tokenMeta.AccessToken,
                RefreshToken: tokenMeta.RefreshToken,
        }
        ctx.JSON(http.StatusOK, payload)
}

func (api *V1AuthAPI) Login(ctx *gin.Context) {
        var userLogin models.UserLogin

        if err := ctx.ShouldBind(&userLogin); err != nil {
                message := api.Validator.Login(err)
                v1.SendError(ctx, http.StatusUnprocessableEntity, errors.New(message))
                return
        }

        user, err := api.UserService.GetByEmail(userLogin.Email)
        if err != nil {
                v1.SendError(ctx, http.StatusNotFound, fmt.Errorf("User %s does not exist", userLogin.Email))
                return
        }

        tokens, err := api.AuthService.Login(userLogin, user)
        if err != nil {
                v1.SendError(ctx, http.StatusBadRequest, err)
                return
        }
        ctx.JSON(http.StatusOK, tokens)
}

func (api *V1AuthAPI) Logout(ctx *gin.Context) {
        accessMeta, err := api.AuthService.ExtractAccessMeta(ctx.Request)
        if err != nil {
                v1.SendError(ctx, http.StatusBadRequest, errors.New("User not logged in"))
                return
        }

        nDeleted, err := api.AuthService.DeleteAuth(accessMeta.AccessUUID.String())
        if err != nil || nDeleted == 0 {
                v1.SendError(ctx, http.StatusUnauthorized, errors.New("Invalid request"))
                return
        }

        ctx.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"})
}
