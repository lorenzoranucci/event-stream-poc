package event_stream

import (
	"github.com/sirupsen/logrus"

	"github.com/ProntoPro/event-stream-golang/internal/pkg/application"
	"github.com/ProntoPro/event-stream-golang/internal/pkg/infrastructure/projectors"
)

func NewReviewCreatedEventConsumer(
	consumer Consumer,
	reviewCreatedEventMarshaller ReviewCreatedEventMarshaller,
	projector *projectors.GetReviewsProjector,
) *ReviewCreatedEventConsumer {
	return &ReviewCreatedEventConsumer{
		consumer:                     consumer,
		reviewCreatedEventMarshaller: reviewCreatedEventMarshaller,
		projector:                    projector,
	}
}

type ReviewCreatedEventConsumer struct {
	consumer                     Consumer
	reviewCreatedEventMarshaller ReviewCreatedEventMarshaller
	projector                    *projectors.GetReviewsProjector
}

func (r *ReviewCreatedEventConsumer) Consume() error {
	messages, err := r.consumer.ConsumeAll()

	if err != nil {
		logrus.Error(err)
		return err
	}

	for message := range messages {
		logrus.Infof("processing message %s", string(message))
		eventMessage, err := r.reviewCreatedEventMarshaller.Unmarshal(message)
		if err != nil {
			// todo handle error
			logrus.Error(err)
		}

		logrus.Infof("processing review created event %#v", eventMessage)

		err = r.projector.Project(
			&application.ReviewCreatedEvent{
				UUID:    eventMessage.Review.UUID,
				Comment: eventMessage.Review.Comment,
				Rating:  eventMessage.Review.Rating,
			},
		)
		if err != nil {
			// todo handle error
			logrus.Error(err)
		}
	}

	return nil
}
