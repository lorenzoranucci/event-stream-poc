package review_created

import (
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/ProntoPro/event-stream-golang/internal/pkg/application/commands"
	"github.com/ProntoPro/event-stream-golang/internal/pkg/infrastructure/event_stream"
)

func NewReviewCreatedEventProducer(
	producer event_stream.Producer,
	reviewCreatedEventMarshaller ReviewCreatedEventMarshaller,
) *ReviewCreatedEventProducer {
	return &ReviewCreatedEventProducer{
		producer:                     producer,
		reviewCreatedEventMarshaller: reviewCreatedEventMarshaller,
	}
}

type ReviewCreatedEventProducer struct {
	producer                     event_stream.Producer
	reviewCreatedEventMarshaller ReviewCreatedEventMarshaller
}

func (r *ReviewCreatedEventProducer) DispatchEvent(event commands.IntegrationEvent) {
	eventPayload, ok := event.Payload.(commands.ReviewCreatedEvent)
	if !ok {
		r.handleErrors(fmt.Errorf("unsupported event payload"))
	}

	messageData, err := r.reviewCreatedEventMarshaller.Marshal(
		&ReviewCreatedEventMessage{
			Review: ReviewMessage{
				UUID:    eventPayload.ReviewUUID,
				Comment: eventPayload.Comment,
				Rating:  eventPayload.Rating,
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

func (r ReviewCreatedEventProducer) handleErrors(err error) {
	// todo handle error so to avoid losing events
	logrus.Errorf("error sending review created message via event stream: %s\n", err.Error())
}
