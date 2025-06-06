package main

import (
	"context"
	"time"

	"github.com/gsouza97/my-bots/config"
	"github.com/gsouza97/my-bots/internal/adapter/bot"
	"github.com/gsouza97/my-bots/internal/adapter/provider"
	"github.com/gsouza97/my-bots/internal/httpserver"
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

	// DB
	db := client.Database(cfg.MongoDBName)

	// Repositories
	billRepo := repository.NewBillRepository(db)
	priceAlertRepo := repository.NewPriceAlertRepository(db)
	poolRepo := repository.NewPoolRepository(db)

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

	// Providers
	binanceProvider := provider.NewBinancePriceProvider()
	revertProvider := provider.NewRevertFeeProvider()
	fearAndGreedProvider := provider.NewAlternativeFearAndGreedProvider()
	fearAndGreedUseCase := usecase.NewGetFearAndGreedIndex(fearAndGreedProvider)

	checkPriceUseCase := usecase.NewCheckPrice(priceAlertRepo, binanceProvider)
	priceAlertsBot := bot.NewPriceAlertsBot(telegramPriceAlertsAdapter, checkPriceUseCase, fearAndGreedUseCase, cfg.BotChatID)

	// Use Cases
	saveUseCase := usecase.NewSaveBill(billRepo)
	generateReportUseCase := usecase.NewGenerateReport(billRepo)
	checkPriceAlertUseCase := usecase.NewCheckPriceAlert(priceAlertRepo, binanceProvider, priceAlertsBot)
	listActivePoolsUseCase := usecase.NewListActivePools(poolRepo)
	getPoolFeesUseCase := usecase.NewGetPoolFees(poolRepo, revertProvider)
	generateDailyAlertUseCase := usecase.NewGenerateDailyAlert(getPoolFeesUseCase, fearAndGreedUseCase, priceAlertRepo, binanceProvider, priceAlertsBot)

	poolsBot := bot.NewPoolsBot(telegramPoolsAdapter, listActivePoolsUseCase, getPoolFeesUseCase, cfg.BotChatID)
	checkPoolsUseCase := usecase.NewCheckPools(poolRepo, binanceProvider, poolsBot, cfg.NotificationCooldown)

	// Bots
	expensesBot := bot.NewExpensesBot(telegramExpensesAdapter, saveUseCase, generateReportUseCase)

	// Schedulers
	alertScheduler := scheduler.NewAlertMonitorScheduler(checkPriceAlertUseCase, cfg.AlertMonitorCron)
	dailyAlertScheduler := scheduler.NewDailyAlertScheduler(generateDailyAlertUseCase, cfg.DailyAlertCron)
	poolsMonitorScheduler := scheduler.NewPoolsMonitorScheduler(checkPoolsUseCase, cfg.PoolsMonitorCron)

	// Health Check server
	go httpserver.StartHealthCheckServer()

	// Start
	go priceAlertsBot.Start()
	go expensesBot.Start()
	go poolsBot.Start()
	go dailyAlertScheduler.Start()
	go poolsMonitorScheduler.Start()
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
