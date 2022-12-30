package repository

import (
	"backend-crowdfunding/src/model"

	"gorm.io/gorm"
)

type UserRepository interface {
	FindByEmailUser(email string) (model.User, error)
	FindByIDUser(id uint) (model.User, error)
	SaveUser(user model.User) (model.User, error)
	UpdateUser(user model.User) (model.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *userRepository {
	return &userRepository{db}
}

func (r *userRepository) SaveUser(user model.User) (model.User, error) {
	data := r.db.Create(&user)

	if data.Error != nil {
		return user, data.Error
	}
	return user, nil
}

func (r *userRepository) FindByEmailUser(email string) (model.User, error) {
	var user model.User
	err := r.db.Where("email = ?", email).Find(&user).Error

	if err != nil {
		return user, err
	}
	return user, nil
}

func (r *userRepository) FindByIDUser(id uint) (model.User, error) {
	var user model.User
	err := r.db.Where("id = ?", id).Find(&user).Error

	if err != nil {
		return user, err
	}
	return user, nil
}

func (r *userRepository) UpdateUser(user model.User) (model.User, error) {
	err := r.db.Save(&user).Error

	if err != nil {
		return user, err
	}
	return user, nil
}
