package password

import (
	"backend-crowdfunding/sdk/errors"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", errors.NewErrorf(500, nil, "hash error", err)
	}
	return string(hashed), nil
}

func ComparePassword(hashedPassword string, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func IsPasswordConfirmed(password string, passwordConf string) error {
	if password != passwordConf {
		return errors.NewErrorf(500, nil, "password_confirmation doesn't match")
	}

	return nil
}
