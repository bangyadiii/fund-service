package service

import (
	"errors"
	"os"

	"github.com/golang-jwt/jwt/v4"
)

type AuthService interface {
	GenerateToken(ID uint, email string) (string, error)
	ValidateToken(token string) (*jwt.Token, error)
}

type jwtServiceImpl struct {
	//
}

var secretKey string = os.Getenv("SECRET_KEY")
var SECRET_KEY = []byte(secretKey)

func (s *jwtServiceImpl) GenerateToken(ID uint, email string) (string, error) {
	//JWT
	//claim = payload
	claim := jwt.MapClaims{}
	claim["_id"] = ID
	claim["email"] = email

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	signedToken, err := token.SignedString(SECRET_KEY)
	if err != nil {
		return signedToken, err
	}

	return signedToken, nil

}
func (s *jwtServiceImpl) ValidateToken(token string) (*jwt.Token, error) {
	//
	decodedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		_, ok := t.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.New("invalid token")
		}
		return []byte(SECRET_KEY), nil
	})

	if err != nil {
		return decodedToken, err
	}

	return decodedToken, nil
}

func NewAuthService() AuthService {
	return &jwtServiceImpl{}
}
