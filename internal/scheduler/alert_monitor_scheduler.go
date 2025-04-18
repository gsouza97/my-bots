package scheduler

import (
	"log"

	"github.com/gsouza97/my-bots/internal/usecase"
	"github.com/robfig/cron/v3"
)

type AlertMonitorScheduler struct {
	checkPriceAlert *usecase.CheckPriceAlert
}

func NewAlertMonitorScheduler(checkPriceAlert *usecase.CheckPriceAlert) *AlertMonitorScheduler {
	return &AlertMonitorScheduler{
		checkPriceAlert: checkPriceAlert,
	}
}

func (s *AlertMonitorScheduler) Start() {
	log.Println("Starting alert monitor scheduler")

	c := cron.New(cron.WithSeconds())

	c.AddFunc("0 */5 * * * *", s.executeCheckPriceAlert)
	c.Start()

	select {}
}

func (s *AlertMonitorScheduler) executeCheckPriceAlert() {
	log.Println("executando check price alert")
	err := s.checkPriceAlert.Execute()
	if err != nil {
		log.Printf("error executing check price alert: %v", err)
	}
}
