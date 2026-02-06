package scheduler

import (
	"github.com/gsouza97/my-bots/internal/logger"
	"github.com/gsouza97/my-bots/internal/usecase"
	"github.com/robfig/cron/v3"
)

type LoansMonitorScheduler struct {
	checkLoans *usecase.CheckLoans
	cron       string
}

func NewLoansMonitorScheduler(checkLoans *usecase.CheckLoans, cron string) *LoansMonitorScheduler {
	return &LoansMonitorScheduler{
		checkLoans: checkLoans,
		cron:       cron,
	}
}

func (s *LoansMonitorScheduler) Start() {
	logger.Log.Info("Starting loans monitor scheduler")

	c := cron.New(cron.WithSeconds())

	c.AddFunc(s.cron, s.executeCheckLoans)
	c.Start()

	select {}
}

func (s *LoansMonitorScheduler) executeCheckLoans() {
	logger.Log.Info("executando check loans")
	err := s.checkLoans.Execute()
	if err != nil {
		logger.Log.Errorf("error executing check loans: %v", err)
	}
}
