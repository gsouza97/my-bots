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

func NewUseCasesComposer(repos *RepositoriesComposer, providers *ProvidersComposer, adapters *AdaptersComposer, cfg *config.Config) *UseCasesComposer {
	uc := &UseCasesComposer{}

	// ========== Providers de negócio (independentes) ==========
	// Estes use cases só precisam de providers, sem repos
	uc.GetFearAndGreedIndex = usecase.NewGetFearAndGreedIndex(providers.FearAndGreedProvider)
	uc.GetAltcoinSeasonIndex = usecase.NewGetAltcoinSeasonIndex(providers.AltcoinSeasonProvider)

	// ========== Bots (que precisam de adapters + outros use cases) ==========
	// Primeiro criamos o bot de price alerts (ele é usado por outros)
	uc.CheckPrice = usecase.NewCheckPrice(repos.PriceAlertRepository, providers.BinancePriceProvider)
	uc.PriceAlertsBot = bot.NewPriceAlertsBot(
		adapters.TelegramPriceAlertsAdapter,
		uc.CheckPrice,
		uc.GetFearAndGreedIndex,
		uc.GetAltcoinSeasonIndex,
		cfg.BotChatID,
	)

	// ========== Use Cases para Bills ==========
	uc.SaveBill = usecase.NewSaveBill(repos.BillRepository)
	uc.GenerateReport = usecase.NewGenerateReport(repos.BillRepository)

	// ========== Use Cases para Price Alerts ==========
	uc.CheckPriceAlert = usecase.NewCheckPriceAlert(
		repos.PriceAlertRepository,
		providers.BinancePriceProvider,
		uc.PriceAlertsBot,
	)

	// ========== Bot de Expenses ==========
	uc.ExpensesBot = bot.NewExpensesBot(adapters.TelegramExpensesAdapter, uc.SaveBill, uc.GenerateReport)

	// ========== Use Cases para Pools ==========
	uc.ListActivePools = usecase.NewListActivePools(repos.PoolRepository)
	uc.GetPoolFees = usecase.NewGetPoolFees(repos.PoolRepository, providers.RevertFeeProvider)

	uc.GenerateDailyAlert = usecase.NewGenerateDailyAlert(
		uc.GetPoolFees,
		uc.GetFearAndGreedIndex,
		uc.GetAltcoinSeasonIndex,
		repos.PriceAlertRepository,
		providers.BinancePriceProvider,
		uc.PriceAlertsBot,
	)

	uc.PoolsBot = bot.NewPoolsBot(adapters.TelegramPoolsAdapter, uc.ListActivePools, uc.GetPoolFees, cfg.BotChatID)
	uc.CheckPools = usecase.NewCheckPools(repos.PoolRepository, providers.BinancePriceProvider, uc.PoolsBot, cfg.NotificationCooldown)

	// ========== Use Cases para Homologacion ==========
	uc.HomologacionBot = bot.NewHomologacionBot(adapters.TelegramHomologacionAdapter, cfg.BotChatID)
	uc.GetHomologacionStatus = usecase.NewGetHomologacionStatus(providers.HomologacionProvider, repos.HomologacionRepository, uc.HomologacionBot)

	// ========== Use Cases para Loans ==========
	uc.GetLoans = usecase.NewGetLoans(repos.LoanRepository, providers.BinancePriceProvider)
	uc.LoansBot = bot.NewLoansBot(adapters.TelegramLoansAdapter, uc.GetLoans, cfg.BotChatID)
	uc.CheckLoans = usecase.NewCheckLoans(repos.LoanRepository, providers.BinancePriceProvider, uc.LoansBot)

	return uc
}
