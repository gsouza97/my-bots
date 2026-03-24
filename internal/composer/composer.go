package composer

import (
	"strconv"

	"github.com/gsouza97/my-bots/config"
	"github.com/gsouza97/my-bots/internal/interfaces/jobs"
	"go.mongodb.org/mongo-driver/mongo"
)

// Composer é a classe principal que orquestra TODAS as dependências
// Segue o padrão Facade de Dependency Injection
// Centraliza todo o wiring da aplicação em um único lugar
type Composer struct {
	Repositories    *RepositoriesComposer
	Providers       *ProvidersComposer
	Adapters        *AdaptersComposer
	EventPublishing *EventPublishingComposer
	UseCases        *UseCasesComposer
	Services        *ServicesComposer
	Handlers        *HandlerComposer
	Routes          *RoutesComposer
	Schedulers      *SchedulersComposer
}

type SchedulersComposer struct {
	AlertMonitorScheduler        *jobs.AlertMonitorScheduler
	DailyAlertScheduler          *jobs.DailyAlertScheduler
	PoolsMonitorScheduler        *jobs.PoolsMonitorScheduler
	HomologacionMonitorScheduler *jobs.HomologacionMonitorScheduler
	LoansMonitorScheduler        *jobs.LoansMonitorScheduler
}

func NewSchedulersComposer(useCases *UseCasesComposer, cfg *config.Config) *SchedulersComposer {
	return &SchedulersComposer{
		AlertMonitorScheduler:        jobs.NewAlertMonitorScheduler(useCases.CheckPriceAlert, cfg.AlertMonitorCron),
		DailyAlertScheduler:          jobs.NewDailyAlertScheduler(useCases.GenerateDailyAlert, cfg.DailyAlertCron),
		PoolsMonitorScheduler:        jobs.NewPoolsMonitorScheduler(useCases.CheckPools, cfg.PoolsMonitorCron),
		HomologacionMonitorScheduler: jobs.NewHomologacionMonitorScheduler(useCases.GetHomologacionStatus, cfg.HomologacionMonitorCron),
		LoansMonitorScheduler:        jobs.NewLoansMonitorScheduler(useCases.CheckLoans, cfg.LoansMonitorCron),
	}
}

// Ordem é importante:
// 1. Repositories (não dependem de nada)
// 2. Providers (não dependem de nada)
// 3. Adapters (dependem de config)
// 4. EventPublisher (depende de adapters para enviar notificações)
// 5. UseCases (dependem de repos + providers + adapters + event publisher)
// 6. Services (dependem de repos + providers)
// 7. Handlers (dependem de services)
// 8. Routes (dependem de handlers)
// 9. Schedulers (dependem de use cases)
func NewComposer(db *mongo.Database, cfg *config.Config) (*Composer, error) {
	chatID, err := strconv.ParseInt(cfg.BotChatID, 10, 64)
	if err != nil {
		return nil, err
	}

	// 1. Inicializa Repositories
	repositories := NewRepositoriesComposer(db)

	// 2. Inicializa Providers
	providers := NewProvidersComposer()

	// 3. Inicializa Adapters (pode falhar se tokens estiverem inválidos)
	adapters, err := NewAdaptersComposer(cfg)
	if err != nil {
		return nil, err
	}

	// 4. Inicializa EventPublisher e registra todos os listeners
	eventPublishing, err := NewEventPublishingComposer(adapters, chatID)
	if err != nil {
		return nil, err
	}

	// 5. Inicializa UseCases (dependem de repos + providers + adapters + event publisher)
	useCases := NewUseCasesComposer(repositories, providers, adapters, eventPublishing, cfg)

	// 6. Inicializa Services (dependem de repos + providers)
	services := NewServicesComposer(repositories, providers, struct{ UserToken string }{UserToken: cfg.UserToken})

	// 7. Inicializa Handlers (dependem de services)
	handlers := NewHandlerComposer(services)

	// 8. Inicializa Routes (dependem de handlers)
	routes := NewRoutesComposer(handlers)

	// 9. Inicializa Schedulers (dependem de use cases)
	schedulers := NewSchedulersComposer(useCases, cfg)

	return &Composer{
		Repositories:    repositories,
		Providers:       providers,
		Adapters:        adapters,
		EventPublishing: eventPublishing,
		UseCases:        useCases,
		Services:        services,
		Handlers:        handlers,
		Routes:          routes,
		Schedulers:      schedulers,
	}, nil
}
