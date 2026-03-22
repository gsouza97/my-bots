package domain

import (
	"context"

	"github.com/gsouza97/my-bots/internal/domain/events"
)

// EventHandler é uma função que processa um evento
// Qualquer domínio que queira escutar eventos implementa essa assinatura
type EventHandler func(ctx context.Context, event events.Event) error

// EventPublisher define a interface para publicação de eventos
// Desacopla use cases de seus subscribers (listeners)
// Permite trocar de implementação (local, kafkan redis, etc) sem mexer em use cases
type EventPublisher interface {
	// Publish publica um evento para todos os handlers registrados
	Publish(ctx context.Context, event events.Event) error

	// Subscribe registra um handler para um tipo de evento específico
	Subscribe(eventType string, handler EventHandler) error

	Close(ctx context.Context) error
}
