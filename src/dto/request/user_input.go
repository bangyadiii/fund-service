package request

type RegisterUserInput struct {
	Name       string `json:"name" validate:"required"`
	Email      string `json:"email" validate:"required,email"`
	Occupation string `json:"occupation" validate:"required"`
	Password   string `json:"password" validate:"required,gt=6"`
	Avatar     string `json:"avatar" validate:"required"`
	Role       string `json:"role" validate:"required"`
}

type LoginUserInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type CheckEmailInput struct {
	Email string `json:"email" validate:"required,email"`
}

type LoginWithGoogleInput struct {
	FirebaseToken string `json:"firebase_token" validate:"required"`
}

type UserParam struct {
	Name            string `json:"name"`
	Email           string `json:"email"`
	IsGoogleAccount bool   `json:"is_google_account"`
}
