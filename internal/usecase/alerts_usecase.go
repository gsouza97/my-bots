package usecase

import (
	"context"

	"github.com/gsouza97/my-bots/internal/domain"
	"github.com/gsouza97/my-bots/internal/dto"
	"github.com/gsouza97/my-bots/internal/logger"
	"github.com/gsouza97/my-bots/internal/repository"
)

type AlertsUseCase struct {
	priceAlertRepository repository.PriceAlertRepository
	priceProvider        domain.CryptoPriceProvider
}

func NewAlertsUseCase(priceAlertRepository repository.PriceAlertRepository, priceProvider domain.CryptoPriceProvider) *AlertsUseCase {
	return &AlertsUseCase{
		priceAlertRepository: priceAlertRepository,
		priceProvider:        priceProvider,
	}
}

func (uc *AlertsUseCase) GetAll() ([]*domain.PriceAlert, error) {
	ctx := context.Background()

	alerts, err := uc.priceAlertRepository.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	return alerts, nil
}

func (uc *AlertsUseCase) UpdateAlert(id string, input dto.UpdatePriceAlertInput) (*domain.PriceAlert, error) {
	ctx := context.Background()
	alert, err := uc.priceAlertRepository.FindById(ctx, id)
	if err != nil {
		return nil, err
	}

	alert.Crypto = input.Crypto
	alert.AlertPrice = input.AlertPrice

	err = uc.priceAlertRepository.Update(ctx, alert)
	if err != nil {
		return nil, err
	}

	return alert, nil
}

func (uc *AlertsUseCase) CreateAlert(input dto.CreatePriceAlertInput) (*domain.PriceAlert, error) {
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
