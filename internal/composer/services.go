package composer

import "github.com/gsouza97/my-bots/internal/application/services"

type ServicesComposer struct {
	AlertsService *services.AlertsService
	LoansService  *services.LoansService
	PoolsService  *services.PoolsService
	AuthService   *services.AuthService
	BillsService  *services.BillsService
}

func NewServicesComposer(
	repos *RepositoriesComposer,
	providers *ProvidersComposer,
	cfg struct{ UserToken string },
) *ServicesComposer {
	return &ServicesComposer{
		AlertsService: services.NewAlertsService(repos.PriceAlertRepository, providers.BinancePriceProvider),
		LoansService:  services.NewLoansService(repos.LoanRepository, providers.BinancePriceProvider),
		PoolsService:  services.NewPoolsService(repos.PoolRepository, providers.BinancePriceProvider),
		AuthService:   services.NewAuthService(repos.UserRepository, cfg.UserToken),
		BillsService:  services.NewBillsService(repos.BillRepository),
	}
}
