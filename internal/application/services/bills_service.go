package services

import (
	"context"

	"github.com/gsouza97/my-bots/internal/domain"
	"github.com/gsouza97/my-bots/internal/infrastructure/repository"
)

type BillsService struct {
	billRepository repository.BillRepository
}

func NewBillsService(billRepository repository.BillRepository) *BillsService {
	return &BillsService{
		billRepository: billRepository,
	}
}

func (s *BillsService) GetAll() ([]*domain.Bill, error) {
	ctx := context.Background()

	bills, err := s.billRepository.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	return bills, nil
}

// func (s *BillsService) UpdateAlert(id string, input dto.UpdatePriceAlertInput) (*domain.PriceAlert, error) {
// 	ctx := context.Background()
// 	alert, err := s.priceAlertRepository.FindById(ctx, id)
// 	if err != nil {
// 		return nil, err
// 	}

// 	alert.Crypto = input.Crypto
// 	alert.AlertPrice = input.AlertPrice

// 	err = s.priceAlertRepository.Update(ctx, alert)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return alert, nil
// }

// func (uc *BillsService) CreateAlert(input dto.CreatePriceAlertInput) (*domain.PriceAlert, error) {
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

// 	createdAlert, err := uc.priceAlertRepository.Create(ctx, alert)
// 	if err != nil {
// 		return nil, err
// 	}

// 	logger.Log.Infof("Alert created: %+v", createdAlert)

// 	return createdAlert, nil
// }

// func (s *BillsService) DeleteAlert(id string) error {
// 	ctx := context.Background()
// 	return s.priceAlertRepository.Delete(ctx, id)
// }
