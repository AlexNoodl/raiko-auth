package models

type RegisterRequest struct {
	Username string `json:"username" bson:"username" validate:"required,min=3,max=20"`
	Email    string `json:"email" bson:"email" validate:"required,email"`
	Password string `json:"password" bson:"password" validate:"required,password"`
}

type LoginRequest struct {
	Login    string `json:"login" bson:"login" validate:"required"`
	Password string `json:"password" bson:"password" validate:"required,password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type SuccessResponse struct {
	Message string `json:"message"`
}
