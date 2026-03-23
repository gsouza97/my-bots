package main

import (
	"context"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gsouza97/my-bots/config"
	"github.com/gsouza97/my-bots/internal/composer"
	"github.com/gsouza97/my-bots/internal/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Inicializa o logger
	logger.Init()
	logger.Log.Info("Iniciando o Projeto My Bots")

	// Configurações
	cfg := config.LoadConfig()

	// Conecta ao MongoDB
	client, err := connectMongoDB(cfg.MongoURI)
	if err != nil {
		logger.Log.Errorf("Erro ao conectar ao MongoDB: %v", err)
		panic(err)
	}
	defer client.Disconnect(context.Background())

	// Composer
	db := client.Database(cfg.MongoDBName)
	comp, err := composer.NewComposer(db, cfg)
	if err != nil {
		logger.Log.Errorf("Erro ao inicializar Composer: %v", err)
		panic(err)
	}

	// Router
	router := gin.New()

	router.Use(gin.LoggerWithConfig(gin.LoggerConfig{
		SkipPaths: []string{"/healthz"},
	}))
	router.Use(gin.Recovery())

	router.Use(cors.New(cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Accept", "Accept-Encoding", "Accept-Language", "Cache-Control", "Connection", "Referer", "User-Agent"},
		AllowOrigins:     cfg.AllowedOrigins,
		AllowCredentials: true,
	}))

	comp.Routes.RegisterRoutes(router, cfg.UserToken)

	// BOTS
	go comp.UseCases.PriceAlertsBot.Start()
	go comp.UseCases.ExpensesBot.Start()
	go comp.UseCases.PoolsBot.Start()
	go comp.UseCases.HomologacionBot.Start()
	go comp.UseCases.LoansBot.Start()

	// SCHEDULERS
	go comp.Schedulers.DailyAlertScheduler.Start()
	go comp.Schedulers.PoolsMonitorScheduler.Start()
	go comp.Schedulers.HomologacionMonitorScheduler.Start()
	go comp.Schedulers.LoansMonitorScheduler.Start()
	go comp.Schedulers.AlertMonitorScheduler.Start()

	router.Run(":8080")
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
