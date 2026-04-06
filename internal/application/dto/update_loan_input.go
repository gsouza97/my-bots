package dto

import "github.com/gsouza97/my-bots/internal/domain"

type UpdateLoanInput struct {
	Description string            `json:"description" bson:"description"`
	Supplies    []domain.LoanItem `json:"supplies" bson:"supplies"`
	Borrows     []domain.LoanItem `json:"borrows" bson:"borrows"`
	LiqLtv      float64           `json:"liqLtv" bson:"liqLtv"`
	AlertRate   float64           `json:"alertRate" bson:"alertRate"`
}
