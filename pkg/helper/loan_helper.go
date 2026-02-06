package helper

import (
	"fmt"

	"github.com/gsouza97/my-bots/internal/domain"
)

func BuildLoanWarningMessage(pool *domain.Loan, currentLtv float64) string {
	message := fmt.Sprintf("⚠️ATENÇÃO: Empréstimo '%s' próximo de ser liquidado!\nLTV Atual: %.2f%%\nLTV Máximo: %.2f%%", pool.Description, currentLtv*100, pool.LiqLtv*100)
	return message
}

func BuildLoansReportMessage(loan domain.Loan, suppliesBalance float64, borrowsBalance float64, currentLtv float64) string {
	return fmt.Sprintf("\n- Empréstimo: %s\n Total Colateral: %.2f USD\n Total Emprestado: %.2f USD\n LTV Atual: %.2f%%\n LTV Liquidação: %.2f%%\n", loan.Description, suppliesBalance, borrowsBalance, currentLtv*100, loan.LiqLtv*100)
}
