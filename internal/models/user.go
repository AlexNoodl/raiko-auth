package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Role string

const (
	AdminRole Role = "admin"
	UserRole  Role = "user"
)

type User struct {
	ID       primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Email    string             `json:"email" bson:"email" validate:"required,email"`
	Password string             `json:"password" bson:"password" validate:"required,min=8,max=20"`
	Username string             `json:"username" bson:"username" validate:"required,min=3,max=20"`
	IsActive bool               `json:"is_active" bson:"is_active"`
	Role     Role               `json:"role" bson:"role"`
}
