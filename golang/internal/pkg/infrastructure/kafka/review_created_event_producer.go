package kafka

import (
	"github.com/sirupsen/logrus"

	"github.com/ProntoPro/event-stream-golang/internal/pkg/application"
)

func NewReviewCreatedEventBus(client Producer) *ReviewCreatedEventBus {
	return &ReviewCreatedEventBus{kafkaProducer: client}
}

type ReviewCreatedEventBus struct {
	kafkaProducer Producer
}

func (r *ReviewCreatedEventBus) DispatchEvent(event application.ReviewCreatedEvent) {
	err := r.kafkaProducer.SendJSONSync(ReviewCreatedEventMessage{
		Review: ReviewMessage{
			UUID:    event.Review.Uuid().String(),
			Comment: event.Review.Comment(),
			Rating:  event.Review.Rating(),
		},
	}, "review_created_event")

	if err != nil {
		// todo handle error
		logrus.Errorf("error sending review created message via Kafka: %s\n", err.Error())
	}
}
