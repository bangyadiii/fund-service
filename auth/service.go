package auth

import (
	"github.com/golang-jwt/jwt/v4"
)

type Service interface {
	GenerateToken(ID int, email string) (string, error)
}

type jwtService struct {
	//
}
var SECRET_KEY = []byte("SECRET_KEY_RAHASIA")

func (s *jwtService) GenerateToken(ID int, email string) (string, error) {
	//JWT
	//claim = payload
	claim := jwt.MapClaims{}
	claim["_id"] = ID
	claim["email"] = email

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	signedToken, err := token.SignedString(SECRET_KEY)
	if err != nil{
		return signedToken, err
	}

	return signedToken, nil
	
}

func NewService() *jwtService{
	return &jwtService{}
}