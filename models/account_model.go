package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Account struct {
	Id           primitive.ObjectID `json:"id,omitempty" bson:"id,omitempty"`
	Username     string             `json:"username" bson:"username" validate:"required,username_valid"`
	Password     string             `json:"password" bson:"password" validate:"required"`
	Role         string             `json:"role" bson:"role" validate:"required,eq=ADMIN|eq=USER"`
	Token        string             `json:"token,omitempty" bson:"token,omitempty"`
	RefreshToken string             `json:"refreshToken,omitempty" bson:"refreshToken,omitempty"`
	ExpiryDate   int64              `json:"expiryDate,omitempty" bson:"expiryDate,omitempty"`
	Disabled     int32              `json:"disabled" bson:"disabled"`
	DeletedBy    primitive.ObjectID `json:"deletedBy,omitempty" bson:"deletedBy,omitempty"`
	DeletedAt    int64              `json:"deletedAt,omitempty" bson:"deletedAt,omitempty"`
	UpdatedBy    primitive.ObjectID `json:"updatedBy,omitempty" bson:"updatedBy,omitempty"`
	UpdatedAt    int64              `json:"updatedAt,omitempty" bson:"updatedAt,omitempty"`
	CreatedBy    primitive.ObjectID `json:"createdBy,omitempty" bson:"createdBy,omitempty"`
	CreatedAt    int64              `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
}

type AccountLogin struct {
	Username string `json:"username"  validate:"required"`
	Password string `json:"password"  validate:"required"`
}

func (u *Account) HidePassWord() {
	u.Password = ""
}
