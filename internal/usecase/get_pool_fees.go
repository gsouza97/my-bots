package usecase

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/gsouza97/my-bots/internal/domain"
	"github.com/gsouza97/my-bots/internal/repository"
	"github.com/gsouza97/my-bots/pkg/helper"
)

type GetPoolFees struct {
	poolRepository repository.PoolRepository
	feesProvider   domain.PoolFeeProvider
}

func NewGetPoolFees(poolRepository repository.PoolRepository, feesProvider domain.PoolFeeProvider) *GetPoolFees {
	return &GetPoolFees{
		poolRepository: poolRepository,
		feesProvider:   feesProvider,
	}
}

func (gpf *GetPoolFees) ExecuteAndUpdateLastFees(ctx context.Context) (string, error) {
	return gpf.execute(ctx, true)
}

func (gpf *GetPoolFees) Execute(ctx context.Context) (string, error) {
	return gpf.execute(ctx, false)
}

func (gpf *GetPoolFees) execute(ctx context.Context, updateLastFeesAmount bool) (string, error) {
	pools, err := gpf.poolRepository.FindAllByActiveIsTrue(ctx)
	if err != nil {
		return "", err
	}

	t := time.Now()

	numWorkers := 8
	poolChannel := make(chan *domain.Pool)
	messageChannel := make(chan string, len(pools))
	totalFeesChannel := make(chan float64, len(pools))

	var wg sync.WaitGroup

	// Inicializar os workers
	for i := 0; i < numWorkers; i++ {

		wg.Add(1)
		go gpf.processPool(ctx, &wg, poolChannel, messageChannel, totalFeesChannel, updateLastFeesAmount)
	}

	// Enviar pools para os workers processarem
	go func() {
		for _, pool := range pools {
			poolChannel <- pool
		}
		close(poolChannel)
	}()

	wg.Wait()
	close(messageChannel)
	close(totalFeesChannel)

	t2 := time.Now()

	tFinal := t2.Sub(t)
	log.Printf("Tempo total get_pool_fees: %s", tFinal)

	totalFeesToCollectSum := 0.0
	feesToCollectMessages := []string{"📌 Fees para coletar:"}

	for msg := range messageChannel {
		feesToCollectMessages = append(feesToCollectMessages, msg)
	}

	for totalFees := range totalFeesChannel {
		totalFeesToCollectSum += totalFees
	}

	feesToCollectMessages = append(feesToCollectMessages, fmt.Sprintf("\nTotal: %.2f USD", totalFeesToCollectSum))

	return strings.Join(feesToCollectMessages, "\n"), nil
}

func (gpf *GetPoolFees) processPool(ctx context.Context, wg *sync.WaitGroup, poolChannel chan *domain.Pool, messageChannel chan string, totalFeesChannel chan float64, updateLastFeesAmount bool) {
	defer wg.Done()
	for pool := range poolChannel {

		data, err := gpf.feesProvider.GetFees(ctx, pool.Chain, pool.NftId)
		if err != nil {
			log.Printf("erro ao buscar fees para pool %s: %s", pool.Description, err.Error())
			continue
		}

		totalFeesToCollect, err := helper.CalculateFeesToCollect(data)
		if err != nil {
			log.Printf("erro ao calcular fees para pool %s: %s", pool.Description, err.Error())
			continue
		}

		// Atualiza o valor total de fees a coletar no BBDD
		if updateLastFeesAmount {
			log.Printf("atualizando lastFeesAmount para pool %s: %.2f", pool.Description, totalFeesToCollect)
			pool.LastFeesAmount = totalFeesToCollect
			err = gpf.poolRepository.Update(ctx, pool)
			if err != nil {
				log.Printf("erro ao atualizar fees para pool %s: %s", pool.Description, err.Error())
				continue
			}
		}

		feesToCollectVariation := totalFeesToCollect - pool.LastFeesAmount
		msg := helper.BuildFeesToCollectMessage(pool, totalFeesToCollect, feesToCollectVariation)
		messageChannel <- msg
		totalFeesChannel <- totalFeesToCollect
	}

}
