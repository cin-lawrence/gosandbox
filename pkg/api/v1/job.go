package v1

import (
	"math/rand"
	"net/http"
	"strconv"
	"time"

	gc "github.com/gocelery/gocelery"
	log "github.com/sirupsen/logrus"

	"github.com/cin-lawrence/gosandbox/pkg/api/error"
	"github.com/cin-lawrence/gosandbox/pkg/config"
	"github.com/cin-lawrence/gosandbox/pkg/models"
	"github.com/cin-lawrence/gosandbox/pkg/services"
	"github.com/cin-lawrence/gosandbox/pkg/wkrce"
	"github.com/gin-gonic/gin"
)

type V1JobAPI struct {
	JobService   services.JobService
	UserService  services.UserService
	WorkerClient *gc.CeleryClient
}

var v1JobAPI *V1JobAPI

func NewV1JobGroup(rg *gin.RouterGroup) *gin.RouterGroup {
	v1JobAPI = &V1JobAPI{
		JobService:   *services.NewJobService(),
		UserService:  *services.NewUserService(),
		WorkerClient: worker.NewWorkerClient(),
	}

	v1JobGroup := rg.Group("/jobs", Authorize())
	v1JobGroup.GET("/", v1JobAPI.ListJobs)
	v1JobGroup.POST("/", v1JobAPI.CreateJob)
	v1JobGroup.GET("/:id", v1JobAPI.GetJob)
	v1JobGroup.PUT("/:id", v1JobAPI.UpdateJob)
	v1JobGroup.DELETE("/:id", v1JobAPI.DeleteJob)
	return v1JobGroup
}

func (api *V1JobAPI) ListJobs(ctx *gin.Context) {
	jobs, err := api.JobService.List()
	if err != nil {
		v1.SendError(ctx, http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, models.JobList{Items: jobs})
}

func (api *V1JobAPI) CreateJob(ctx *gin.Context) {
	var job models.Job
	if err := ctx.ShouldBindJSON(&job); err != nil {
		v1.SendError(ctx, http.StatusBadRequest, err)
		return
	}

	_, err := api.UserService.Get(
		strconv.FormatUint(uint64(job.UserID), 10),
	)
	if err != nil {
		v1.SendError(ctx, http.StatusNotFound, err)
		return
	}

	job.Status = models.JobStatusPending
	job, err = api.JobService.Create(job)
	if err != nil {
		v1.SendError(ctx, http.StatusInternalServerError, err)
	}

	rand.Seed(time.Now().UnixNano())
	randomNumber := rand.Intn(1000)
	log.Infof("Sending job %d with parameter %d", job.ID, randomNumber)
	_, err = api.WorkerClient.Delay(
		config.Config.CeleryTaskName,
		randomNumber,
	)
	if err != nil {
		v1.SendError(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusCreated, job)
}

func (api *V1JobAPI) GetJob(ctx *gin.Context) {
	id := ctx.Param("id")
	job, err := api.JobService.Get(id)
	if err != nil {
		v1.SendError(ctx, http.StatusNotFound, err)
		return
	}

	ctx.JSON(http.StatusOK, job)
}

func (api *V1JobAPI) UpdateJob(ctx *gin.Context) {
	id := ctx.Param("id")
	var job models.Job
	if err := ctx.ShouldBindJSON(&job); err != nil {
		v1.SendError(ctx, http.StatusNotFound, err)
		return
	}

	job, err := api.JobService.Update(id, job)
	if err != nil {
		v1.SendError(ctx, http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, job)
}

func (api *V1JobAPI) DeleteJob(ctx *gin.Context) {
	id := ctx.Param("id")

	err := api.JobService.Delete(id)
	if err != nil {
		v1.SendError(ctx, http.StatusNotFound, err)
		return
	}

	ctx.JSON(http.StatusOK, nil)
}
