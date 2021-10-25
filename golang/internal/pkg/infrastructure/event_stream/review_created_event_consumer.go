package event_stream

import (
	"github.com/sirupsen/logrus"
)

func NewReviewCreatedEventConsumer(
	consumer Consumer,
	reviewCreatedEventMarshaller ReviewCreatedEventMarshaller,
) *ReviewCreatedEventConsumer {
	return &ReviewCreatedEventConsumer{
		consumer:                     consumer,
		reviewCreatedEventMarshaller: reviewCreatedEventMarshaller,
	}
}

type ReviewCreatedEventConsumer struct {
	consumer                     Consumer
	reviewCreatedEventMarshaller ReviewCreatedEventMarshaller
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
	}

	return nil
}
