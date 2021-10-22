package kafka

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/Shopify/sarama"
)

func NewClient(broker string) (*Client, error) {
	dataCollector, err := newDataCollector([]string{broker})
	if err != nil {
		return nil, err
	}

	return &Client{SyncProducer: dataCollector}, nil
}

type Client struct {
	SyncProducer sarama.SyncProducer
}

func (c Client) SendJSONSync(message interface{}, topic string) error {
	marshalledMessage, err := json.Marshal(message)
	if err != nil {
		return err
	}

	_, _, err = c.SyncProducer.SendMessage(
		&sarama.ProducerMessage{
			Topic: topic,
			Value: sarama.StringEncoder(marshalledMessage),
		},
	)

	return err
}

func newDataCollector(brokerList []string) (sarama.SyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll // Wait for all in-sync replicas to ack the message
	config.Producer.Retry.Max = 10                   // Retry up to 10 times to produce the message
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(brokerList, config)
	if err != nil {
		log.Fatalf("Failed to start Sarama producer on broker %s: %s\n", strings.Join(brokerList, ","), err)
	}

	return producer, err
}
