package application

type IntegrationEvent struct {
	UUID        string
	AggregateID string
	Name        string
	Payload     interface{}
	Version     string
}

type IntegrationEventOutboxRepository interface {
	Add(event IntegrationEvent, transaction Transaction) error
}

type IntegrationEventBus interface {
	DispatchEvent(event IntegrationEvent)
}
