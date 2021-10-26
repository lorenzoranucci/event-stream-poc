package event_stream

import (
	"github.com/sirupsen/logrus"

	"github.com/ProntoPro/event-stream-golang/internal/pkg/application"
)

func NewReviewCreatedEventProducer(
	producer Producer,
	reviewCreatedEventMarshaller ReviewCreatedEventMarshaller,
) *ReviewCreatedEventProducer {
	return &ReviewCreatedEventProducer{
		producer:                     producer,
		reviewCreatedEventMarshaller: reviewCreatedEventMarshaller,
	}
}

type ReviewCreatedEventProducer struct {
	producer                     Producer
	reviewCreatedEventMarshaller ReviewCreatedEventMarshaller
}

func (r *ReviewCreatedEventProducer) DispatchEvent(event application.ReviewCreatedEvent) {
	messageData, err := r.reviewCreatedEventMarshaller.Marshal(
		&ReviewCreatedEventMessage{
			Review: ReviewMessage{
				UUID:    event.UUID,
				Comment: event.Comment,
				Rating:  event.Rating,
			},
		},
	)

	if err != nil {
		r.handleErrors(err)
	}

	err = r.producer.Dispatch(
		messageData,
	)

	if err != nil {
		r.handleErrors(err)
	}
}

func (r ReviewCreatedEventProducer) handleErrors(err error) {
	// todo handle error so to avoid losing events
	logrus.Errorf("error sending review created message via event stream: %s\n", err.Error())
}
