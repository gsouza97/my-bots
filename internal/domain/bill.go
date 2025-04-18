package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Bill struct {
	ID           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Description  string             `json:"description" bson:"description"`
	Amount       float64            `json:"amount" bson:"amount"`
	PurchaseDate time.Time          `json:"purchaseDate" bson:"purchaseDate"`
	Timestamp    time.Time          `json:"timestamp" bson:"timestamp"`
}
