package pulsar

import (
	"time"

	"github.com/apache/pulsar-client-go/pulsar"
)

func NewConsumer(broker string, topic string) (*Consumer, error) {
	client, err := pulsar.NewClient(pulsar.ClientOptions{
		URL:               broker,
		OperationTimeout:  10 * time.Second,
		ConnectionTimeout: 10 * time.Second,
	})
	if err != nil {
		return nil, err
	}

	consumer, err := client.Subscribe(pulsar.ConsumerOptions{
		Topic:            topic,
		SubscriptionName: topic,
	})
	if err != nil {
		return nil, err
	}

	return &Consumer{consumer: consumer}, nil
}

type Consumer struct {
	consumer pulsar.Consumer
}

func (c *Consumer) ConsumeAll() (<-chan []byte, error) {
	channel := c.consumer.Chan()

	var returnChannel = make(chan []byte)
	go func(channel <-chan pulsar.ConsumerMessage) {
		for message := range channel {
			returnChannel <- message.Payload()
		}
	}(channel)

	return returnChannel, nil
}

func (c *Consumer) Close() error {
	c.consumer.Close()

	return nil
}
