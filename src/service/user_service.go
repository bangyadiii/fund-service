package service

import (
	"backend-crowdfunding/src/model"
	"backend-crowdfunding/src/repository"
	"backend-crowdfunding/src/request"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	RegisterUser(input request.RegisterUserInput) (model.User, error)
	Login(input request.LoginUserInput) (model.User, error)
	IsEmailAvailable(input request.CheckEmailInput) (bool, error)
	SaveAvatar(ID uint, file string) (model.User, error)
	FindByID(ID uint) (model.User, error)
}

type userService struct {
	repository repository.UserRepository
}

func NewService(r repository.UserRepository) UserService {
	return &userService{r}
}

func (s *userService) RegisterUser(input request.RegisterUserInput) (model.User, error) {
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

	newUser, err := s.repository.SaveUser(user)

	if err != nil {
		return newUser, err
	}

	return newUser, nil

}

func (s *userService) Login(input request.LoginUserInput) (model.User, error) {
	email := input.Email
	password := input.Password

	user, err := s.repository.FindByEmailUser(email)
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

func (s *userService) IsEmailAvailable(input request.CheckEmailInput) (bool, error) {
	email := input.Email

	user, err := s.repository.FindByEmailUser(email)
	if err != nil {
		return false, err
	}
	if user.ID == 0 {
		return true, nil
	}

	return false, nil
}

func (s *userService) SaveAvatar(ID uint, fileName string) (model.User, error) {

	user, err := s.repository.FindByIDUser(ID)
	if err != nil {
		return user, err
	}
	user.Avatar = fileName

	updatedUser, err := s.repository.UpdateUser(user)
	if err != nil {
		return updatedUser, err
	}
	return updatedUser, nil
}
func (s *userService) FindByID(ID uint) (model.User, error) {

	user, err := s.repository.FindByIDUser(ID)
	if err != nil {
		return user, err
	}
	return user, nil
}
