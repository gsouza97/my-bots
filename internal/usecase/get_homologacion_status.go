package usecase

import (
	"context"
	"fmt"
	"net/url"

	"github.com/gsouza97/my-bots/internal/adapter/provider"
	"github.com/gsouza97/my-bots/internal/domain"
	"github.com/gsouza97/my-bots/internal/logger"
	"github.com/gsouza97/my-bots/internal/repository"
)

type GetHomologacionStatus struct {
	homologacionProvider   domain.HomologacionProvider
	homologacionRepository repository.HomologacionRepository
	notifier               domain.Notifier
}

func NewGetHomologacionStatus(homologacionProvider domain.HomologacionProvider, homologacionRepository repository.HomologacionRepository, notifier domain.Notifier) *GetHomologacionStatus {
	return &GetHomologacionStatus{
		homologacionProvider:   homologacionProvider,
		homologacionRepository: homologacionRepository,
		notifier:               notifier,
	}
}

func (uc *GetHomologacionStatus) Execute() error {
	ctx := context.Background()
	homologacionConfigParams, err := uc.homologacionRepository.FindOne(ctx)
	if err != nil {
		return err
	}

	currentStatus := homologacionConfigParams.CurrentStatus

	data, err := uc.homologacionProvider.GetHomologacionStatus(buildParams(homologacionConfigParams))
	if err != nil {
		return err
	}
	if len(data) == 0 {
		return fmt.Errorf("nenhum dado retornado do provider")
	}

	logger.Log.Infof("Homologacion Status: %s", data[0].Estado)

	msg := buildHomologacionMessage(data, currentStatus)
	if msg != "" {
		err := uc.notifier.SendMessage(msg)
		if err != nil {
			return fmt.Errorf("error sending notification: %w", err)
		}
	}
	return nil
}

func buildHomologacionMessage(data provider.HomologacionResponse, currentStatus string) string {
	if data[0].Estado != currentStatus {
		return fmt.Sprintf(
			"Estado da homologação alterado para: %s",
			data[0].Estado,
		)
	}
	return ""
}

func buildParams(homologacionConfigParams *domain.HomologacionConfigParams) string {
	values := url.Values{}
	values.Set("fechaNacimiento", homologacionConfigParams.DateOfBirth)
	values.Set("nombreApellidosLogin", homologacionConfigParams.Fullname)
	values.Set("apellidos", "")
	values.Set("numeroDocumentoIdentidadSolicitud", homologacionConfigParams.DocumentNumber)
	values.Set("numeroDocumentoIdentidadLogin", homologacionConfigParams.DocumentNumber)
	values.Set("accesoClave", "true")

	return "?" + values.Encode()
}
