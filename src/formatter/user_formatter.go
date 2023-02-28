package formatter

import "backend-crowdfunding/src/model"

type UserLoginFormatter struct {
	User  UserFormatter `json:"user"`
	Token string        `json:"token"`
}
type UserFormatter struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Avatar     string `json:"avatar"`
	Occupation string `json:"occupation"`
}

func FormatUserLogin(user model.User, token string) UserLoginFormatter {
	userFormat := UserFormatter{
		ID:         user.ID,
		Name:       user.Name,
		Email:      user.Email,
		Avatar:     user.Avatar,
		Occupation: user.Occupation,
	}
	formatter := UserLoginFormatter{
		User:  userFormat,
		Token: token,
	}

	return formatter
}
