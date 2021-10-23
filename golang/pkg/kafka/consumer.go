package kafka

import (
	"github.com/Shopify/sarama"
)

func NewConsumer(broker string, readPartition int32, topic string) (*Consumer, error) {
	consumer, err := sarama.NewConsumer([]string{broker}, sarama.NewConfig())
	if err != nil {
		return nil, err
	}
	return &Consumer{consumer: consumer, readPartition: readPartition, topic: topic}, nil
}

type Consumer struct {
	consumer      sarama.Consumer
	readPartition int32
	topic         string
}

func (c *Consumer) ConsumeAll() (<-chan []byte, error) {
	return c.ConsumeFromOffset(sarama.OffsetOldest)
}

func (c *Consumer) ConsumeNew() (<-chan []byte, error) {
	return c.ConsumeFromOffset(sarama.OffsetNewest)
}

func (c *Consumer) ConsumeFromOffset(offset int64) (<-chan []byte, error) {
	partitionConsumer, err := c.consumer.ConsumePartition(c.topic, c.readPartition, offset)
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
