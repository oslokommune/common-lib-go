package queue

import (
	MQTT "github.com/eclipse/paho.mqtt.golang"
)

type MQTTClient struct {
	client MQTT.Client
}

func NewMQTTClient(broker, username, password string) (*MQTTClient, error) {
	opts := MQTT.NewClientOptions()
	opts.AddBroker(broker)
	opts.SetClientID("id")
	opts.SetUsername(username)
	opts.SetPassword(password)
	opts.SetCleanSession(false)

	client := MQTT.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}

	return &MQTTClient{
		client,
	}, nil
}

func (m *MQTTClient) Publish(topic, message string) error {
	m.client.Publish(topic, byte(0), false, message)
	return nil
}

func (m *MQTTClient) Disconnect() {
	m.client.Disconnect(255)
}
