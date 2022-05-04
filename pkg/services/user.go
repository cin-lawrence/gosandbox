package services

import (
        log "github.com/sirupsen/logrus"

        "github.com/cin-lawrence/gosandbox/pkg/db"
        "github.com/cin-lawrence/gosandbox/pkg/models"
)

type UserService struct {
        *EntityService
}

func NewUserService() *UserService {
        srv := &UserService{
                &EntityService{db: db.DB},
        }
        db.DB.AutoMigrate(&models.User{})

        return srv
}

func (srv *UserService) List() ([]models.User, error) {
        var models []models.User
        session := srv.newSession()
        result := session.Find(&models)
        if result.Error != nil {
                log.Error("List models failed")
                return nil, result.Error
        }
        return models, nil
}

func (srv *UserService) Create(model models.User) (models.User, error) {
        session := srv.newSession()
        result := session.Create(&model)
        if result.Error != nil {
                log.Errorf("Create user failed: %v", result.Error)
                return models.User{}, result.Error
        }

        session.First(&model, model.ID)
        return model, nil
}

func (srv *UserService) Get(id string) (models.User, error) {
        var model models.User
        session := srv.newSession()
        result := session.First(&model, id)
        if result.Error != nil {
                log.Errorf("Get user %s failed: %v", id, result.Error)
                return models.User{}, result.Error
        }
        return model, nil
}

func (srv *UserService) Update(id string, model models.User) (models.User, error) {
        modelInDB, err := srv.Get(id)
        if err != nil {
                log.Errorf("User %s not found: %v", id, err)
                return models.User{}, err
        }
        session := srv.newSession()
        result := session.Model(&modelInDB).Updates(model)
        if result.Error != nil {
                log.Errorf("Update user % failed: %v", id, result.Error)
                return models.User{}, result.Error
        }
        return modelInDB, nil
}

func (srv *UserService) Delete(id string) error {
        modelInDB, err := srv.Get(id)
        if err != nil {
                log.Errorf("User %s not found: %v",  id, err)
                return err
        }
        session := srv.newSession()
        result := session.Delete(&modelInDB)
        if result.Error != nil {
                log.Errorf("Delete user %s failed: %v", id, result.Error)
                return result.Error
        }
        return nil
}
