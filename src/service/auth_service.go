package service

import (
	"backend-crowdfunding/config"
	"backend-crowdfunding/sdk/errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type AuthService interface {
	GenerateToken(ID string, email string) (string, error)
	ValidateToken(token string) (*jwt.Token, error)
	ConvertTokenToCustomClaims(accessToken *jwt.Token) (*CustomClaims, bool)
}

type JwtServiceImpl struct {
	env config.Config
}

type CustomClaims struct {
	jwt.RegisteredClaims
	Email string `json:"email"`
}

func (s *JwtServiceImpl) GenerateToken(ID string, email string) (string, error) {
	//JWT
	SecretKey := []byte(s.env.GetWithDefault("ACCESS_SECRET_KEY", "rahasiabanget"))
	//claim = payload
	registerClaims := jwt.RegisteredClaims{
		ID:        ID,
		Issuer:    "fund-platform",
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 2)),
	}

	claim := CustomClaims{
		registerClaims,
		email,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	signedToken, err := token.SignedString(SecretKey)
	if err != nil {
		return signedToken, errors.NewErrorf(500, nil, "failed to sign token", err)
	}

	return signedToken, nil

}
func (s *JwtServiceImpl) ValidateToken(token string) (*jwt.Token, error) {
	//
	SecretKey := []byte(s.env.GetWithDefault("ACCESS_SECRET_KEY", "rahasiabanget"))
	decodedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		_, ok := t.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.NewErrorf(400, nil, "validate token error")
		}
		return SecretKey, nil
	})

	if err != nil {
		return decodedToken, err
	}
	return decodedToken, nil
}

func (s *JwtServiceImpl) ConvertTokenToCustomClaims(accessToken *jwt.Token) (*CustomClaims, bool) {
	// konversi dari map[string]interface{} ke *service.CustomClaims
	claims := accessToken.Claims.(jwt.MapClaims)
	payload := &CustomClaims{
		Email: claims["email"].(string),
	}

	// melakukan assertion pada waktu iat
	if val, ok := claims["iat"].(float64); ok {
		payload.IssuedAt = jwt.NewNumericDate(time.Unix(int64(val), 0))
	}

	// melakukan assertion pada waktu exp
	if val, ok := claims["exp"].(float64); ok {
		payload.ExpiresAt = jwt.NewNumericDate(time.Unix(int64(val), 0))
	}

	// melakukan assertion pada string-string lain yang diperlukan
	if val, ok := claims["iss"].(string); ok {
		payload.Issuer = val
	}

	if val, ok := claims["jti"].(string); ok {
		payload.ID = val
	}
	return payload, true
}

func NewAuthService(cfg *config.Config) AuthService {
	return &JwtServiceImpl{
		env: *cfg,
	}
}
