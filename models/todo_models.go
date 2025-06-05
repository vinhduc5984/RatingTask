package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type ToDo struct {
	Id          primitive.ObjectID `json:"id,omitempty" bson:"id,omitempty"`
	Title       string             `json:"title,omitempty" bson:"title,omitempty" validate:"required"`
	Description string             `json:"description,omitempty" bson:"description,omitempty" validate:"required"`
	StartDate   int64              `json:"startDate,omitempty" bson:"startDate,omitempty" validate:"required"`
	EndDate     int64              `json:"endDate,omitempty" bson:"endDate,omitempty" validate:"required"`
	Email       string             `json:"email,omitempty" bson:"email,omitempty" validate:"required,email"`
	DeletedBy   primitive.ObjectID `json:"deletedBy,omitempty" bson:"deletedBy,omitempty"`
	DeletedAt   int64              `json:"deletedAt,omitempty" bson:"deletedAt,omitempty"`
	UpdatedBy   primitive.ObjectID `json:"updatedBy,omitempty" bson:"updatedBy,omitempty"`
	UpdatedAt   int64              `json:"updatedAt,omitempty" bson:"updatedAt,omitempty"`
	CreatedBy   primitive.ObjectID `json:"createdBy,omitempty" bson:"createdBy,omitempty"`
	CreatedAt   int64              `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
}
