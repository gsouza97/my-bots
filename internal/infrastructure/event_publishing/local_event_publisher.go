package event_publishing

import (
	"context"
	"fmt"
	"sync"

	"github.com/gsouza97/my-bots/internal/domain"
	"github.com/gsouza97/my-bots/internal/domain/events"
	"github.com/gsouza97/my-bots/internal/logger"
)

// LocalEventPublisher implementa EventPublisher em memória
type LocalEventPublisher struct {
	handlers map[string][]domain.EventHandler
	mu       sync.RWMutex
}

// NewLocalEventPublisher cria um novo publisher local
func NewLocalEventPublisher() *LocalEventPublisher {
	return &LocalEventPublisher{
		handlers: make(map[string][]domain.EventHandler),
	}
}

// Publish executa todos os handlers registrados para um evento
func (lep *LocalEventPublisher) Publish(ctx context.Context, event events.Event) error {
	lep.mu.RLock()
	handlers, exists := lep.handlers[event.EventType()]
	lep.mu.RUnlock()

	if !exists {
		logger.Log.Debugf("No handlers registered for event type: %s", event.EventType())
		return nil
	}

	logger.Log.Debugf("Publishing event: %s (aggregate: %s)", event.EventType(), event.AggregateID())

	// Executar cada handler de forma assíncrona
	for _, handler := range handlers {
		go func(h domain.EventHandler) {
			if err := h(ctx, event); err != nil {
				logger.Log.Errorf(
					"Error handling event %s (aggregate: %s): %v",
					event.EventType(),
					event.AggregateID(),
					err,
				)
			}
		}(handler)
	}

	return nil
}

// Subscribe registra um handler para um tipo de evento
func (lep *LocalEventPublisher) Subscribe(eventType string, handler domain.EventHandler) error {
	if eventType == "" {
		return fmt.Errorf("event type cannot be empty")
	}

	if handler == nil {
		return fmt.Errorf("handler cannot be nil")
	}

	lep.mu.Lock()
	defer lep.mu.Unlock()

	lep.handlers[eventType] = append(lep.handlers[eventType], handler)
	logger.Log.Infof("Handler subscribed to event: %s (total handlers: %d)", eventType, len(lep.handlers[eventType]))

	return nil
}

func (lep *LocalEventPublisher) Close(ctx context.Context) error {
	logger.Log.Info("LocalEventPublisher closed")
	return nil
}

func (lep *LocalEventPublisher) GetHandlerCount(eventType string) int {
	lep.mu.RLock()
	defer lep.mu.RUnlock()
	return len(lep.handlers[eventType])
}
