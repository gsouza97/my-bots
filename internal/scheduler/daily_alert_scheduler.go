package scheduler

import (
	"time"

	"github.com/gsouza97/my-bots/internal/logger"
	"github.com/gsouza97/my-bots/internal/usecase"
	"github.com/robfig/cron/v3"
)

type DailyAlertScheduler struct {
	generateDailyAlert *usecase.GenerateDailyAlert
	cron               string
}

func NewDailyAlertScheduler(generateDailyAlert *usecase.GenerateDailyAlert, cron string) *DailyAlertScheduler {
	return &DailyAlertScheduler{
		generateDailyAlert: generateDailyAlert,
		cron:               cron,
	}
}

func (s *DailyAlertScheduler) Start() {
	logger.Log.Info("Starting daily alert scheduler")

	loc, err := time.LoadLocation("Europe/Madrid")
	if err != nil {
		logger.Log.Fatalf("failed to load location: %v", err)
	}

	c := cron.New(cron.WithSeconds(), cron.WithLocation(loc))

	c.AddFunc(s.cron, s.executeDailyAlert)
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
