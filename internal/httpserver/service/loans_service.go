package service

import (
	"context"

	"github.com/gsouza97/my-bots/internal/domain"
	"github.com/gsouza97/my-bots/internal/dto"
	"github.com/gsouza97/my-bots/internal/logger"
	"github.com/gsouza97/my-bots/internal/repository"
	"github.com/gsouza97/my-bots/pkg/helper"
)

type LoansService struct {
	loanRepository repository.LoanRepository
	priceProvider  domain.CryptoPriceProvider
}

func NewLoansService(loanRepository repository.LoanRepository, priceProvider domain.CryptoPriceProvider) *LoansService {
	return &LoansService{
		loanRepository: loanRepository,
		priceProvider:  priceProvider,
	}
}

func (s *LoansService) GetAll() ([]*dto.LoanOutput, error) {
	ctx := context.Background()

	var output []*dto.LoanOutput

	loans, err := s.loanRepository.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	assetsList := helper.ExtractLoansAssets(loans)
	logger.Log.Infof("Unique cryptos in loans: %v", assetsList)

	assetsPrices, err := s.priceProvider.GetMultiplePrices(assetsList)
	if err != nil {
		return nil, err
	}

	logger.Log.Infof("Assets prices: %v", assetsPrices)

	for _, loan := range loans {

		loanOutput := &dto.LoanOutput{
			ID:          loan.ID,
			Description: loan.Description,
			LiqLtv:      loan.LiqLtv,
			CurrentLtv:  loan.CurrentLtv,
			AlertRate:   loan.AlertRate,
		}

		for _, supply := range loan.Supplies {
			supplyOutput := dto.LoanItemOutput{
				Asset:            supply.Asset,
				Amount:           supply.Amount,
				CurrentItemValue: assetsPrices[supply.Asset] * supply.Amount,
			}
			loanOutput.Supplies = append(loanOutput.Supplies, supplyOutput)
		}

		for _, borrow := range loan.Borrows {
			borrowOutput := dto.LoanItemOutput{
				Asset:            borrow.Asset,
				Amount:           borrow.Amount,
				CurrentItemValue: assetsPrices[borrow.Asset] * borrow.Amount,
			}
			loanOutput.Borrows = append(loanOutput.Borrows, borrowOutput)
		}

		for _, supply := range loanOutput.Supplies {
			loanOutput.TotalSupplyValue += supply.CurrentItemValue
		}

		for _, borrow := range loanOutput.Borrows {
			loanOutput.TotalBorrowValue += borrow.CurrentItemValue
		}

		output = append(output, loanOutput)
	}

	// for _, loan := range loans {
	// 	supplyValue, err := calculateSupplyValue(loan.Supplies, s.priceProvider)
	// 	if err != nil {
	// 		logger.Log.Errorf("Error calculating supply value for loan %s: %v", loan.ID, err)
	// 		continue
	// 	}
	// 	loan.SupplyValue = supplyValue
	// }

	return output, nil
}

// func (s *LoansService) UpdateAlert(id string, input dto.UpdatePriceAlertInput) (*domain.PriceAlert, error) {
// 	ctx := context.Background()
// 	alert, err := s.loanRepository.FindById(ctx, id)
// 	if err != nil {
// 		return nil, err
// 	}

// 	alert.Crypto = input.Crypto
// 	alert.AlertPrice = input.AlertPrice

// 	err = s.loanRepository.Update(ctx, alert)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return alert, nil
// }

// func (uc *LoansService) CreateAlert(input dto.CreatePriceAlertInput) (*domain.PriceAlert, error) {
// 	ctx := context.Background()

// 	alert := &domain.PriceAlert{
// 		Crypto:     input.Crypto,
// 		AlertPrice: input.AlertPrice,
// 		Active:     true,
// 	}

// 	price, err := uc.priceProvider.GetPrice(alert.Crypto)
// 	if err != nil {
// 		return nil, err
// 	}

// 	alert.PriceStatus = domain.UnderPrice
// 	if price >= alert.AlertPrice {
// 		alert.PriceStatus = domain.OverPrice
// 	}

// 	createdAlert, err := uc.loanRepository.Create(ctx, alert)
// 	if err != nil {
// 		return nil, err
// 	}

// 	logger.Log.Infof("Alert created: %+v", createdAlert)

// 	return createdAlert, nil
// }

// func createLoansCryptoSet(loans []*domain.Loan) map[string]struct{} {
// 	var loansCryptos []string
// 	for _, loan := range loans {
// 		for _, supplyCrypto := range loan.Supplies {
// 			loansCryptos = append(loansCryptos, supplyCrypto.Asset)
// 		}

// 		for _, borrowCrypto := range loan.Borrows {
// 			loansCryptos = append(loansCryptos, borrowCrypto.Asset)
// 		}
// 	}

// 	return util.CreateSetFromSlice(loansCryptos)
// }

// func calculateSupplyValue(supplies []domain.CryptoAmount, priceProvider domain.CryptoPriceProvider) (float64, error) {
// 	var totalValue float64
// 	for _, supply := range supplies {
// 		price, err := priceProvider.GetPrice(supply.Asset)
// 		if err != nil {
// 			return 0, err
// 		}
// 		totalValue += supply.Amount * price
// 	}
// 	return totalValue, nil
// }
