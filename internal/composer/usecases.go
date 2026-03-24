package composer

import (
	"github.com/gsouza97/my-bots/config"
	"github.com/gsouza97/my-bots/internal/adapter/bot"
	"github.com/gsouza97/my-bots/internal/application/usecases/bills"
	"github.com/gsouza97/my-bots/internal/application/usecases/dailyalert"
	"github.com/gsouza97/my-bots/internal/application/usecases/homologacion"
	"github.com/gsouza97/my-bots/internal/application/usecases/loans"
	"github.com/gsouza97/my-bots/internal/application/usecases/marketindices"
	"github.com/gsouza97/my-bots/internal/application/usecases/pools"
	"github.com/gsouza97/my-bots/internal/application/usecases/pricealerts"
)

type UseCasesComposer struct {
	// Bots
	PriceAlertsBot  *bot.PriceAlertsBot
	ExpensesBot     *bot.ExpensesBot
	PoolsBot        *bot.PoolsBot
	HomologacionBot *bot.HomologacionBot
	LoansBot        *bot.LoansBot

	// Use Cases para scheduler
	CheckPriceAlert       *pricealerts.CheckPriceAlert
	GenerateDailyAlert    *dailyalert.GenerateDailyAlert
	CheckPools            *pools.CheckPools
	GetHomologacionStatus *homologacion.GetHomologacionStatus
	CheckLoans            *loans.CheckLoans

	// Use Cases
	SaveBill              *bills.SaveBill
	GenerateReport        *bills.GenerateReport
	ListActivePools       *pools.ListActivePools
	GetPoolFees           *pools.GetPoolFees
	GetLoans              *loans.GetLoans
	GetFearAndGreedIndex  *marketindices.GetFearAndGreedIndex
	GetAltcoinSeasonIndex *marketindices.GetAltcoinSeasonIndex
	CheckPrice            *pricealerts.CheckPrice
}

func NewUseCasesComposer(repos *RepositoriesComposer, providers *ProvidersComposer, adapters *AdaptersComposer, eventPublishing *EventPublishingComposer, cfg *config.Config) *UseCasesComposer {
	uc := &UseCasesComposer{}

	// ==================== LEVEL 1: Use Cases Independentes ====================
	// Usam apenas providers (sem repos, sem event publisher)
	uc.GetFearAndGreedIndex = marketindices.NewGetFearAndGreedIndex(providers.FearAndGreedProvider)
	uc.GetAltcoinSeasonIndex = marketindices.NewGetAltcoinSeasonIndex(providers.AltcoinSeasonProvider)

	// ==================== LEVEL 2: Use Cases Simples ====================
	// Usam apenas repos + providers (sem event publisher, sem bots)
	uc.CheckPrice = pricealerts.NewCheckPrice(repos.PriceAlertRepository, providers.BinancePriceProvider)
	uc.SaveBill = bills.NewSaveBill(repos.BillRepository)
	uc.GenerateReport = bills.NewGenerateReport(repos.BillRepository)
	uc.ListActivePools = pools.NewListActivePools(repos.PoolRepository)
	uc.GetPoolFees = pools.NewGetPoolFees(repos.PoolRepository, providers.RevertFeeProvider)
	uc.GetLoans = loans.NewGetLoans(repos.LoanRepository, providers.BinancePriceProvider)

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
	uc.CheckPriceAlert = pricealerts.NewCheckPriceAlert(
		repos.PriceAlertRepository,
		providers.BinancePriceProvider,
		eventPublishing.EventPublisher,
	)
	uc.CheckPools = pools.NewCheckPools(
		repos.PoolRepository,
		providers.BinancePriceProvider,
		eventPublishing.EventPublisher,
		cfg.NotificationCooldown,
	)
	uc.CheckLoans = loans.NewCheckLoans(
		repos.LoanRepository,
		providers.BinancePriceProvider,
		eventPublishing.EventPublisher,
	)
	uc.GetHomologacionStatus = homologacion.NewGetHomologacionStatus(
		providers.HomologacionProvider,
		repos.HomologacionRepository,
		eventPublishing.EventPublisher,
	)
	uc.GenerateDailyAlert = dailyalert.NewGenerateDailyAlert(
		uc.GetPoolFees,
		uc.GetFearAndGreedIndex,
		uc.GetAltcoinSeasonIndex,
		repos.PriceAlertRepository,
		providers.BinancePriceProvider,
		eventPublishing.EventPublisher,
	)

	return uc
}
