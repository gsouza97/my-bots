package service

import (
	"context"

	"github.com/gsouza97/my-bots/internal/domain"
	"github.com/gsouza97/my-bots/internal/repository"
)

type PoolsService struct {
	poolsRepository repository.PoolRepository
	priceProvider   domain.CryptoPriceProvider
}

func NewPoolsService(poolsRepository repository.PoolRepository, priceProvider domain.CryptoPriceProvider) *PoolsService {
	return &PoolsService{
		poolsRepository: poolsRepository,
		priceProvider:   priceProvider,
	}
}

func (s *PoolsService) GetAll() ([]*domain.Pool, error) {
	ctx := context.Background()

	pools, err := s.poolsRepository.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	return pools, nil
}

// func (s *AlertsService) UpdateAlert(id string, input dto.UpdatePriceAlertInput) (*domain.PriceAlert, error) {
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

// func (uc *AlertsService) CreateAlert(input dto.CreatePriceAlertInput) (*domain.PriceAlert, error) {
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
