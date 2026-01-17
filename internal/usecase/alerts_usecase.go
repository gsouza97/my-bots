package usecase

import (
	"context"

	"github.com/gsouza97/my-bots/internal/domain"
	"github.com/gsouza97/my-bots/internal/dto"
	"github.com/gsouza97/my-bots/internal/repository"
)

type AlertsUseCase struct {
	priceAlertRepository repository.PriceAlertRepository
}

func NewAlertsUseCase(priceAlertRepository repository.PriceAlertRepository) *AlertsUseCase {
	return &AlertsUseCase{
		priceAlertRepository: priceAlertRepository,
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
	alert.Active = input.Active

	err = uc.priceAlertRepository.Update(ctx, alert)
	if err != nil {
		return nil, err
	}

	return alert, nil
}
