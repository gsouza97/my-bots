package domain

import "github.com/gsouza97/my-bots/internal/infrastructure/providers"

type HomologacionProvider interface {
	GetHomologacionStatus(params string) (providers.HomologacionResponse, error)
}
