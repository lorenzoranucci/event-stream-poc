package kafka

import (
	"github.com/sirupsen/logrus"

	"github.com/ProntoPro/event-stream-golang/internal/pkg/application"
)

func NewReviewCreatedEventProducer(
	kafkaProducer Producer,
	reviewCreatedEventMarshaller ReviewCreatedEventMarshaller,
) *ReviewCreatedEventProducer {
	return &ReviewCreatedEventProducer{
		kafkaProducer:                kafkaProducer,
		reviewCreatedEventMarshaller: reviewCreatedEventMarshaller,
	}
}

type ReviewCreatedEventProducer struct {
	kafkaProducer                Producer
	reviewCreatedEventMarshaller ReviewCreatedEventMarshaller
}

func (r *ReviewCreatedEventProducer) DispatchEvent(event application.ReviewCreatedEvent) {
	messageData, err := r.reviewCreatedEventMarshaller.Marshal(
		&ReviewCreatedEventMessage{
			Review: ReviewMessage{
				UUID:    event.Review.Uuid().String(),
				Comment: event.Review.Comment(),
				Rating:  event.Review.Rating(),
			},
		},
	)

	if err != nil {
		r.handleErrors(err)
	}

	err = r.kafkaProducer.Dispatch(
		messageData,
	)

	if err != nil {
		r.handleErrors(err)
	}
}

func (r ReviewCreatedEventProducer) handleErrors(err error) {
	// todo handle error so to avoid losing events
	logrus.Errorf("error sending review created message via Kafka: %s\n", err.Error())
}
