package kafka

type Consumer interface {
	ReadAllFromTopic(topic string) (<-chan []byte, error)
	ReadNewFromTopic(topic string) (<-chan []byte, error)
	ReadTopicFromOffset(topic string, offset int64) (<-chan []byte, error)
}
