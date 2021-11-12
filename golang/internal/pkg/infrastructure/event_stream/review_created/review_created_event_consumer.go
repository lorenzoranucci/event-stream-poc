package review_created

import (
	"github.com/sirupsen/logrus"

	"github.com/ProntoPro/event-stream-golang/internal/pkg/application/commands"
	"github.com/ProntoPro/event-stream-golang/internal/pkg/infrastructure/event_stream"
	"github.com/ProntoPro/event-stream-golang/internal/pkg/infrastructure/projectors/get_reviews"
)

func NewReviewCreatedEventConsumer(
	consumer event_stream.Consumer,
	reviewCreatedEventMarshaller ReviewCreatedEventMarshaller,
	projector *get_reviews.ReviewCreatedProjector,
) *ReviewCreatedEventConsumer {
	return &ReviewCreatedEventConsumer{
		consumer:                     consumer,
		reviewCreatedEventMarshaller: reviewCreatedEventMarshaller,
		projector:                    projector,
	}
}

type ReviewCreatedEventConsumer struct {
	consumer                     event_stream.Consumer
	reviewCreatedEventMarshaller ReviewCreatedEventMarshaller
	projector                    *get_reviews.ReviewCreatedProjector
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
			&commands.ReviewCreatedEvent{
				ReviewUUID: eventMessage.Review.UUID,
				Comment:    eventMessage.Review.Comment,
				Rating:     eventMessage.Review.Rating,
			},
		)
		if err != nil {
			// todo handle error
			logrus.Error(err)
		}
	}

	return nil
}
