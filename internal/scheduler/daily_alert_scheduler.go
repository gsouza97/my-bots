package scheduler

import (
	"time"

	"github.com/gsouza97/my-bots/internal/logger"
	"github.com/gsouza97/my-bots/internal/usecase"
	"github.com/robfig/cron/v3"
)

type DailyAlertScheduler struct {
	generateDailyAlert *usecase.GenerateDailyAlert
}

func NewDailyAlertScheduler(generateDailyAlert *usecase.GenerateDailyAlert) *DailyAlertScheduler {
	return &DailyAlertScheduler{
		generateDailyAlert: generateDailyAlert,
	}
}

func (s *DailyAlertScheduler) Start() {
	logger.Log.Info("Starting daily alert scheduler")

	loc, err := time.LoadLocation("Europe/Madrid")
	if err != nil {
		logger.Log.Fatalf("failed to load location: %v", err)
	}

	c := cron.New(cron.WithSeconds(), cron.WithLocation(loc))

	c.AddFunc("0 0 8,15,20 * * *", s.executeDailyAlert)
	c.Start()

	select {}
}

func (s *DailyAlertScheduler) executeDailyAlert() {
	logger.Log.Info("executando daily alert")
	err := s.generateDailyAlert.Execute()
	if err != nil {
		logger.Log.Errorf("error executing daily alert: %v", err)
	}
}
