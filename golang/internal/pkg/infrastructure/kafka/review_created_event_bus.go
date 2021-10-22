package kafka

import (
	"fmt"

	"github.com/ProntoPro/event-stream-golang/internal/pkg/application"
	"github.com/ProntoPro/event-stream-golang/pkg/kafka"
)

func NewReviewCreatedEventBus(client *kafka.Client) *ReviewCreatedEventBus {
	return &ReviewCreatedEventBus{kafkaClient: client}
}

type ReviewCreatedEventBus struct {
	kafkaClient *kafka.Client
}

type ReviewCreatedEventMessage struct {
	Review ReviewMessage `json:"review"`
}

type ReviewMessage struct {
	UUID    string `json:"uuid"`
	Comment string `json:"comment"`
	Rating  int    `json:"rating"`
}

func (r *ReviewCreatedEventBus) DispatchEvent(event application.ReviewCreatedEvent) {
	err := r.kafkaClient.SendJSONSync(ReviewCreatedEventMessage{
		Review: ReviewMessage{
			UUID:    event.Review.Uuid().String(),
			Comment: event.Review.Comment(),
			Rating:  event.Review.Rating(),
		},
	}, "review_created_event")

	if err != nil {
		// todo handle error
		fmt.Printf("error sending review created message via Kafka: %s\n", err.Error())
	}
}
