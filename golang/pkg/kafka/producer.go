package kafka

import (
	"encoding/json"

	"github.com/Shopify/sarama"
)

func NewProducer(broker string) (*Producer, error) {
	syncProducer, err := newSyncProducer([]string{broker})
	if err != nil {
		return nil, err
	}

	return &Producer{syncProducer: syncProducer}, nil
}

type Producer struct {
	syncProducer sarama.SyncProducer
}

func (p *Producer) SendJSONSync(message interface{}, topic string) error {
	marshalledMessage, err := json.Marshal(message)
	if err != nil {
		return err
	}

	_, _, err = p.syncProducer.SendMessage(
		&sarama.ProducerMessage{
			Topic: topic,
			Value: sarama.StringEncoder(marshalledMessage),
		},
	)

	return err
}

func newSyncProducer(brokerList []string) (sarama.SyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll // Wait for all in-sync replicas to ack the message
	config.Producer.Retry.Max = 10                   // Retry up to 10 times to produce the message
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(brokerList, config)
	if err != nil {
		return nil, err
	}

	return producer, err
}

func (p *Producer) Close() error {
	return p.syncProducer.Close()
}
