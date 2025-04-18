package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type AlertPriceStatus string

const (
	UnderPrice AlertPriceStatus = "UNDER_PRICE"
	OverPrice  AlertPriceStatus = "OVER_PRICE"
)

type PriceAlert struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Crypto      string             `json:"crypto" bson:"crypto"`
	AlertPrice  float64            `json:"alertPrice" bson:"alertPrice"`
	PriceStatus AlertPriceStatus   `json:"priceStatus" bson:"priceStatus"`
	Active      bool               `json:"active" bson:"active"`
}
