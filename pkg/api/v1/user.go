package v1

import (
	"net/http"

	"github.com/cin-lawrence/gosandbox/pkg/api/error"
	"github.com/cin-lawrence/gosandbox/pkg/models"
	"github.com/cin-lawrence/gosandbox/pkg/services"
	"github.com/gin-gonic/gin"
)

type V1UserAPI struct {
	UserService services.UserService
}

type UserList struct {
	Items []models.User `json:"items"`
}

type UserInput struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required"`
}

var v1UserAPI *V1UserAPI

func NewV1UserGroup(rg *gin.RouterGroup) *gin.RouterGroup {
	v1UserAPI = &V1UserAPI{
		UserService: *services.NewUserService(),
	}

	v1UserGroup := rg.Group("/users")
	v1UserGroup.GET("/", v1UserAPI.ListUsers)
	v1UserGroup.POST("/", v1UserAPI.CreateUser)
	v1UserGroup.GET("/:id", v1UserAPI.GetUser)
	v1UserGroup.PUT("/:id", v1UserAPI.UpdateUser)
	v1UserGroup.DELETE("/:id", v1UserAPI.DeleteUser)
	return v1UserGroup
}

func (api *V1UserAPI) ListUsers(ctx *gin.Context) {
	users, err := api.UserService.List()
	if err != nil {
		v1.SendError(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, UserList{Items: users})
}

func (api *V1UserAPI) CreateUser(ctx *gin.Context) {
	var user models.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		v1.SendError(ctx, http.StatusBadRequest, err)
		return
	}

	user, err := api.UserService.Create(user)
	if err != nil {
		v1.SendError(ctx, http.StatusBadRequest, err)
		return
	}
	ctx.JSON(http.StatusCreated, user)
}

func (api *V1UserAPI) GetUser(ctx *gin.Context) {
	id := ctx.Param("id")
	user, err := api.UserService.Get(id)
	if err != nil {
		v1.SendError(ctx, http.StatusNotFound, err)
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func (api *V1UserAPI) UpdateUser(ctx *gin.Context) {
	id := ctx.Param("id")
	var user models.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		v1.SendError(ctx, http.StatusNotFound, err)
		return
	}

	user, err := api.UserService.Update(id, user)
	if err != nil {
		v1.SendError(ctx, http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, user)
}

func (api *V1UserAPI) DeleteUser(ctx *gin.Context) {
	id := ctx.Param("id")

	err := api.UserService.Delete(id)
	if err != nil {
		v1.SendError(ctx, http.StatusNotFound, err)
		return
	}

	ctx.JSON(http.StatusOK, nil)
}
