package auth

type SignUpRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	PhoneNumber string `json:"phone_number" binding:"required"`
	Nationality string `json:"nationality" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
}


type SignInRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}