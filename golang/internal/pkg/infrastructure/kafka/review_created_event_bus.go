package kafka

import (
	"github.com/ProntoPro/golang-kafka/internal/pkg/application"
)

type ReviewCreatedEventBus struct {
}

func (r *ReviewCreatedEventBus) DispatchEvent(event application.ReviewCreatedEvent) {
	// todo implement
}
