package kafka

type Producer interface {
	SendJSONSync(message interface{}, topic string) error
}
