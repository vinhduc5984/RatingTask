package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Creator struct {
	Id          primitive.ObjectID `json:"id,omitempty"`
	Name       	string             `json:"name,omitempty" validate:"required"`
	UserId 		primitive.ObjectID `json:"userId" validate:"required"`
	CreatedAt   int64              `json:"CreatedAt,omitempty" validate:"required"`
}