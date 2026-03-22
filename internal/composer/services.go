package composer

import (
	"github.com/gsouza97/my-bots/internal/httpserver/service"
)

type ServicesComposer struct {
	AlertsService *service.AlertsService
	LoansService  *service.LoansService
	PoolsService  *service.PoolsService
	AuthService   *service.AuthService
}

func NewServicesComposer(
	repos *RepositoriesComposer,
	providers *ProvidersComposer,
	cfg struct{ UserToken string },
) *ServicesComposer {
	return &ServicesComposer{
		AlertsService: service.NewAlertsService(repos.PriceAlertRepository, providers.BinancePriceProvider),
		LoansService:  service.NewLoansService(repos.LoanRepository, providers.BinancePriceProvider),
		PoolsService:  service.NewPoolsService(repos.PoolRepository, providers.BinancePriceProvider),
		AuthService:   service.NewAuthService(repos.UserRepository, cfg.UserToken),
	}
}
