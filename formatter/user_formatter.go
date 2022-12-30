package formatter

import "backend-crowdfunding/src/model"

type UserFormatter struct {
	ID         uint   `json:"id"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Occupation string `json:"occupation"`
	Token      string `json:"token"`
}

func FormatUser(user model.User, token string) UserFormatter {
	formatter := UserFormatter{
		ID:         user.ID,
		Name:       user.Name,
		Email:      user.Email,
		Occupation: user.Occupation,
		Token:      token,
	}

	return formatter
}
