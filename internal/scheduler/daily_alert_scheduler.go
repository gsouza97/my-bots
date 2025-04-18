package scheduler

import (
	"log"

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
	log.Println("Starting daily alert scheduler")

	c := cron.New(cron.WithSeconds())

	c.AddFunc("0 */10 * * * *", s.executeDailyAlert)
	c.Start()

	select {}
}

func (s *DailyAlertScheduler) executeDailyAlert() {
	log.Println("executando daily alert")
	err := s.generateDailyAlert.Execute()
	if err != nil {
		log.Printf("error executing daily alert: %v", err)
	}
}
