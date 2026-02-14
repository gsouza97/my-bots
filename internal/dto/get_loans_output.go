package dto

import "go.mongodb.org/mongo-driver/bson/primitive"

type LoanItemOutput struct {
	Asset            string  `json:"asset" bson:"asset"`
	Amount           float64 `json:"amount" bson:"amount"`
	CurrentItemValue float64 `json:"currentItemValue" bson:"currentItemValue"`
}

type LoanOutput struct {
	ID               primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Description      string             `json:"description" bson:"description"`
	Supplies         []LoanItemOutput   `json:"supplies" bson:"supplies"`
	TotalSupplyValue float64            `json:"totalSupplyValue" bson:"totalSupplyValue"`
	Borrows          []LoanItemOutput   `json:"borrows" bson:"borrows"`
	TotalBorrowValue float64            `json:"totalBorrowValue" bson:"totalBorrowValue"`
	LiqLtv           float64            `json:"liqLtv" bson:"liqLtv"`
	CurrentLtv       float64            `json:"currentLtv" bson:"currentLtv"`
	AlertRate        float64            `json:"alertRate" bson:"alertRate"`
}
