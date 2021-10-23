package kafka

type Producer interface {
	Dispatch(message []byte) error
}
