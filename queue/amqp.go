package queue

import (
	"context"

	amqp "github.com/Azure/go-amqp"
)

type AMQPClient struct {
	conn *amqp.Conn
}

func (a *AMQPClient) Disconnect() {
	a.conn.Close()
}

func (a *AMQPClient) Publish(topic string, message string) error {
	// open a session
	session, err := a.conn.NewSession(context.Background(), nil)
	if err != nil {
		return err
	}

	// create a sender
	sender, err := session.NewSender(context.Background(), topic, nil)
	if err != nil {
		return err
	}

	return sender.Send(context.Background(), amqp.NewMessage([]byte(message)), nil)
}

func NewAMQPClient(broker, username, password string) (*AMQPClient, error) {
	conn, err := amqp.Dial(context.Background(), broker, &amqp.ConnOptions{
		SASLType: amqp.SASLTypePlain(username, password),
	})
	if err != nil {
		return nil, err
	}

	return &AMQPClient{conn}, nil
}
