package service

import (
	"context"

	"github.com/gsouza97/my-bots/internal/domain"
	"github.com/gsouza97/my-bots/internal/repository"
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

func (s *LoansService) GetAll() ([]*domain.Loan, error) {
	ctx := context.Background()

	loans, err := s.loanRepository.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	return loans, nil
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
