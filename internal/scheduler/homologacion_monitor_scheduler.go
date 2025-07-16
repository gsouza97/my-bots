package scheduler

import (
	"github.com/gsouza97/my-bots/internal/logger"
	"github.com/gsouza97/my-bots/internal/usecase"
	"github.com/robfig/cron/v3"
)

type HomologacionMonitorScheduler struct {
	getHomologacionStatus *usecase.GetHomologacionStatus
	cron                  string
}

func NewHomologacionMonitorScheduler(getHomologacionStatus *usecase.GetHomologacionStatus, cron string) *HomologacionMonitorScheduler {
	return &HomologacionMonitorScheduler{
		getHomologacionStatus: getHomologacionStatus,
		cron:                  cron,
	}
}

func (s *HomologacionMonitorScheduler) Start() {
	logger.Log.Info("Starting homologacion monitor scheduler")

	c := cron.New(cron.WithSeconds())

	c.AddFunc(s.cron, s.executeCheckHomologacionStatus)
	c.Start()

	select {}
}

func (s *HomologacionMonitorScheduler) executeCheckHomologacionStatus() {
	logger.Log.Info("executando get homologacion status")
	err := s.getHomologacionStatus.Execute()
	if err != nil {
		logger.Log.Errorf("error executing get homologacion status: %v", err)
	}
}
