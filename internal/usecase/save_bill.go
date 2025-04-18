package usecase

import (
	"context"
	"fmt"

	"github.com/gsouza97/my-bots/internal/repository"
	"github.com/gsouza97/my-bots/pkg/parser"
)

type SaveBill struct {
	billRepository repository.BillRepository
}

func NewSaveBill(billRepo repository.BillRepository) *SaveBill {
	return &SaveBill{
		billRepository: billRepo,
	}
}

func (sb *SaveBill) Execute(ctx context.Context, message string) (string, error) {
	bill, err := parser.ParseBillMessage(message)
	if err != nil {
		return "", err
	}
	_, err = sb.billRepository.Save(ctx, bill)
	if err != nil {
		return "Erro ao salvar conta", fmt.Errorf("error saving bill: %w", err)
	}

	return "Salvo com sucesso!", nil
}
