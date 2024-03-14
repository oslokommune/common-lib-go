package queue

type ActiveMQApi interface {
	Publish(topic, message string, header ...Header) error
	Disconnect()
}

type Header struct {
	Key   string
	Value string
}
