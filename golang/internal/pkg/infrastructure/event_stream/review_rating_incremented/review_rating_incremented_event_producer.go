package review_rating_incremented

import (
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/ProntoPro/event-stream-golang/internal/pkg/application/commands"
	"github.com/ProntoPro/event-stream-golang/internal/pkg/infrastructure/event_stream"
)

func NewReviewRatingIncrementedEventProducer(
	producer event_stream.Producer,
	reviewCreatedEventMarshaller ReviewRatingIncrementedEventMarshaller,
) *ReviewRatingIncrementedEventProducer {
	return &ReviewRatingIncrementedEventProducer{
		producer:                     producer,
		reviewCreatedEventMarshaller: reviewCreatedEventMarshaller,
	}
}

type ReviewRatingIncrementedEventProducer struct {
	producer                     event_stream.Producer
	reviewCreatedEventMarshaller ReviewRatingIncrementedEventMarshaller
}

func (r *ReviewRatingIncrementedEventProducer) DispatchEvent(event commands.IntegrationEvent) {
	eventPayload, ok := event.Payload.(commands.ReviewRatingIncrementedEvent)
	if !ok {
		r.handleErrors(fmt.Errorf("unsupported event payload"))
	}

	messageData, err := r.reviewCreatedEventMarshaller.Marshal(
		&ReviewRatingIncrementedEventMessage{
			Review: ReviewMessage{
				UUID: eventPayload.ReviewUUID,
			},
		},
	)

	if err != nil {
		r.handleErrors(err)
	}

	// todo what happens if we fail before here...
	err = r.producer.Dispatch(
		messageData,
	)
	// todo what happens if we fail here...

	if err != nil {
		r.handleErrors(err)
	}
}

func (r ReviewRatingIncrementedEventProducer) handleErrors(err error) {
	// todo handle error so to avoid losing events
	logrus.Errorf("error sending review created message via event stream: %s\n", err.Error())
}
