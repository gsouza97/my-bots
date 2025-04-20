package usecase

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/gsouza97/my-bots/internal/domain"
	"github.com/gsouza97/my-bots/internal/logger"
	"github.com/gsouza97/my-bots/internal/repository"
	"github.com/gsouza97/my-bots/pkg/helper"
)

type CheckPools struct {
	poolRepository       repository.PoolRepository
	priceProvider        domain.CryptoPriceProvider
	notifier             domain.Notifier
	notificationCooldown string
}

func NewCheckPools(poolRepository repository.PoolRepository, priceProvider domain.CryptoPriceProvider, notifier domain.Notifier, notificationCooldown string) *CheckPools {
	return &CheckPools{
		poolRepository:       poolRepository,
		priceProvider:        priceProvider,
		notifier:             notifier,
		notificationCooldown: notificationCooldown,
	}
}

func (cp *CheckPools) Execute() error {
	ctx := context.Background()
	pools, err := cp.poolRepository.FindAllByActiveIsTrue(ctx)
	if err != nil {
		return err
	}

	numWorkers := 8
	poolChannel := make(chan *domain.Pool, numWorkers)
	errChannel := make(chan error, 1)
	t := time.Now()
	var wg sync.WaitGroup

	for _, pool := range pools {
		wg.Add(1)

		poolChannel <- pool

		go func(pool *domain.Pool) {
			defer wg.Done()
			defer func() { <-poolChannel }()

			err := cp.processPool(ctx, pool)
			if err != nil {
				select {
				case errChannel <- err:
				default:
				}
			}
		}(pool)

	}
	wg.Wait()

	t2 := time.Now()
	logger.Log.Infof("Tempo total check_pools: %s", t2.Sub(t))

	select {
	case err := <-errChannel:
		return err
	default:
		return nil
	}
}

func (cp *CheckPools) processPool(ctx context.Context, pool *domain.Pool) error {
	cryptoPool := pool.Crypto1 + pool.Crypto2
	price, err := cp.priceProvider.GetPrice(pool.Crypto1, pool.Crypto2)
	if err != nil {
		return fmt.Errorf("error getting price for %s: %w", cryptoPool, err)
	}
	logger.Log.Infof("price for %s: %f", cryptoPool, price)

	percentToMax := helper.CalculatePercentToMaxPrice(price, pool.MaxPrice)
	percentToMin := helper.CalculatePercentToMinPrice(price, pool.MinPrice)

	logger.Log.Infof("percentToMax: %f", percentToMax)
	logger.Log.Infof("percentToMin: %f", percentToMin)

	outOfRange := price > pool.MaxPrice || price < pool.MinPrice

	// Se o status mudou, manda mensagem
	if pool.OutOfRange != outOfRange {
		logger.Log.Infof("Pool %s out of range: %t", cryptoPool, outOfRange)
		pool.OutOfRange = outOfRange
		err := cp.poolRepository.Update(ctx, pool)
		if err != nil {
			return fmt.Errorf("error updating pool: %w", err)
		}

		message := helper.BuildRangeMessage(pool, outOfRange, price)
		cp.notifier.SendMessage(message)
	}

	if !outOfRange {
		notificationCooldownFloat, err := strconv.ParseFloat(cp.notificationCooldown, 64)
		if err != nil {
			return fmt.Errorf("error parsing notification cooldown: %w", err)
		}

		cooldown := time.Since(pool.LastNotificationTime).Seconds()

		if cooldown > notificationCooldownFloat {

			if percentToMax < pool.RiskRate {
				logger.Log.Infof("Risk rate reached. Sending notification.")
				maxWarningMessage := true
				diffToMax := price * percentToMax
				msgMax := helper.BuildWarningMessage(pool, price, percentToMax, diffToMax, maxWarningMessage)
				cp.notifier.SendMessage(msgMax)
				pool.LastNotificationTime = time.Now()
				cp.poolRepository.Update(ctx, pool)
			}

			if percentToMin < pool.RiskRate {
				logger.Log.Infof("Risk rate reached. Sending notification.")
				maxWarningMessage := false
				diffToMin := price * percentToMin
				msgMin := helper.BuildWarningMessage(pool, price, percentToMax, diffToMin, maxWarningMessage)
				cp.notifier.SendMessage(msgMin)
				pool.LastNotificationTime = time.Now()
				cp.poolRepository.Update(ctx, pool)
			}

		} else {
			logger.Log.Infof("Skipping notification. Cooldown time not reached yet.")
			return nil
		}
	}

	return nil
}
