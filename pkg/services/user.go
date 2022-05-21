package services

import (
	"errors"

	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"

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

func (srv *UserService) List() (users []models.User, err error) {
	session := srv.newSession()
	result := session.Find(&users)
	if result.Error != nil {
		log.Error("List models failed")
		return nil, result.Error
	}
	return users, nil
}

func (srv *UserService) Create(userInput models.UserInput) (user models.User, err error) {
	session := srv.newSession()
	bytePassword := []byte(userInput.Password)
	hashedPassword, err := bcrypt.GenerateFromPassword(bytePassword, bcrypt.DefaultCost)
	if err != nil {
		return user, errors.New("Can't generate hashed password")
	}
	user = models.User{
		Name:           userInput.Name,
		Email:          userInput.Email,
		HashedPassword: string(hashedPassword),
		IsActive:       true,
	}

	result := session.Create(&user)
	if result.Error != nil {
		log.Errorf("Create user failed: %v", result.Error)
		return models.User{}, result.Error
	}

	return user, nil
}

func (srv *UserService) Get(id string) (user models.User, err error) {
	session := srv.newSession()
	result := session.First(&user, id)
	if result.Error != nil {
		log.Errorf("Get user %s failed: %v", id, result.Error)
		return models.User{}, result.Error
	}
	return user, nil
}

func (srv *UserService) Update(id string, userIn models.UserUpdate) (user models.User, err error) {
	user, err = srv.Get(id)
	if err != nil {
		log.Errorf("User %s not found: %v", id, err)
		return user, err
	}
	session := srv.newSession()
	result := session.Model(&user).Updates(userIn)
	if result.Error != nil {
		log.Errorf("Update user %s failed: %v", id, result.Error)
		return user, result.Error
	}
	return user, nil
}

func (srv *UserService) Delete(id string) error {
	modelInDB, err := srv.Get(id)
	if err != nil {
		log.Errorf("User %s not found: %v", id, err)
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

func (srv *UserService) GetByEmail(email string) (user models.User, err error) {
	session := srv.newSession()
	res := session.Model(models.User{Email: email}).First(&user)
	if res.Error != nil {
		return user, res.Error
	}
	return user, nil
}
