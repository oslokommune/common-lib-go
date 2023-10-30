package queue

type ActiveMQApi interface {
	Publish(topic, message string) error
	Disconnect()
}
