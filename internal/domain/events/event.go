package events

import "time"

type Event interface {
	EventType() string
	AggregateID() string
	OccurredAt() time.Time
}

type BaseEvent struct {
	eventType   string
	aggregateID string
	occurredAt  time.Time
}

func (be *BaseEvent) EventType() string {
	return be.eventType
}

func (be *BaseEvent) AggregateID() string {
	return be.aggregateID
}

func (be *BaseEvent) OccurredAt() time.Time {
	return be.occurredAt
}

func NewBaseEvent(eventType, aggregateID string) BaseEvent {
	return BaseEvent{
		eventType:   eventType,
		aggregateID: aggregateID,
		occurredAt:  time.Now(),
	}
}

// PriceAlertTriggeredEvent é emitido quando um alerta de preço é acionado
type PriceAlertTriggeredEvent struct {
	BaseEvent
	AlertID      string
	Crypto       string
	CurrentPrice float64
	AlertPrice   float64
	Status       string
}

func NewPriceAlertTriggeredEvent(alertID string, crypto string, currentPrice float64, alertPrice float64, status string) *PriceAlertTriggeredEvent {
	return &PriceAlertTriggeredEvent{
		BaseEvent:    NewBaseEvent("PriceAlertTriggered", alertID),
		AlertID:      alertID,
		Crypto:       crypto,
		CurrentPrice: currentPrice,
		AlertPrice:   alertPrice,
		Status:       status,
	}
}

// DailyAlertTriggeredEvent é emitido diariamente com resumo
type DailyAlertTriggeredEvent struct {
	BaseEvent
	Message string
}

func NewDailyAlertTriggeredEvent(message string) *DailyAlertTriggeredEvent {
	return &DailyAlertTriggeredEvent{
		BaseEvent: NewBaseEvent("DailyAlertTriggered", "system"),
		Message:   message,
	}
}

// PoolAlertTriggeredEvent é emitido quando uma pool sai do range
type PoolAlertTriggeredEvent struct {
	BaseEvent
	PoolID      string
	Description string
	Message     string
}

func NewPoolAlertTriggeredEvent(poolID string, description string, message string) *PoolAlertTriggeredEvent {
	return &PoolAlertTriggeredEvent{
		BaseEvent:   NewBaseEvent("PoolAlertTriggered", poolID),
		PoolID:      poolID,
		Description: description,
		Message:     message,
	}
}

// LoanAlertTriggeredEvent é emitido quando um empréstimo precisa atenção
type LoanAlertTriggeredEvent struct {
	BaseEvent
	LoanID  string
	Message string
}

func NewLoanAlertTriggeredEvent(loanID string, message string) *LoanAlertTriggeredEvent {
	return &LoanAlertTriggeredEvent{
		BaseEvent: NewBaseEvent("LoanAlertTriggered", loanID),
		LoanID:    loanID,
		Message:   message,
	}
}

// HomologacionAlertTriggeredEvent é emitido para alertas de homologação
type HomologacionAlertTriggeredEvent struct {
	BaseEvent
	HomologacionID string
	Message        string
}

func NewHomologacionAlertTriggeredEvent(homologacionID string, message string) *HomologacionAlertTriggeredEvent {
	return &HomologacionAlertTriggeredEvent{
		BaseEvent:      NewBaseEvent("HomologacionAlertTriggered", "homologacion"),
		HomologacionID: homologacionID,
		Message:        message,
	}
}
