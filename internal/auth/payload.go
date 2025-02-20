package auth

type LoginRequest struct {
	Email string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type RegisterRequest struct {
	Username     string    `json:"username" validate:"required"`
	Email        string    `json:"email" validate:"required,email"`
	Password		 string    `json:"password" validate:"required"`
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	AvatarURL    string    `json:"avatar_url"`
}

type RegisterResponse struct {
	Token string `json:"token"`
}
