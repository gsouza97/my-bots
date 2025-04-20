package scheduler

import (
	"github.com/gsouza97/my-bots/internal/logger"
	"github.com/gsouza97/my-bots/internal/usecase"
	"github.com/robfig/cron/v3"
)

type PoolsMonitorScheduler struct {
	checkPools *usecase.CheckPools
}

func NewPoolsMonitorScheduler(checkPools *usecase.CheckPools) *PoolsMonitorScheduler {
	return &PoolsMonitorScheduler{
		checkPools: checkPools,
	}
}

func (s *PoolsMonitorScheduler) Start() {
	logger.Log.Info("Starting pool monitor scheduler")

	c := cron.New(cron.WithSeconds())

	c.AddFunc("0 */10 * * * *", s.executeCheckPools)
	c.Start()

	select {}
}

func (s *PoolsMonitorScheduler) executeCheckPools() {
	logger.Log.Info("executando check pools")
	err := s.checkPools.Execute()
	if err != nil {
		logger.Log.Errorf("error executing check pools: %v", err)
	}
}
