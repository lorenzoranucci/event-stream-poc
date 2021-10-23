package kafka

import (
	"encoding/json"

	"github.com/sirupsen/logrus"
)

func NewReviewCreatedEventConsumer(client Consumer) *ReviewCreatedEventConsumer {
	return &ReviewCreatedEventConsumer{kafkaClient: client}
}

type ReviewCreatedEventConsumer struct {
	kafkaClient Consumer
}

func (r *ReviewCreatedEventConsumer) Consume() error {
	messages, err := r.kafkaClient.ReadAllFromTopic("review_created_event")

	if err != nil {
		logrus.Error(err)
		return err
	}

	for message := range messages {
		logrus.Infof("processing message %s", string(message))
		reviewCreatedEvent := &ReviewCreatedEventMessage{}
		err := json.Unmarshal(message, reviewCreatedEvent)
		if err != nil {
			return err
		}
		logrus.Infof("processing review created event %#v", reviewCreatedEvent)
	}

	return nil
}
