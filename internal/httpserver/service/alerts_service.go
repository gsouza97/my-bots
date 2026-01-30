package service

import (
	"context"

	"github.com/gsouza97/my-bots/internal/domain"
	"github.com/gsouza97/my-bots/internal/dto"
	"github.com/gsouza97/my-bots/internal/logger"
	"github.com/gsouza97/my-bots/internal/repository"
)

type AlertsService struct {
	priceAlertRepository repository.PriceAlertRepository
	priceProvider        domain.CryptoPriceProvider
}

func NewAlertsService(priceAlertRepository repository.PriceAlertRepository, priceProvider domain.CryptoPriceProvider) *AlertsService {
	return &AlertsService{
		priceAlertRepository: priceAlertRepository,
		priceProvider:        priceProvider,
	}
}

func (s *AlertsService) GetAll() ([]*domain.PriceAlert, error) {
	ctx := context.Background()

	alerts, err := s.priceAlertRepository.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	return alerts, nil
}

func (s *AlertsService) UpdateAlert(id string, input dto.UpdatePriceAlertInput) (*domain.PriceAlert, error) {
	ctx := context.Background()
	alert, err := s.priceAlertRepository.FindById(ctx, id)
	if err != nil {
		return nil, err
	}

	alert.Crypto = input.Crypto
	alert.AlertPrice = input.AlertPrice

	err = s.priceAlertRepository.Update(ctx, alert)
	if err != nil {
		return nil, err
	}

	return alert, nil
}

func (uc *AlertsService) CreateAlert(input dto.CreatePriceAlertInput) (*domain.PriceAlert, error) {
	ctx := context.Background()

	alert := &domain.PriceAlert{
		Crypto:     input.Crypto,
		AlertPrice: input.AlertPrice,
		Active:     true,
	}

	price, err := uc.priceProvider.GetPrice(alert.Crypto)
	if err != nil {
		return nil, err
	}

	alert.PriceStatus = domain.UnderPrice
	if price >= alert.AlertPrice {
		alert.PriceStatus = domain.OverPrice
	}

	createdAlert, err := uc.priceAlertRepository.Create(ctx, alert)
	if err != nil {
		return nil, err
	}

	logger.Log.Infof("Alert created: %+v", createdAlert)

	return createdAlert, nil
}
