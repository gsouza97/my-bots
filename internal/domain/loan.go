package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type LoanItem struct {
	Asset  string  `json:"asset" bson:"asset"`
	Amount float64 `json:"amount" bson:"amount"`
}

type Loan struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Description string             `json:"description" bson:"description"`
	Supplies    []LoanItem         `json:"supplies" bson:"supplies"`
	Borrows     []LoanItem         `json:"borrows" bson:"borrows"`
	LiqLtv      float64            `json:"liqLtv" bson:"liqLtv"`
	CurrentLtv  float64            `json:"currentLtv" bson:"currentLtv"`
	AlertRate   float64            `json:"alertRate" bson:"alertRate"`
}
