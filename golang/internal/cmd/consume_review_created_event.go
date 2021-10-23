package cmd

import (
	"fmt"

	"github.com/urfave/cli"

	"github.com/ProntoPro/event-stream-golang/internal/pkg/infrastructure/kafka"
	kafka2 "github.com/ProntoPro/event-stream-golang/pkg/kafka"
)

func consumeReviewCreatedEvent(baseFlags []cli.Flag) cli.Command {
	return cli.Command{
		Name:   "consume-review-created",
		Action: consumeReviewCreated,
		Flags:  baseFlags,
	}
}

func consumeReviewCreated(c *cli.Context) error {
	kafkaConsumer, err := kafka2.NewConsumer(c.String("kafka-url"), 0, "review_created_event")
	if err != nil {
		return err
	}

	consumer := kafka.NewReviewCreatedEventConsumer(kafkaConsumer, &kafka.ReviewCreatedEventJSONMarshaller{})

	err = consumer.Consume()
	if err != nil {
		fmt.Println(err)
	}

	return nil
}
