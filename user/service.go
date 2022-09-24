package user

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)
type Service interface {
	RegisterUser(input RegisterUserInput) (User, error)
	Login(input LoginUserInput) (User, error)
	IsEmailAvailable(input CheckEmailInput) (bool, error)
	SaveAvatar(id int, file string) (User, error)
}

type service struct {
	repository Repository
}

func NewService(r Repository) *service {
	return &service{r}
}

func (s *service) RegisterUser(input RegisterUserInput) (User, error) {
	user := User{}
	user.Name = input.Name
	user.Email = input.Email
	user.Avatar = input.Avatar
	user.Occupation = input.Occupation
	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.MinCost)

	if err != nil {
		return user, err
	}
	user.Password = string(hash)
	newUser, err := s.repository.Save(user)

	if err != nil {
		return newUser, err
	}

	return newUser, nil

}

func (s *service) Login(input LoginUserInput) (User, error){
	email := input.Email
	password := input.Password

	user, err := s.repository.FindByEmail(email)
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

func (s *service) IsEmailAvailable(input CheckEmailInput) (bool, error){
	email := input.Email

	user, err := s.repository.FindByEmail(email)
	if err != nil {
		return false, err 
	}
	if user.ID == 0 {
		return true, nil 
	}

	return false, nil
}

func (s *service) SaveAvatar(ID int, fileName string) (User, error){

	user, err := s.repository.FindByID(ID)
	if err != nil {
		return user, err 
	}
	user.Avatar = fileName

	updatedUser, err := s.repository.Update(user)
	if err != nil {
		return updatedUser, err 
	}
	return updatedUser, nil

	

}