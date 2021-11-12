package review_rating_incremented

import (
	"github.com/sirupsen/logrus"

	"github.com/ProntoPro/event-stream-golang/internal/pkg/application/commands"
	"github.com/ProntoPro/event-stream-golang/internal/pkg/infrastructure/event_stream"
	"github.com/ProntoPro/event-stream-golang/internal/pkg/infrastructure/projectors/get_reviews"
)

func NewReviewRatingIncrementedEventConsumer(
	consumer event_stream.Consumer,
	reviewCreatedEventMarshaller ReviewRatingIncrementedEventMarshaller,
	projector *get_reviews.ReviewRatingIncrementedProjector,
) *ReviewRatingIncrementedEventConsumer {
	return &ReviewRatingIncrementedEventConsumer{
		consumer:                     consumer,
		reviewCreatedEventMarshaller: reviewCreatedEventMarshaller,
		projector:                    projector,
	}
}

type ReviewRatingIncrementedEventConsumer struct {
	consumer                     event_stream.Consumer
	reviewCreatedEventMarshaller ReviewRatingIncrementedEventMarshaller
	projector                    *get_reviews.ReviewRatingIncrementedProjector
}

func (r *ReviewRatingIncrementedEventConsumer) Consume() error {
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
			&commands.ReviewRatingIncrementedEvent{
				ReviewUUID: eventMessage.Review.UUID,
			},
		)
		if err != nil {
			// todo handle error
			logrus.Error(err)
		}
	}

	return nil
}
