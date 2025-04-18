package usecase

import (
	"context"
	"fmt"

	"github.com/gsouza97/my-bots/internal/repository"
	"github.com/gsouza97/my-bots/pkg/helper"
)

type ListActivePools struct {
	poolRepository repository.PoolRepository
}

func NewListActivePools(poolRepo repository.PoolRepository) *ListActivePools {
	return &ListActivePools{
		poolRepository: poolRepo,
	}
}

func (lap *ListActivePools) Execute(ctx context.Context) (string, error) {
	pools, err := lap.poolRepository.FindAllByActiveIsTrue(ctx)
	if err != nil {
		return "", fmt.Errorf("erro ao buscar pools: %w", err)
	}

	if len(pools) == 0 {
		return "Nenhuma pool encontrada", nil
	}

	msg := helper.BuildPoolResponseMessage(pools)

	return msg, nil
}
