package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Id           primitive.ObjectID `json:"id" bson:"_id"`
	FirstName    string             `json:"firstName" bson:"firstName" validate:"required"`
	LastName     string             `json:"lastName" bson:"lastName" validate:"required"`
	Mobile       string             `json:"mobile" bson:"mobile" validate:"required,min=10,max=10"`
	Email        string             `json:"email" bson:"email" validate:"required,email"`
	Password     string             `json:"password" bson:"password" validate:"required,min=5,max=50"`
	Token        string             `json:"token" bson:"token"`
	RefreshToken string             `json:"refreshToken" bson:"refreshToken"`
	IsActive     bool               `json:"isActive" bson:"isActive"`
	CreatedAt    primitive.DateTime `json:"createdAt" bson:"createdAt"`
	UpdatedAt    primitive.DateTime `json:"updatedAt" bson:"updatedAt"`
}

type Response struct {
	StatusCode int         `json:"errorCode"`
	Msg        string      `json:"msg"`
	Data       interface{} `json:"data"`
}

type LoginRequest struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

type JwtTokenExp struct {
	Exp int64 `json:"exp"`
}
