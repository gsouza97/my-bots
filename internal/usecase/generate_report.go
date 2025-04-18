package usecase

import (
	"context"
	"fmt"

	"github.com/gsouza97/my-bots/internal/repository"
	"github.com/gsouza97/my-bots/pkg/helper"
	"github.com/gsouza97/my-bots/pkg/parser"
)

type GenerateReport struct {
	billRepository repository.BillRepository
}

func NewGenerateReport(billRepo repository.BillRepository) *GenerateReport {
	return &GenerateReport{
		billRepository: billRepo,
	}
}

func (gr *GenerateReport) Execute(ctx context.Context, monthString string) (string, error) {
	month, err := parser.ParseMonth(monthString)
	if err != nil {
		return "Mês inválido", nil
	}

	bills, err := gr.billRepository.FindByMonth(ctx, month)
	if err != nil {
		return "", fmt.Errorf("erro ao buscar contas: %w", err)
	}

	if len(bills) == 0 {
		return "Nenhuma conta encontrada", nil
	}

	report := helper.BuildReport(bills, monthString)
	return report, nil
}
