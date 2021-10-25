package event_stream

type Producer interface {
	Dispatch(message []byte) error
}
