package request

type RegisterUserInput struct {
	Name       string `json:"name" binding:"required"`
	Email      string `json:"email" binding:"required,email"`
	Occupation string `json:"occupation" binding:"required"`
	Password   string `json:"password" binding:"required"`
	Avatar     string `json:"avatar" binding:"required"`
	Role       string `json:"role" binding:"required"`
}

type LoginUserInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type CheckEmailInput struct {
	Email string `json:"email" binding:"required,email"`
}

type LoginWithGoogleInput struct {
	FirebaseToken string `json:"firebase_token" binding:"required"`
}

type UserParam struct {
	Name            string `json:"name"`
	Email           string `json:"email"`
	IsGoogleAccount bool   `json:"is_google_account"`
}
