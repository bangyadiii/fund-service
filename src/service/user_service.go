package service

import (
	"backend-crowdfunding/src/model"
	"backend-crowdfunding/src/repository"
	"backend-crowdfunding/src/request"
	"context"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	RegisterUser(ctx context.Context, input request.RegisterUserInput) (model.User, error)
	Login(ctx context.Context, input request.LoginUserInput) (model.User, error)
	IsEmailAvailable(ctx context.Context, input request.CheckEmailInput) (bool, error)
	SaveAvatar(ctx context.Context, ID uint, file string) (model.User, error)
	FindByID(ctx context.Context, ID uint) (model.User, error)
}

type userService struct {
	repository repository.UserRepository
	timeout    time.Duration
}

func NewUserService(r repository.UserRepository) UserService {
	return &userService{
		repository: r,
		timeout:    2 * time.Second,
	}
}

func (s *userService) RegisterUser(ctx context.Context, input request.RegisterUserInput) (model.User, error) {
	context, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	user := model.User{}
	user.Name = input.Name
	user.Email = input.Email
	user.Avatar = input.Avatar
	user.Occupation = input.Occupation
	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.MinCost)

	if err != nil {
		return user, err
	}
	user.Password = string(hash)

	newUser, err := s.repository.SaveUser(context, user)

	if err != nil {
		return newUser, err
	}

	return newUser, nil

}

func (s *userService) Login(ctx context.Context, input request.LoginUserInput) (model.User, error) {
	email := input.Email
	password := input.Password

	user, err := s.repository.FindByEmailUser(ctx, email)
	if err != nil {
		return user, err
	}
	if user.ID == 0 {
		return user, errors.New("there is no user with this email")
	}
	err1 := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	if err1 != nil {
		return user, err1
	}

	return user, nil
}

func (s *userService) IsEmailAvailable(ctx context.Context, input request.CheckEmailInput) (bool, error) {
	email := input.Email

	user, err := s.repository.FindByEmailUser(ctx, email)
	if err != nil {
		return false, err
	}
	if user.ID == 0 {
		return true, nil
	}

	return false, nil
}

func (s *userService) SaveAvatar(ctx context.Context, ID uint, fileName string) (model.User, error) {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	user, err := s.repository.FindByIDUser(c, ID)
	if err != nil {
		return user, err
	}
	user.Avatar = fileName

	updatedUser, err := s.repository.UpdateUser(c, user)
	if err != nil {
		return updatedUser, err
	}
	return updatedUser, nil
}
func (s *userService) FindByID(ctx context.Context, ID uint) (model.User, error) {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()
	
	user, err := s.repository.FindByIDUser(c, ID)
	if err != nil {
		return user, err
	}
	return user, nil
}
