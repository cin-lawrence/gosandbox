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
	v1JobGroup.DELETE("/:id", v1JobAPI.DeleteJob)
	return v1JobGroup
}

// ListJobs godoc
// @Summary	List all jobs
// @Tags	jobs
// @Produce	json
// @Success	200	{object} models.JobList
// @Failure	500	{object} error.APIError
// @Router	/api/v1/jobs/ [get]
// @Security    OAuth2Password
func (api *V1JobAPI) ListJobs(ctx *gin.Context) {
	jobs, err := api.JobService.List()
	if err != nil {
		v1.SendError(ctx, http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, models.JobList{Items: jobs})
}

// CreateUsers godoc
// @Summary	Create a new job
// @Description Create a simple Celery task that does a random fibonacci calculation.
// @Tags	jobs
// @Accept	json
// @Produce	json
// @Param	meta	body models.JobCreate true "Job meta"
// @Success	201	{object} models.Job
// @Failure	404	{object} error.APIError
// @Failure	422	{object} error.APIError
// @Failure	500	{object} error.APIError
// @Router	/api/v1/jobs/ [post]
// @Security    OAuth2Password
func (api *V1JobAPI) CreateJob(ctx *gin.Context) {
	var jobIn models.JobCreate
	if err := ctx.ShouldBindJSON(&jobIn); err != nil {
		v1.SendError(ctx, http.StatusUnprocessableEntity, err)
		return
	}

	_, err := api.UserService.Get(
		strconv.FormatUint(uint64(jobIn.UserID), 10),
	)
	if err != nil {
		v1.SendError(ctx, http.StatusNotFound, err)
		return
	}

	job, err := api.JobService.Create(jobIn)
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

// GetJob godoc
// @Summary	Retrieve a job
// @Tags	jobs
// @Produce	json
// @Param	id	path integer true "Job ID"
// @Success	200	{object} models.Job
// @Failure	404	{object} error.APIError
// @Router	/api/v1/jobs/{id} [get]
// @Security    OAuth2Password
func (api *V1JobAPI) GetJob(ctx *gin.Context) {
	id := ctx.Param("id")
	job, err := api.JobService.Get(id)
	if err != nil {
		v1.SendError(ctx, http.StatusNotFound, err)
		return
	}

	ctx.JSON(http.StatusOK, job)
}

// DeleteJob godoc
// @Summary	Delete a job
// @Tags	jobs
// @Produce	json
// @Param	id	path integer true "jobs ID"
// @Success	200	{object} models.Job
// @Failure	404	{object} error.APIError
// @Router	/api/v1/jobs/{id} [delete]
// @Security	OAuth2Password
func (api *V1JobAPI) DeleteJob(ctx *gin.Context) {
	id := ctx.Param("id")

	err := api.JobService.Delete(id)
	if err != nil {
		v1.SendError(ctx, http.StatusNotFound, err)
		return
	}

	ctx.JSON(http.StatusOK, nil)
}
