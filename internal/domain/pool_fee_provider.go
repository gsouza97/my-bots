package domain

import (
	"context"

	"github.com/gsouza97/my-bots/internal/infrastructure/providers"
)

type PoolFeeProvider interface {
	GetFees(ctx context.Context, chain string, nftId string) (providers.RevertPoolDataResponse, error)
}
