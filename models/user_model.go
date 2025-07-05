package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	Id          primitive.ObjectID `json:"id,omitempty" bson:"id,omitempty"`
	AccountId   primitive.ObjectID `json:"accountId,omitempty" bson:"accountId,omitempty"`
	FirstName   string             `json:"firstName" bson:"firstName" validate:"required"`
	LastName    string             `json:"lastName" bson:"lastName" validate:"required"`
	Email       string             `json:"email" bson:"email" validate:"required"`
	Phone       string             `json:"phone" bson:"phone"`
	AvataUrl    string             `json:"avataUrl" bson:"avataUrl"`
	DeviceToken string             `json:"deviceToken" bson:"deviceToken"`
	Ip          string             `json:"ip" bson:"ip"`
	Status      int32              `json:"status" bson:"status"`
	Role        string             `json:"role" bson:"role" validate:"required,eq=ADMIN|eq=USER"`
	Disabled    int32              `json:"disabled" bson:"disabled" validate:"required"`
	DeletedBy   primitive.ObjectID `json:"deletedBy,omitempty" bson:"deletedBy,omitempty"`
	DeletedAt   int64              `json:"deletedAt,omitempty" bson:"deletedAt,omitempty"`
	UpdatedBy   primitive.ObjectID `json:"updatedBy,omitempty" bson:"updatedBy,omitempty"`
	UpdatedAt   int64              `json:"updatedAt,omitempty" bson:"updatedAt,omitempty"`
	CreatedBy   primitive.ObjectID `json:"createdBy,omitempty" bson:"createdBy,omitempty"`
	CreatedAt   int64              `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
}
