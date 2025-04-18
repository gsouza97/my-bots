package domain

import (
	"context"

	"github.com/gsouza97/my-bots/internal/adapter/provider"
)

type PoolFeeProvider interface {
	GetFees(ctx context.Context, chain string, nftId string) (provider.RevertPoolDataResponse, error)
}
