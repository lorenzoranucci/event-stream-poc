package kafka

import (
	"github.com/sirupsen/logrus"
)

func NewReviewCreatedEventConsumer(
	kafkaClient Consumer,
	reviewCreatedEventMarshaller ReviewCreatedEventMarshaller,
) *ReviewCreatedEventConsumer {
	return &ReviewCreatedEventConsumer{
		kafkaClient:                  kafkaClient,
		reviewCreatedEventMarshaller: reviewCreatedEventMarshaller,
	}
}

type ReviewCreatedEventConsumer struct {
	kafkaClient                  Consumer
	reviewCreatedEventMarshaller ReviewCreatedEventMarshaller
}

func (r *ReviewCreatedEventConsumer) Consume() error {
	messages, err := r.kafkaClient.ConsumeAll()

	if err != nil {
		logrus.Error(err)
		return err
	}

	for message := range messages {
		logrus.Info("processing message")
		eventMessage, err := r.reviewCreatedEventMarshaller.Unmarshal(message)
		if err != nil {
			return err
		}

		logrus.Infof("processing review created event %#v", eventMessage)
	}

	return nil
}
