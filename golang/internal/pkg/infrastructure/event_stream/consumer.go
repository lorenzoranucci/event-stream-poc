package event_stream

type Consumer interface {
	ConsumeAll() (<-chan []byte, error)
}
