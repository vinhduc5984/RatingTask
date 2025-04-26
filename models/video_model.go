package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Video struct {
	Id          primitive.ObjectID 	`json:"id,omitempty"`
	Title       string             	`json:"title" validate:"required"`
	CreatorId 	primitive.ObjectID  `json:"creatorId" validate:"required"`
	Score		float64				`json:"score,omitempty"`
	CreatedAt   int64              	`json:"CreatedAt,omitempty" validate:"required"`
	
}