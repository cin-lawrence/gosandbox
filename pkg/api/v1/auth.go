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
	Validator   *validator.UserValidator
}

var v1AuthAPI *V1AuthAPI

func NewV1AuthGroup(rg *gin.RouterGroup) *gin.RouterGroup {
	v1AuthAPI = &V1AuthAPI{
		AuthService: *services.NewAuthService(),
		UserService: *services.NewUserService(),
		Validator:   new(validator.UserValidator),
	}

	v1AuthGroup := rg.Group("/auth")
	v1AuthGroup.POST("/login", v1AuthAPI.Login)
	v1AuthGroup.POST("/refresh", v1AuthAPI.RefreshToken)
	v1AuthGroup.POST("/logout", v1AuthAPI.Logout)
	return v1AuthGroup
}

// RefreshToken godoc
// @Summary	Refresh an access token
// @Tags	auth
// @Accept	json
// @Produce	json
// @Param	token	body models.RefreshTokenInput true "Refresh token"
// @Success	200	{object} models.Tokens
// @Failure	401	{object} error.APIError
// @Failure	403	{object} error.APIError
// @Failure	422	{object} error.APIError
// @Failure	500	{object} error.APIError
// @Router	/api/v1/auth/refresh/ [post]
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
		AccessToken:  tokenMeta.AccessToken,
		RefreshToken: tokenMeta.RefreshToken,
	}
	ctx.JSON(http.StatusOK, payload)
}

// Login godoc
// @Summary	Log in
// @Tags	auth
// @Accept	mpfd
// @Produce	json
// @Param	info	formData models.UserLogin true "Login information"
// @Success	200	{object} models.Tokens
// @Failure	404	{object} error.APIError
// @Failure	422	{object} error.APIError
// @Failure	500	{object} error.APIError
// @Router	/api/v1/auth/login/ [post]
func (api *V1AuthAPI) Login(ctx *gin.Context) {
	var userLogin models.UserLogin

	if err := ctx.ShouldBind(&userLogin); err != nil {
		message := api.Validator.Login(err)
		v1.SendError(ctx, http.StatusUnprocessableEntity, errors.New(message))
		return }
	user, err := api.UserService.GetByEmail(userLogin.Username)
	if err != nil {
		v1.SendError(ctx, http.StatusNotFound, fmt.Errorf("User %s does not exist", userLogin.Username))
		return
	}

	tokens, err := api.AuthService.Login(userLogin, user)
	if err != nil {
		v1.SendError(ctx, http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, tokens)
}

// Logout godoc
// @Summary	Log out
// @Tags	auth
// @Produce	json
// @Success	200	{string} string "OK"
// @Failure	400	{object} error.APIError
// @Failure	401	{object} error.APIError
// @Router	/api/v1/auth/logout/ [post]
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

	ctx.String(http.StatusOK, "OK")
}
