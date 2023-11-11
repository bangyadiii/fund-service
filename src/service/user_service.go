package service

import (
	"backend-crowdfunding/sdk/errors"
	"backend-crowdfunding/sdk/password"
	"backend-crowdfunding/src/dto/request"
	"backend-crowdfunding/src/dto/response"
	"backend-crowdfunding/src/model"
	"backend-crowdfunding/src/repository"
	"context"
	"time"
)

type UserService interface {
	RegisterUser(ctx context.Context, input request.RegisterUserInput) (model.User, string, error)
	Login(ctx context.Context, input request.LoginUserInput) (response.UserLoginResponse, error)
	IsEmailAvailable(ctx context.Context, input request.CheckEmailInput) (bool, error)
	SaveAvatar(ctx context.Context, ID string, fileName string) (model.User, error)
	FindByID(ctx context.Context, ID string) (response.UserResponse, error)
	LoginWithGoogle(ctx context.Context, input request.LoginWithGoogleInput) (response.UserLoginResponse, error)
}

type userServiceImpl struct {
	repository  repository.UserRepository
	authService AuthService
	timeout     time.Duration
}

// NewUserService is a function that takes a `UserRepository` and an `AuthService` and returns a `UserService`
func NewUserService(r repository.UserRepository, authService AuthService) UserService {
	return &userServiceImpl{
		repository:  r,
		authService: authService,
		timeout:     2 * time.Second,
	}
}

func (s *userServiceImpl) RegisterUser(ctx context.Context, input request.RegisterUserInput) (model.User, string, error) {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	user := model.User{}
	user.Name = input.Name
	user.Email = input.Email
	user.Avatar = input.Avatar
	user.Occupation = input.Occupation
	hash, err := password.HashPassword(input.Password)
	var token string

	if err != nil {
		return user, token, err
	}

	user.Password = hash

	newUser, err := s.repository.SaveUser(c, user)

	if err != nil {
		return newUser, token, err
	}

	token, err = s.authService.GenerateToken(newUser.ID, newUser.Email)

	return newUser, token, nil

}

// Login This function is used to log in with email and password.
func (s *userServiceImpl) Login(ctx context.Context, input request.LoginUserInput) (response.UserLoginResponse, error) {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	var loginResponse response.UserLoginResponse

	email := input.Email
	pwd := input.Password

	user, err := s.repository.FindByEmailUser(c, email)
	if err != nil {
		return loginResponse, err
	}

	if user.ID == "" {
		return loginResponse, errors.NewErrorf(400, nil, "there is no user with this email")
	}

	if user.IsGoogleAccount {
		return loginResponse, errors.NewErrorf(400, nil, "this is google account")
	}

	err = password.ComparePassword(user.Password, pwd)

	if err != nil {
		return loginResponse, err
	}
	newToken, err := s.authService.GenerateToken(user.ID, user.Email)

	if err != nil {
		return loginResponse, err
	}
	loginResponse = response.FormatUserLogin(&user, newToken)
	return loginResponse, nil
}

// LoginWithGoogle A function that is used to log in with Google account.
func (s *userServiceImpl) LoginWithGoogle(ctx context.Context, input request.LoginWithGoogleInput) (response.UserLoginResponse, error) {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	var userFormat response.UserLoginResponse

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
	userFormat = response.FormatUserLogin(&user, token)
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

// IsEmailAvailable Checking if the email is available or not.
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

// FindByID A function that is used to find user by ID.
func (s *userServiceImpl) FindByID(ctx context.Context, ID string) (response.UserResponse, error) {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()
	var userRes response.UserResponse

	user, err := s.repository.FindByIDUser(c, ID)
	if err != nil {
		return userRes, err
	}
	userRes = response.FormatUserResponse(&user)
	return userRes, nil
}
