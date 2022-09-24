package user

import (
	"gorm.io/gorm"
)

type Repository interface {
	FindByEmail(email string) (User, error)
	Save(user User) (User, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository{
	return &repository{db}
}

func (r *repository) Save(user User) (User, error){
	data := r.db.Create(&user)

	if data.Error != nil {
		return user, data.Error
	}
	return user, nil
}

func (r *repository) FindByEmail(email string) (User, error){
	var user User
	err := r.db.Where("email = ?", email).Find(&user).Error

	if err != nil {
		return user, err
	}
	return user, nil

}