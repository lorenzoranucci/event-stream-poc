package kafka

import (
	"github.com/ProntoPro/event-stream-golang/internal/pkg/application"
)

type ReviewCreatedEventBus struct {
}

func (r *ReviewCreatedEventBus) DispatchEvent(event application.ReviewCreatedEvent) {
	// todo implement
}
