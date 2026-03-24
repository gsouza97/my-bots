package dto

type CreatePriceAlertInput struct {
	Crypto     string  `json:"crypto" bson:"crypto"`
	AlertPrice float64 `json:"alertPrice" bson:"alertPrice"`
}
