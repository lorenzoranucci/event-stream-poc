package pulsar

import (
	"context"
	"time"

	"github.com/apache/pulsar-client-go/pulsar"
)

func NewProducer(broker string, topic string) (*Producer, error) {
	client, err := pulsar.NewClient(pulsar.ClientOptions{
		URL:               broker,
		OperationTimeout:  10 * time.Second,
		ConnectionTimeout: 10 * time.Second,
	})
	if err != nil {
		return nil, err
	}

	producer, err := client.CreateProducer(pulsar.ProducerOptions{
		Topic: topic,
	})
	if err != nil {
		return nil, err
	}

	return &Producer{producer: producer}, nil
}

type Producer struct {
	producer pulsar.Producer
}

func (p *Producer) Dispatch(message []byte) error {
	_, err := p.producer.Send(
		context.Background(),
		&pulsar.ProducerMessage{
			Payload: message,
		},
	)

	return err
}

func (p *Producer) Close() error {
	p.producer.Close()

	return nil
}
