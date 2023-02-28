package service

import (
	"backend-crowdfunding/sdk/password"
	"backend-crowdfunding/src/formatter"
	"backend-crowdfunding/src/model"
	"backend-crowdfunding/src/repository"
	"backend-crowdfunding/src/request"
	"context"
	"errors"
	"time"
)

type UserService interface {
	RegisterUser(ctx context.Context, input request.RegisterUserInput) (model.User, error)
	Login(ctx context.Context, input request.LoginUserInput) (model.User, error)
	IsEmailAvailable(ctx context.Context, input request.CheckEmailInput) (bool, error)
	SaveAvatar(ctx context.Context, ID string, file string) (model.User, error)
	FindByID(ctx context.Context, ID string) (model.User, error)
	LoginWithGoogle(ctx context.Context, input request.LoginWithGoogleInput) (formatter.UserLoginFormatter, error)
}

type userServiceImpl struct {
	repository  repository.UserRepository
	authService AuthService
	timeout     time.Duration
}

func NewUserService(r repository.UserRepository, authService AuthService) UserService {
	return &userServiceImpl{
		repository:  r,
		authService: authService,
		timeout:     2 * time.Second,
	}
}

func (s *userServiceImpl) RegisterUser(ctx context.Context, input request.RegisterUserInput) (model.User, error) {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	user := model.User{}
	user.Name = input.Name
	user.Email = input.Email
	user.Avatar = input.Avatar
	user.Occupation = input.Occupation
	hash, err := password.HashPassword(input.Password)

	if err != nil {
		return user, err
	}
	user.Password = hash

	newUser, err := s.repository.SaveUser(c, user)

	if err != nil {
		return newUser, err
	}

	return newUser, nil

}

func (s *userServiceImpl) Login(ctx context.Context, input request.LoginUserInput) (model.User, error) {
	email := input.Email
	pwd := input.Password

	user, err := s.repository.FindByEmailUser(ctx, email)
	if err != nil {
		return user, err
	}

	if user.ID == "" {
		return user, errors.New("there is no user with this email")
	}

	if user.IsGoogleAccount {
		return model.User{}, errors.New("this is google account")
	}

	err = password.ComparePassword(user.Password, pwd)

	if err != nil {
		return user, err
	}

	return user, nil
}

func (s *userServiceImpl) LoginWithGoogle(ctx context.Context, input request.LoginWithGoogleInput) (formatter.UserLoginFormatter, error) {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	var userFormat formatter.UserLoginFormatter

	firebaseToken, err := s.repository.VerifyFirebaseToken(c, input.FirebaseToken)
	if err != nil {
		return userFormat, err
	}
	email := firebaseToken.Claims["email"].(string)
	name := firebaseToken.Claims["name"].(string)

	user, err := s.repository.Get(c, request.UserParam{
		Email:           email,
		IsGoogleAccount: true,
	})
	if err != nil && err.Error() == "NOT_FOUND" {
		user, err = s.registerFromGoogleAccount(c, request.RegisterUserInput{Email: email, Name: name})
		if err != nil {
			return userFormat, err
		}
	} else if err != nil {
		return userFormat, err
	}

	token, err := s.authService.GenerateToken(user.ID, user.Email)
	if err != nil {
		return userFormat, err
	}
	userFormat = formatter.FormatUserLogin(user, token)
	return userFormat, nil

}

func (s userServiceImpl) registerFromGoogleAccount(ctx context.Context, param request.RegisterUserInput) (model.User, error) {
	var user model.User
	user = model.User{
		Email:           param.Email,
		Name:            param.Name,
		IsGoogleAccount: true,
	}
	user, err := s.repository.SaveUser(ctx, user)
	if err != nil {
		return user, err
	}
	return user, nil
}

func (s *userServiceImpl) IsEmailAvailable(ctx context.Context, input request.CheckEmailInput) (bool, error) {
	email := input.Email

	user, err := s.repository.FindByEmailUser(ctx, email)
	if err != nil {
		return false, err
	}
	if user.ID == "" {
		return true, nil
	}

	return false, nil
}

func (s *userServiceImpl) SaveAvatar(ctx context.Context, ID string, fileName string) (model.User, error) {
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

func (s *userServiceImpl) FindByID(ctx context.Context, ID string) (model.User, error) {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	user, err := s.repository.FindByIDUser(c, ID)
	if err != nil {
		return user, err
	}
	return user, nil
}
