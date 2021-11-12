package commands

type IntegrationEventStatus int

const (
	ToBeDispatched IntegrationEventStatus = iota
	Dispatched
)

type IntegrationEvent struct {
	UUID        string
	AggregateID string
	Name        string
	Payload     interface{}
	Version     string
	Status      IntegrationEventStatus
}

type IntegrationEventOutboxRepository interface {
	Save(event IntegrationEvent, transaction Transaction) error
}

type IntegrationEventBus interface {
	DispatchEvent(event IntegrationEvent)
}
