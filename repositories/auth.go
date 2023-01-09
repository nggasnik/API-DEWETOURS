package repositories

import (
	"erlangga/models"
	"errors"
)

type AuthRepository interface {
	Register(newUser models.User) (models.User, error)
	Login(email string) (models.User, error)
	GetUserLogin (id int) (models.User, error)
}

func (r *repository) Register(newUser models.User) (models.User, error) {
	var user models.User
	errCekUser := r.db.First(&user, "email=?", newUser.Email).Error
	if errCekUser == nil {
		return user, errors.New("user already registered")
	}

	errAddUser := r.db.Create(&newUser).Error
	return newUser, errAddUser
}

func (r *repository) Login(email string) (models.User, error) {
	var user models.User
	err := r.db.First(&user, "email=?", email).Error
	return user, err
}

func (r *repository) GetUserLogin (id int) (models.User, error) {
	var user models.User
	err := r.db.First(&user, "id=?", id).Error
	return user, err
}
