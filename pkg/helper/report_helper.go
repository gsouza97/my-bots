package helper

import (
	"fmt"

	"github.com/gsouza97/my-bots/internal/domain"
)

func BuildReport(bills []*domain.Bill, monthStr string) string {
	var total float64
	report := fmt.Sprintf("Relatório de contas para %s:\n", monthStr)

	for _, bill := range bills {
		report += fmt.Sprintf("%s - %.2f €\n", bill.Description, bill.Amount)
		total += bill.Amount
	}

	report += fmt.Sprintf("\nTotal: %.2f €", total)
	return report
}
