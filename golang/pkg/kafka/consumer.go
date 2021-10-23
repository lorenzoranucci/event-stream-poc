package kafka

import (
	"github.com/Shopify/sarama"
)

func NewConsumer(broker string, readPartition int32) (*Consumer, error) {
	consumer, err := sarama.NewConsumer([]string{broker}, sarama.NewConfig())
	if err != nil {
		return nil, err
	}
	return &Consumer{consumer: consumer, readPartition: readPartition}, nil
}

type Consumer struct {
	consumer      sarama.Consumer
	readPartition int32
}

func (c *Consumer) ReadAllFromTopic(topic string) (<-chan []byte, error) {
	return c.ReadTopicFromOffset(topic, sarama.OffsetOldest)
}

func (c *Consumer) ReadNewFromTopic(topic string) (<-chan []byte, error) {
	return c.ReadTopicFromOffset(topic, sarama.OffsetNewest)
}

func (c *Consumer) ReadTopicFromOffset(topic string, offset int64) (<-chan []byte, error) {
	partitionConsumer, err := c.consumer.ConsumePartition(topic, c.readPartition, offset)
	if err != nil {
		return nil, err
	}

	var returnChannel = make(chan []byte)
	go func(channel <-chan *sarama.ConsumerMessage) {
		for message := range channel {
			returnChannel <- message.Value
		}
	}(partitionConsumer.Messages())

	return returnChannel, nil
}

func (c *Consumer) Close() error {
	return c.consumer.Close()
}
