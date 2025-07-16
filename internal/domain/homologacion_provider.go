package domain

import (
	"github.com/gsouza97/my-bots/internal/adapter/provider"
)

type HomologacionProvider interface {
	GetHomologacionStatus(params string) (provider.HomologacionResponse, error)
}
