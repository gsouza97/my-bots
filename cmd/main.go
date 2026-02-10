package main

import (
	"context"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gsouza97/my-bots/config"
	"github.com/gsouza97/my-bots/internal/adapter/bot"
	"github.com/gsouza97/my-bots/internal/adapter/provider"
	"github.com/gsouza97/my-bots/internal/handlers"
	"github.com/gsouza97/my-bots/internal/httpserver/routes"
	"github.com/gsouza97/my-bots/internal/httpserver/service"
	"github.com/gsouza97/my-bots/internal/logger"
	"github.com/gsouza97/my-bots/internal/repository"
	"github.com/gsouza97/my-bots/internal/scheduler"
	"github.com/gsouza97/my-bots/internal/usecase"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Inicializa o logger
	logger.Init()
	logger.Log.Info("Iniciando o Projeto My Bots")

	cfg := config.LoadConfig()

	client, err := connectMongoDB(cfg.MongoURI)
	if err != nil {
		logger.Log.Error("Erro ao conectar ao MongoDB: ", err)
		panic(err)
	}
	defer client.Disconnect(context.Background())

	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Accept", "Accept-Encoding", "Accept-Language", "Cache-Control", "Connection", "Referer", "User-Agent"},
		AllowOrigins:     cfg.AllowedOrigins,
		AllowCredentials: true,
	}))

	// DB
	db := client.Database(cfg.MongoDBName)

	// Repositories
	billRepo := repository.NewBillRepository(db)
	priceAlertRepo := repository.NewPriceAlertRepository(db)
	poolRepo := repository.NewPoolRepository(db)
	homologacionRepo := repository.NewHomologacionRepository(db)
	userRepo := repository.NewUserRepository(db)
	loanRepo := repository.NewLoanRepository(db)

	// Adapters
	telegramExpensesAdapter, err := bot.NewTelegramAdapter(cfg.ExpensesBotToken)
	if err != nil {
		panic(err)
	}
	telegramPriceAlertsAdapter, err := bot.NewTelegramAdapter(cfg.PriceAlertsBotToken)
	if err != nil {
		panic(err)
	}
	telegramPoolsAdapter, err := bot.NewTelegramAdapter(cfg.PoolsBotToken)
	if err != nil {
		panic(err)
	}
	telegramHomologacionAdapter, err := bot.NewTelegramAdapter(cfg.HomologacionBotToken)
	if err != nil {
		panic(err)
	}
	telegramLoansAdapter, err := bot.NewTelegramAdapter(cfg.LoansBotToken)
	if err != nil {
		panic(err)
	}

	// Providers
	binanceProvider := provider.NewBinancePriceProvider()
	revertProvider := provider.NewRevertFeeProvider()
	fearAndGreedProvider := provider.NewAlternativeFearAndGreedProvider()
	altcoinSeasonProvider := provider.NewCmcAltcoinSeasonProvider()
	homologacionProvider := provider.NewHomologacionProvider()

	fearAndGreedUseCase := usecase.NewGetFearAndGreedIndex(fearAndGreedProvider)
	getAltcoinSeasonUseCase := usecase.NewGetAltcoinSeasonIndex(altcoinSeasonProvider)
	checkPriceUseCase := usecase.NewCheckPrice(priceAlertRepo, binanceProvider)
	priceAlertsBot := bot.NewPriceAlertsBot(telegramPriceAlertsAdapter, checkPriceUseCase, fearAndGreedUseCase, getAltcoinSeasonUseCase, cfg.BotChatID)

	// Use Cases
	saveUseCase := usecase.NewSaveBill(billRepo)
	generateReportUseCase := usecase.NewGenerateReport(billRepo)
	checkPriceAlertUseCase := usecase.NewCheckPriceAlert(priceAlertRepo, binanceProvider, priceAlertsBot)
	listActivePoolsUseCase := usecase.NewListActivePools(poolRepo)
	getPoolFeesUseCase := usecase.NewGetPoolFees(poolRepo, revertProvider)
	generateDailyAlertUseCase := usecase.NewGenerateDailyAlert(getPoolFeesUseCase, fearAndGreedUseCase, getAltcoinSeasonUseCase, priceAlertRepo, binanceProvider, priceAlertsBot)
	getLoansUseCase := usecase.NewGetLoans(loanRepo, binanceProvider)

	homologacionBot := bot.NewHomologacionBot(telegramHomologacionAdapter, cfg.BotChatID)
	getHomologacionStatusUseCase := usecase.NewGetHomologacionStatus(homologacionProvider, homologacionRepo, homologacionBot)

	poolsBot := bot.NewPoolsBot(telegramPoolsAdapter, listActivePoolsUseCase, getPoolFeesUseCase, cfg.BotChatID)
	checkPoolsUseCase := usecase.NewCheckPools(poolRepo, binanceProvider, poolsBot, cfg.NotificationCooldown)

	loansBot := bot.NewLoansBot(telegramLoansAdapter, getLoansUseCase, cfg.BotChatID)
	checkLoansUseCase := usecase.NewCheckLoans(loanRepo, binanceProvider, loansBot)

	// Bots
	expensesBot := bot.NewExpensesBot(telegramExpensesAdapter, saveUseCase, generateReportUseCase)

	// Schedulers
	alertScheduler := scheduler.NewAlertMonitorScheduler(checkPriceAlertUseCase, cfg.AlertMonitorCron)
	dailyAlertScheduler := scheduler.NewDailyAlertScheduler(generateDailyAlertUseCase, cfg.DailyAlertCron)
	poolsMonitorScheduler := scheduler.NewPoolsMonitorScheduler(checkPoolsUseCase, cfg.PoolsMonitorCron)
	homologacionMonitorScheduler := scheduler.NewHomologacionMonitorScheduler(getHomologacionStatusUseCase, cfg.HomologacionMonitorCron)
	loansMonitorScheduler := scheduler.NewLoansMonitorScheduler(checkLoansUseCase, cfg.LoansMonitorCron)

	// Services
	alertsService := service.NewAlertsService(priceAlertRepo, binanceProvider)
	authService := service.NewAuthService(userRepo, cfg.UserToken)
	loansService := service.NewLoansService(loanRepo, binanceProvider)
	poolsService := service.NewPoolsService(poolRepo, binanceProvider)

	// Handlers
	alertsHandler := handlers.NewAlertsHandler(alertsService)
	loginHandler := handlers.NewLoginHandler(authService)
	loansHandler := handlers.NewLoansHandler(loansService)
	poolsHandler := handlers.NewPoolsHandler(poolsService)

	// Routes
	alertRoutes := routes.NewAlertsRoutes(alertsHandler)
	loginRoutes := routes.NewLoginRoutes(loginHandler)
	loanRoutes := routes.NewLoansRoutes(loansHandler)
	poolsRoutes := routes.NewPoolsRoutes(poolsHandler)

	// Servers
	alertRoutes.StartAlertsRoutes(router)
	routes.StartHealthRoutes(router)
	loginRoutes.StartLoginRoutes(router)
	loanRoutes.StartLoansRoutes(router)
	poolsRoutes.StartPoolsRoutes(router)
	go router.Run(":8080")

	// Start
	go priceAlertsBot.Start()
	go expensesBot.Start()
	go poolsBot.Start()
	go homologacionBot.Start()
	go dailyAlertScheduler.Start()
	go poolsMonitorScheduler.Start()
	go homologacionMonitorScheduler.Start()
	go loansMonitorScheduler.Start()
	go loansBot.Start()
	alertScheduler.Start()

}

func connectMongoDB(uri string) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	// Verifica a conexão
	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	logger.Log.Info("Conectado ao MongoDB com sucesso!")
	return client, nil
}
