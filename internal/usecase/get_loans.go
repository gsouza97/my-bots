package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/gsouza97/my-bots/internal/domain"
	"github.com/gsouza97/my-bots/internal/logger"
	"github.com/gsouza97/my-bots/internal/repository"
	"github.com/gsouza97/my-bots/pkg/helper"
)

type GetLoans struct {
	loanRepository repository.LoanRepository
	priceProvider  domain.CryptoPriceProvider
}

func NewGetLoans(loanRepository repository.LoanRepository, priceProvider domain.CryptoPriceProvider) *GetLoans {
	return &GetLoans{
		loanRepository: loanRepository,
		priceProvider:  priceProvider,
	}
}

func (gl *GetLoans) Execute(ctx context.Context) (string, error) {
	t := time.Now()

	loans, err := gl.loanRepository.FindAll(ctx)
	if err != nil {
		return "", err
	}

	assetsSet := make(map[string]bool)
	for _, loan := range loans {
		for _, borrow := range loan.Borrows {
			assetsSet[borrow.Asset] = true
		}
		for _, supply := range loan.Supplies {
			assetsSet[supply.Asset] = true
		}
	}

	var assetsList []string
	for asset := range assetsSet {
		assetsList = append(assetsList, asset)
	}

	assetsPrices, err := gl.priceProvider.GetMultiplePrices(assetsList)
	if err != nil {
		return "", err
	}

	var suppliesBalance float64
	var borrowsBalance float64
	loansMessage := fmt.Sprintf("📌 Empréstimos:\n")

	for _, loan := range loans {
		for _, supply := range loan.Supplies {
			price := assetsPrices[supply.Asset]
			suppliesBalance += float64(supply.Amount) * price
		}
		for _, borrow := range loan.Borrows {
			price := assetsPrices[borrow.Asset]
			borrowsBalance += float64(borrow.Amount) * price
		}

		currentLtv := borrowsBalance / suppliesBalance

		logger.Log.Infof("Total Colateral para %s: %.2f", loan.Description, suppliesBalance)
		logger.Log.Infof("Total Emprestado para %s: %.2f", loan.Description, borrowsBalance)
		logger.Log.Infof("LTV Atual para %s: %.2f%%", loan.Description, currentLtv*100)

		loansMessage += helper.BuildLoansReportMessage(*loan, suppliesBalance, borrowsBalance, currentLtv)

		// Atualizar Base de Dados
		loan.CurrentLtv = currentLtv
		err = gl.loanRepository.Update(ctx, loan)
		if err != nil {
			logger.Log.Errorf("Erro ao atualizar LTV do empréstimo '%s': %s", loan.Description, err.Error())
		}

		suppliesBalance = 0
		borrowsBalance = 0
	}

	t2 := time.Now()

	tFinal := t2.Sub(t)
	logger.Log.Infof("Tempo total get_loans: %s", tFinal)

	return loansMessage, nil
}
