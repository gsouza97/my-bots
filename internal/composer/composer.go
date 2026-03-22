package composer

import (
	"github.com/gsouza97/my-bots/config"
	"github.com/gsouza97/my-bots/internal/scheduler"
	"go.mongodb.org/mongo-driver/mongo"
)

// Composer é a classe principal que orquestra TODAS as dependências
// Segue o padrão Facade de Dependency Injection
// Centraliza todo o wiring da aplicação em um único lugar
type Composer struct {
	Repositories *RepositoriesComposer
	Providers    *ProvidersComposer
	Adapters     *AdaptersComposer
	UseCases     *UseCasesComposer
	Services     *ServicesComposer
	Handlers     *HandlerComposer
	Routes       *RoutesComposer
	Schedulers   *SchedulersComposer
}

type SchedulersComposer struct {
	AlertMonitorScheduler        *scheduler.AlertMonitorScheduler
	DailyAlertScheduler          *scheduler.DailyAlertScheduler
	PoolsMonitorScheduler        *scheduler.PoolsMonitorScheduler
	HomologacionMonitorScheduler *scheduler.HomologacionMonitorScheduler
	LoansMonitorScheduler        *scheduler.LoansMonitorScheduler
}

func NewSchedulersComposer(useCases *UseCasesComposer, cfg *config.Config) *SchedulersComposer {
	return &SchedulersComposer{
		AlertMonitorScheduler:        scheduler.NewAlertMonitorScheduler(useCases.CheckPriceAlert, cfg.AlertMonitorCron),
		DailyAlertScheduler:          scheduler.NewDailyAlertScheduler(useCases.GenerateDailyAlert, cfg.DailyAlertCron),
		PoolsMonitorScheduler:        scheduler.NewPoolsMonitorScheduler(useCases.CheckPools, cfg.PoolsMonitorCron),
		HomologacionMonitorScheduler: scheduler.NewHomologacionMonitorScheduler(useCases.GetHomologacionStatus, cfg.HomologacionMonitorCron),
		LoansMonitorScheduler:        scheduler.NewLoansMonitorScheduler(useCases.CheckLoans, cfg.LoansMonitorCron),
	}
}

// Ordem é importante:
// 1. Repositories (não dependem de nada)
// 2. Providers (não dependem de nada)
// 3. Adapters (dependem de config)
// 4. UseCases (dependem de repos + providers + adapters)
// 5. Services (dependem de repos + providers)
// 6. Handlers (dependem de services)
// 7. Routes (dependem de handlers)
// 8. Schedulers (dependem de use cases)
func NewComposer(db *mongo.Database, cfg *config.Config) (*Composer, error) {
	// 1. Inicializa Repositories
	repositories := NewRepositoriesComposer(db)

	// 2. Inicializa Providers
	providers := NewProvidersComposer()

	// 3. Inicializa Adapters (pode falhar se tokens estiverem inválidos)
	adapters, err := NewAdaptersComposer(cfg)
	if err != nil {
		return nil, err
	}

	// 4. Inicializa UseCases (dependem de repos + providers + adapters)
	useCases := NewUseCasesComposer(repositories, providers, adapters, cfg)

	// 5. Inicializa Services (dependem de repos + providers)
	services := NewServicesComposer(repositories, providers, struct{ UserToken string }{UserToken: cfg.UserToken})

	// 6. Inicializa Handlers (dependem de services)
	handlers := NewHandlerComposer(services)

	// 7. Inicializa Routes (dependem de handlers)
	routes := NewRoutesComposer(handlers)

	// 8. Inicializa Schedulers (dependem de use cases)
	schedulers := NewSchedulersComposer(useCases, cfg)

	return &Composer{
		Repositories: repositories,
		Providers:    providers,
		Adapters:     adapters,
		UseCases:     useCases,
		Services:     services,
		Handlers:     handlers,
		Routes:       routes,
		Schedulers:   schedulers,
	}, nil
}
