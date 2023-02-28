package repository

import (
	"backend-crowdfunding/database"
	"backend-crowdfunding/insfrastructure/firebase"
	"backend-crowdfunding/sdk/id"
	"backend-crowdfunding/src/model"
	"backend-crowdfunding/src/request"
	"context"
	"errors"
	firebase_auth "firebase.google.com/go/auth"
)

type UserRepository interface {
	FindByEmailUser(ctx context.Context, email string) (model.User, error)
	FindByIDUser(ctx context.Context, id string) (model.User, error)
	Get(ctx context.Context, param request.UserParam) (model.User, error)
	SaveUser(ctx context.Context, user model.User) (model.User, error)
	UpdateUser(ctx context.Context, user model.User) (model.User, error)
	VerifyFirebaseToken(ctx context.Context, firebaseToken string) (*firebase_auth.Token, error)
}

type userRepoImpl struct {
	db          *database.DB
	firebase    *firebase.Firebase
	idGenerator id.IDGenerator
}

func NewUserRepository(db *database.DB, firebase *firebase.Firebase, idGenerator id.IDGenerator) UserRepository {
	return &userRepoImpl{
		db,
		firebase,
		idGenerator,
	}
}

func (r *userRepoImpl) SaveUser(ctx context.Context, user model.User) (model.User, error) {
	user.ID = r.idGenerator.Generate()
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

func (r userRepoImpl) Get(ctx context.Context, param request.UserParam) (model.User, error) {
	var user model.User
	res := r.db.WithContext(ctx).Where(param).First(&user)
	if res.RowsAffected == 0 {
		return user, errors.New("NOT_FOUND")
	} else if res.Error != nil {
		return user, res.Error
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

func (r *userRepoImpl) VerifyFirebaseToken(ctx context.Context, firebaseToken string) (*firebase_auth.Token, error) {
	return r.firebase.Auth.VerifyIDToken(ctx, firebaseToken)
}
