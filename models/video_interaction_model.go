package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// type of VideoInteraction
// 1: view
// 2: like
// 3: comment
// 4: share
// 5: watch_time
type VideoInteraction struct {
	Id          primitive.ObjectID 	`json:"id,omitempty"`
	VideoId     primitive.ObjectID 	`json:"videoId,omitempty"`
	UserId      primitive.ObjectID 	`json:"userId,omitempty"`
	Type       	int32             	`json:"type" validate:"required"`
	Value 		float64             `json:"value" validate:"required"`
	CreatedAt   int64              	`json:"CreatedAt,omitempty" validate:"required"`
}