package kafka

type Consumer interface {
	ConsumeAll() (<-chan []byte, error)
	ConsumeNew() (<-chan []byte, error)
	ConsumeFromOffset(offset int64) (<-chan []byte, error)
}
