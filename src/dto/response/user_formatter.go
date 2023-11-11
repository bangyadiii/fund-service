package response

import "backend-crowdfunding/src/model"

type UserLoginResponse struct {
	User  UserResponse `json:"user"`
	Token string       `json:"token"`
}

type UserResponse struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Avatar     string `json:"avatar"`
	Occupation string `json:"occupation"`
}

func FormatUserResponse(user *model.User) UserResponse {
	return UserResponse{
		ID:         user.ID,
		Name:       user.Name,
		Email:      user.Email,
		Avatar:     user.Avatar,
		Occupation: user.Occupation,
	}
}

func FormatUserLogin(user *model.User, token string) UserLoginResponse {
	userFormat := FormatUserResponse(user)
	formatter := UserLoginResponse{
		User:  userFormat,
		Token: token,
	}

	return formatter
}
