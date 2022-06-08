package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	Id         primitive.ObjectID `json:"id,omitempty"`
	FirstName  string             `json:"firstName,omitempty"`
	LastName   string             `json:"lastName,omitempty"`
	Country    string             `json:"country,omitempty"`
	Email      string             `json:"email,omitempty" validate:"required"`
	Password   string             `json:"password,omitempty" validate:"required"`
	PublicKey  string             `json:"publicKey,omitempty"`
	PrivateKey string             `json:"privateKey,omitempty"`
}
