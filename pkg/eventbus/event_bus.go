package eventbus

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"backend-sembako/pkg/logger"

	"github.com/google/uuid"
)

type DomainEvent struct {
	ID          uuid.UUID              `json:"id"`
	Name        string                 `json:"name"`
	Aggregate   string                 `json:"aggregate"`
	AggregateID uuid.UUID              `json:"aggregate_id"`
	Payload     map[string]interface{} `json:"payload"`
	OccurredAt  time.Time              `json:"occurred_at"`
}

type EventHandler func(ctx context.Context, event DomainEvent) error

type EventBus interface {
	Subscribe(eventName string, handler EventHandler)
	Publish(ctx context.Context, event DomainEvent) error
}

type inMemoryEventBus struct {
	handlers map[string][]EventHandler
	mu       sync.RWMutex
}

func NewEventBus() EventBus {
	return &inMemoryEventBus{
		handlers: make(map[string][]EventHandler),
	}
}

func (b *inMemoryEventBus) Subscribe(eventName string, handler EventHandler) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.handlers[eventName] = append(b.handlers[eventName], handler)
}

func (b *inMemoryEventBus) Publish(ctx context.Context, event DomainEvent) error {
	b.mu.RLock()
	handlers, exists := b.handlers[event.Name]
	b.mu.RUnlock()

	if !exists || len(handlers) == 0 {
		return nil
	}

	// Dispatch handlers asynchronously so calling service is never blocked!
	go func(evt DomainEvent, hList []EventHandler) {
		asyncCtx := context.Background()
		for _, h := range hList {
			if err := h(asyncCtx, evt); err != nil {
				logger.Log.Error().
					Err(err).
					Str("event_name", evt.Name).
					Str("aggregate", evt.Aggregate).
					Msg("Failed to execute event subscriber handler")
			}
		}
	}(event, handlers)

	return nil
}

func (e *DomainEvent) ToJSON() string {
	bytes, _ := json.Marshal(e.Payload)
	return string(bytes)
}
