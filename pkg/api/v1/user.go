package v1

import (
	"errors"
	"net/http"

	"github.com/cin-lawrence/gosandbox/pkg/api/error"
	"github.com/cin-lawrence/gosandbox/pkg/models"
	"github.com/cin-lawrence/gosandbox/pkg/services"
	"github.com/cin-lawrence/gosandbox/pkg/validator"
	"github.com/gin-gonic/gin"
)

type V1UserAPI struct {
	UserService services.UserService
	Validator   *validator.UserValidator
}

var v1UserAPI *V1UserAPI

func NewV1UserGroup(rg *gin.RouterGroup) *gin.RouterGroup {
	v1UserAPI = &V1UserAPI{
		UserService: *services.NewUserService(),
		Validator:   new(validator.UserValidator),
	}

	// v1UserGroup := rg.Group("/users", Authorize())
	v1UserGroup := rg.Group("/users")
	v1UserGroup.GET("/", v1UserAPI.ListUsers)
	v1UserGroup.POST("/", v1UserAPI.CreateUser)
	v1UserGroup.GET("/:id", v1UserAPI.GetUser)
	v1UserGroup.PUT("/:id", v1UserAPI.UpdateUser)
	v1UserGroup.DELETE("/:id", v1UserAPI.DeleteUser)
	return v1UserGroup
}

// ListUsers godoc
// @Summary	List all users
// @Tags	users
// @Produce	json
// @Success	200	{object} models.UserList
// @Failure	500	{object} error.APIError
// @Router	/api/v1/users/ [get]
// @Security    OAuth2Password
func (api *V1UserAPI) ListUsers(ctx *gin.Context) {
	users, err := api.UserService.List()
	if err != nil {
		v1.SendError(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, models.UserList{Items: users})
}

// CreateUsers godoc
// @Summary	Create a new user
// @Tags	users
// @Accept	mpfd
// @Produce	json
// @Param	info	formData models.UserInput true "User info"
// @Success	201	{object} models.User
// @Failure	422	{object} error.APIError
// @Failure	500	{object} error.APIError
// @Router	/api/v1/users/ [post]
// @Security    OAuth2Password
func (api *V1UserAPI) CreateUser(ctx *gin.Context) {
	var userInput models.UserInput
	if err := ctx.ShouldBind(&userInput); err != nil {
		message := api.Validator.CreateUser(err)
		v1.SendError(ctx, http.StatusUnprocessableEntity, errors.New(message))
		return
	}

	user, err := api.UserService.Create(userInput)
	if err != nil {
		v1.SendError(ctx, http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusCreated, user)
}

// GetUser godoc
// @Summary	Retrieve a user information
// @Tags	users
// @Produce	json
// @Param	id	path integer true "User ID"
// @Success	200	{object} models.User
// @Failure	404	{object} error.APIError
// @Router	/api/v1/users/{id} [get]
// @Security    OAuth2Password
func (api *V1UserAPI) GetUser(ctx *gin.Context) {
	id := ctx.Param("id")
	user, err := api.UserService.Get(id)
	if err != nil {
		v1.SendError(ctx, http.StatusNotFound, err)
		return
	}

	ctx.JSON(http.StatusOK, user)
}

// UpdateUser godoc
// @Summary	Update a user
// @Tags	users
// @Accept	json
// @Produce	json
// @Param	id	path integer true "User ID"
// @Param	info	body models.UserUpdate true "User information"
// @Success	200	{object} models.User
// @Failure	404	{object} error.APIError
// @Failure	500	{object} error.APIError
// @Router	/api/v1/users/{id} [put]
// @Security    OAuth2Password
func (api *V1UserAPI) UpdateUser(ctx *gin.Context) {
	id := ctx.Param("id")
	var userIn models.UserUpdate
	if err := ctx.ShouldBindJSON(&userIn); err != nil {
		v1.SendError(ctx, http.StatusNotFound, err)
		return
	}

	user, err := api.UserService.Update(id, userIn)
	if err != nil {
		v1.SendError(ctx, http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, user)
}

// DeleteUser godoc
// @Summary	Delete a user
// @Tags	users
// @Produce	json
// @Param	id	path integer true "User ID"
// @Success	200	{object} models.User
// @Failure	404	{object} error.APIError
// @Router	/api/v1/users/{id} [delete]
// @Security    OAuth2Password
func (api *V1UserAPI) DeleteUser(ctx *gin.Context) {
	id := ctx.Param("id")

	err := api.UserService.Delete(id)
	if err != nil {
		v1.SendError(ctx, http.StatusNotFound, err)
		return
	}

	ctx.JSON(http.StatusOK, nil)
}
