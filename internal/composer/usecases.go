package composer

import (
	"github.com/gsouza97/my-bots/config"
	"github.com/gsouza97/my-bots/internal/adapter/bot"
	"github.com/gsouza97/my-bots/internal/usecase"
)

type UseCasesComposer struct {
	// Bots
	PriceAlertsBot  *bot.PriceAlertsBot
	ExpensesBot     *bot.ExpensesBot
	PoolsBot        *bot.PoolsBot
	HomologacionBot *bot.HomologacionBot
	LoansBot        *bot.LoansBot

	// Use Cases para scheduler
	CheckPriceAlert       *usecase.CheckPriceAlert
	GenerateDailyAlert    *usecase.GenerateDailyAlert
	CheckPools            *usecase.CheckPools
	GetHomologacionStatus *usecase.GetHomologacionStatus
	CheckLoans            *usecase.CheckLoans

	// Use Cases
	SaveBill              *usecase.SaveBill
	GenerateReport        *usecase.GenerateReport
	ListActivePools       *usecase.ListActivePools
	GetPoolFees           *usecase.GetPoolFees
	GetLoans              *usecase.GetLoans
	GetFearAndGreedIndex  *usecase.GetFearAndGreedIndex
	GetAltcoinSeasonIndex *usecase.GetAltcoinSeasonIndex
	CheckPrice            *usecase.CheckPrice
}

func NewUseCasesComposer(repos *RepositoriesComposer, providers *ProvidersComposer, adapters *AdaptersComposer, eventPublishing *EventPublishingComposer, cfg *config.Config) *UseCasesComposer {
	uc := &UseCasesComposer{}

	// ==================== LEVEL 1: Use Cases Independentes ====================
	// Usam apenas providers (sem repos, sem event publisher)
	uc.GetFearAndGreedIndex = usecase.NewGetFearAndGreedIndex(providers.FearAndGreedProvider)
	uc.GetAltcoinSeasonIndex = usecase.NewGetAltcoinSeasonIndex(providers.AltcoinSeasonProvider)

	// ==================== LEVEL 2: Use Cases Simples ====================
	// Usam apenas repos + providers (sem event publisher, sem bots)
	uc.CheckPrice = usecase.NewCheckPrice(repos.PriceAlertRepository, providers.BinancePriceProvider)
	uc.SaveBill = usecase.NewSaveBill(repos.BillRepository)
	uc.GenerateReport = usecase.NewGenerateReport(repos.BillRepository)
	uc.ListActivePools = usecase.NewListActivePools(repos.PoolRepository)
	uc.GetPoolFees = usecase.NewGetPoolFees(repos.PoolRepository, providers.RevertFeeProvider)
	uc.GetLoans = usecase.NewGetLoans(repos.LoanRepository, providers.BinancePriceProvider)

	// ==================== LEVEL 3: Bots ====================
	// Usam adapters + use cases independentes
	uc.PriceAlertsBot = bot.NewPriceAlertsBot(
		adapters.TelegramPriceAlertsAdapter,
		uc.CheckPrice,
		uc.GetFearAndGreedIndex,
		uc.GetAltcoinSeasonIndex,
		cfg.BotChatID,
	)
	uc.ExpensesBot = bot.NewExpensesBot(
		adapters.TelegramExpensesAdapter,
		uc.SaveBill,
		uc.GenerateReport,
	)
	uc.PoolsBot = bot.NewPoolsBot(
		adapters.TelegramPoolsAdapter,
		uc.ListActivePools,
		uc.GetPoolFees,
		cfg.BotChatID,
	)
	uc.HomologacionBot = bot.NewHomologacionBot(
		adapters.TelegramHomologacionAdapter,
		cfg.BotChatID,
	)
	uc.LoansBot = bot.NewLoansBot(
		adapters.TelegramLoansAdapter,
		uc.GetLoans,
		cfg.BotChatID,
	)

	// ==================== LEVEL 4: Use Cases com Event Publishing ====================
	// Usam event publisher para notificar via listeners (ex: Telegram)
	uc.CheckPriceAlert = usecase.NewCheckPriceAlert(
		repos.PriceAlertRepository,
		providers.BinancePriceProvider,
		eventPublishing.EventPublisher,
	)
	uc.CheckPools = usecase.NewCheckPools(
		repos.PoolRepository,
		providers.BinancePriceProvider,
		eventPublishing.EventPublisher,
		cfg.NotificationCooldown,
	)
	uc.CheckLoans = usecase.NewCheckLoans(
		repos.LoanRepository,
		providers.BinancePriceProvider,
		eventPublishing.EventPublisher,
	)
	uc.GetHomologacionStatus = usecase.NewGetHomologacionStatus(
		providers.HomologacionProvider,
		repos.HomologacionRepository,
		eventPublishing.EventPublisher,
	)
	uc.GenerateDailyAlert = usecase.NewGenerateDailyAlert(
		uc.GetPoolFees,
		uc.GetFearAndGreedIndex,
		uc.GetAltcoinSeasonIndex,
		repos.PriceAlertRepository,
		providers.BinancePriceProvider,
		eventPublishing.EventPublisher,
	)

	return uc
}
