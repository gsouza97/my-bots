package scheduler

import (
	"github.com/gsouza97/my-bots/internal/logger"
	"github.com/gsouza97/my-bots/internal/usecase"
	"github.com/robfig/cron/v3"
)

type AlertMonitorScheduler struct {
	checkPriceAlert *usecase.CheckPriceAlert
	cron            string
}

func NewAlertMonitorScheduler(checkPriceAlert *usecase.CheckPriceAlert, cron string) *AlertMonitorScheduler {
	return &AlertMonitorScheduler{
		checkPriceAlert: checkPriceAlert,
		cron:            cron,
	}
}

func (s *AlertMonitorScheduler) Start() {
	logger.Log.Info("Starting alert monitor scheduler")

	c := cron.New(cron.WithSeconds())

	c.AddFunc(s.cron, s.executeCheckPriceAlert)
	c.Start()

	select {}
}

func (s *AlertMonitorScheduler) executeCheckPriceAlert() {
	logger.Log.Info("executando check price alert")
	err := s.checkPriceAlert.Execute()
	if err != nil {
		logger.Log.Errorf("error executing check price alert: %v", err)
	}
}
