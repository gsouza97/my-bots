package dto

type UpdatePriceAlertInput struct {
	Crypto     string  `json:"crypto" bson:"crypto"`
	AlertPrice float64 `json:"alertPrice" bson:"alertPrice"`
	Active     bool    `json:"active" bson:"active"`
}
