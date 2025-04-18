package scheduler

import (
	"log"

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
	log.Println("Starting pool monitor scheduler")

	c := cron.New(cron.WithSeconds())

	c.AddFunc("0 */10 * * * *", s.executeCheckPools)
	c.Start()

	select {}
}

func (s *PoolsMonitorScheduler) executeCheckPools() {
	log.Println("executando check pools")
	err := s.checkPools.Execute()
	if err != nil {
		log.Printf("error executing check pools: %v", err)
	}
}
