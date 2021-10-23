package cmd

import (
	"fmt"

	"github.com/urfave/cli"
)

func consumeReviewCreatedEvent(baseFlags []cli.Flag) cli.Command {
	return cli.Command{
		Name:   "consume-review-created",
		Action: consumeReviewCreated,
		Flags:  baseFlags,
	}
}

func consumeReviewCreated(c *cli.Context) error {
	serviceLocator, err := NewServiceLocator(c.String("kafka-url"))
	if err != nil {
		return err
	}

	useJSON := true
	switch c.String("messaging-protocol") {
	case "protobuf":
		useJSON = false
	}

	consumer := serviceLocator.ReviewCreatedConsumerWithKafkaAndJSON()
	if !useJSON {
		consumer = serviceLocator.ReviewCreatedConsumerWithKafkaAndProtobuf()
	}

	err = consumer.Consume()
	if err != nil {
		fmt.Println(err)
	}

	return nil
}
