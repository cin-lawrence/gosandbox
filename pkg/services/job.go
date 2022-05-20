package services

import (
	"github.com/cin-lawrence/gosandbox/pkg/db"
	"github.com/cin-lawrence/gosandbox/pkg/models"
	log "github.com/sirupsen/logrus"
)

type JobService struct {
	*EntityService
}

func NewJobService() *JobService {
	srv := &JobService{
		&EntityService{db: db.DB.Preload("User")},
	}
	db.DB.Exec(models.CREATE_ENUM_JOB_STATUS)
	db.DB.AutoMigrate(&models.Job{})

	return srv
}

func (srv *JobService) List() ([]models.Job, error) {
	var models []models.Job
	session := srv.newSession()
	result := session.Find(&models)
	if result.Error != nil {
		log.Error("List models failed")
		return nil, result.Error
	}
	return models, nil
}

func (srv *JobService) Create(jobIn models.JobCreate) (job models.Job, err error) {
	session := srv.newSession()
	job = models.Job{
		Status: models.JobStatusPending,
		UserID: jobIn.UserID,
	}
	result := session.Create(&job)
	if result.Error != nil {
		log.Errorf("Create job failed: %v", result.Error)
		return models.Job{}, result.Error
	}

	session.First(&job, job.ID)
	return job, nil
}

func (srv *JobService) Get(id string) (models.Job, error) {
	var model models.Job
	session := srv.newSession()
	result := session.First(&model, id)
	if result.Error != nil {
		log.Errorf("Get job %s failed: %v", id, result.Error)
		return models.Job{}, result.Error
	}
	return model, nil
}

func (srv *JobService) Update(id string, model models.Job) (models.Job, error) {
	modelInDB, err := srv.Get(id)
	if err != nil {
		log.Errorf("Job %s not found: %v", id, err)
		return models.Job{}, err
	}
	session := srv.newSession()
	result := session.Model(&modelInDB).Updates(model)
	if result.Error != nil {
		log.Errorf("Update job % failed: %v", id, result.Error)
		return models.Job{}, result.Error
	}
	return modelInDB, nil
}

func (srv *JobService) Delete(id string) error {
	modelInDB, err := srv.Get(id)
	if err != nil {
		log.Errorf("Job %s not found: %v", id, err)
		return err
	}
	session := srv.newSession()
	result := session.Delete(&modelInDB)
	if result.Error != nil {
		log.Errorf("Delete job %s failed: %v", id, result.Error)
		return result.Error
	}
	return nil
}
