package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type HomologacionConfigParams struct {
	ID             primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	DateOfBirth    string             `json:"dateOfBirth" bson:"dateOfBirth"`
	Fullname       string             `json:"fullname" bson:"fullname"`
	DocumentNumber string             `json:"documentNumber" bson:"documentNumber"`
	CurrentStatus  string             `json:"currentStatus" bson:"currentStatus"`
}
