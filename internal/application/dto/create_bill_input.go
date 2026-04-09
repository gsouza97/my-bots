package dto

type CreateBillInput struct {
	Description  string  `json:"description" bson:"description"`
	Amount       float64 `json:"amount" bson:"amount"`
	Category     string  `json:"category" bson:"category"`
	PurchaseDate string  `json:"purchaseDate" bson:"purchaseDate"`
}
