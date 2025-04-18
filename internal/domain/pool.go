package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Pool struct {
	ID                   primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Description          string             `json:"description" bson:"description"`
	Crypto1              string             `json:"crypto1" bson:"crypto1"`
	Crypto2              string             `json:"crypto2" bson:"crypto2"`
	MinPrice             float64            `json:"minPrice" bson:"minPrice"`
	MaxPrice             float64            `json:"maxPrice" bson:"maxPrice"`
	RiskRate             float64            `json:"riskRate" bson:"riskRate"`
	OutOfRange           bool               `json:"outOfRange" bson:"outOfRange"`
	LastNotificationTime time.Time          `json:"lastNotificationTime" bson:"lastNotificationTime"`
	Chain                string             `json:"chain" bson:"chain"`
	NftId                string             `json:"nftId" bson:"nftId"`
	Active               bool               `json:"active" bson:"active"`
	LastFeesAmount       float64            `json:"lastFeesAmount" bson:"lastFeesAmount"`
}
