package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type ToDo struct {
	Id          primitive.ObjectID `json:"id,omitempty"`
	Title       string             `json:"title,omitempty" validate:"required"`
	Description string             `json:"description,omitempty" validate:"required"`
	StartDate   int64              `json:"startDate,omitempty" validate:"required"`
	EndDate     int64              `json:"endDate,omitempty" validate:"required"`
	Email 		string 			   `json:"email,omitempty" validate:"required,email"`
}

type ToDo1 struct {
	Id          primitive.ObjectID `json:"id,omitempty"`
	Title       string             `json:"title,omitempty" validate:"required"`
	Description string             `json:"description,omitempty" validate:"required"`
}
