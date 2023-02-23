package repository

import (
	"backend-crowdfunding/database"
	"backend-crowdfunding/src/model"
	"backend-crowdfunding/src/util/id"
	"context"
)

type UserRepository interface {
	FindByEmailUser(ctx context.Context, email string) (model.User, error)
	FindByIDUser(ctx context.Context, id string) (model.User, error)
	SaveUser(ctx context.Context, user model.User) (model.User, error)
	UpdateUser(ctx context.Context, user model.User) (model.User, error)
}

type userRepoImpl struct {
	db          *database.DB
	idGenerator id.IDGenerator
}

func NewUserRepository(db *database.DB, idGenerator id.IDGenerator) UserRepository {
	return &userRepoImpl{
		db,
		idGenerator,
	}
}

func (r *userRepoImpl) SaveUser(ctx context.Context, user model.User) (model.User, error) {
	id := r.idGenerator.Generate()
	user.ID = id
	data := r.db.WithContext(ctx).Create(&user)

	if data.Error != nil {
		return user, data.Error
	}
	return user, nil
}

func (r *userRepoImpl) FindByEmailUser(ctx context.Context, email string) (model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).Where("email = ?", email).Find(&user).Error

	if err != nil {
		return user, err
	}
	return user, nil
}

func (r *userRepoImpl) FindByIDUser(ctx context.Context, id string) (model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).Where("id = ?", id).Find(&user).Error

	if err != nil {
		return user, err
	}
	return user, nil
}

func (r *userRepoImpl) UpdateUser(ctx context.Context, user model.User) (model.User, error) {
	err := r.db.WithContext(ctx).Save(&user).Error

	if err != nil {
		return user, err
	}
	return user, nil
}
