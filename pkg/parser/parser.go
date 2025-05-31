package parser

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gsouza97/my-bots/internal/domain"
)

func ParseBillMessage(message string) (*domain.Bill, error) {
	parts := strings.Split(message, " ")

	if len(parts) < 2 {
		return nil, errors.New("formato inválido. Use /save descricao valor data(opcional)")
	}

	description := parts[0]
	amount, err := parseAmount(parts[1])
	if err != nil {
		return nil, fmt.Errorf("erro ao processar valor: %s", parts[1])
	}

	purchaseDate := time.Now()
	if len(parts) == 3 {
		purchaseDate, err = time.Parse("02/01/2006", parts[2])
		if err != nil {
			return nil, fmt.Errorf("erro ao processar data: %s", parts[2])
		}
	}

	return &domain.Bill{
		Description:  description,
		Amount:       amount,
		PurchaseDate: purchaseDate,
	}, nil

}

func ParseMonth(monthStr string) (time.Month, error) {
	if monthStr == "" {
		return time.Now().Month(), nil
	}
	monthNames := []string{
		"janeiro", "fevereiro", "março", "abril", "maio", "junho", "julho", "agosto", "setembro", "outubro", "novembro", "dezembro",
	}

	for i, name := range monthNames {
		if strings.EqualFold(name, monthStr) {
			return time.Month(i + 1), nil
		}
	}

	return 0, errors.New("mês inválido")
}

func ParseCheckPriceMesage(message string) (string, error) {
	parts := strings.Split(message, " ")

	if len(parts) < 1 {
		return "", errors.New("formato inválido. Use /price cripto")
	}

	description := parts[0]
	return description, nil
}

func parseAmount(input string) (float64, error) {
	return strconv.ParseFloat(strings.Replace(input, ",", ".", 1), 64)
}
