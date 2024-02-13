package queue

import (
	"crypto/tls"
	"errors"
	"time"

	stomp "github.com/go-stomp/stomp/v3"
	"github.com/go-stomp/stomp/v3/frame"
)

type StompClient struct {
	conn                       *stomp.Conn
	broker, username, password string
}

func NewStompClient(broker, username, password string) *StompClient {
	return &StompClient{username: username, broker: broker, password: password}
}

// Connect connets client to broker
func (s *StompClient) Connect() error {
	options := []func(*stomp.Conn) error{
		stomp.ConnOpt.Login(s.username, s.password),
		stomp.ConnOpt.HeartBeat(3*time.Second, 3*time.Second),
		stomp.ConnOpt.UseStomp,
	}

	tlsConfig := &tls.Config{
		InsecureSkipVerify: false, // Set to false for production use with a valid certificate
	}

	// Establish a TLS connection
	tlsConnection, err := tls.Dial("tcp", s.broker, tlsConfig)
	if err != nil {
		return err
	}

	conn, err := stomp.Connect(tlsConnection, options...)
	if err != nil {
		return err
	}

	s.conn = conn
	return nil
}

// Publish pulishes message to queue. Will reconnect if connection is already closed
func (s *StompClient) Publish(destination string, msg string) error {
	err := s.send(destination, msg)
	if err != nil {
		switch {
		case errors.As(err, &stomp.ErrAlreadyClosed):
			{
				connErr := s.Connect()
				if connErr != nil {
					return connErr
				}
				return s.send(destination, msg)
			}
		default:
			return err
		}
	}
	return nil
}

func (s *StompClient) send(destination string, msg string) error {
	contentType := "application/json;charset=utf-8"
	if len(msg) > 0 && msg[0] == '<' {
		contentType = "application/xml;charset=utf-8"
	}

	return s.conn.Send(
		destination, // destination
		contentType,
		[]byte(msg), // body
		stomp.SendOpt.Receipt,
		func(f *frame.Frame) error {
			f.Header.Del(frame.ContentLength)
			f.Header.Add("persistent", "true")
			f.Header.Add("Destination", destination)
			return nil
		})
}

// Disconnect closes connection to broker
func (s *StompClient) Disconnect() {
	s.conn.Disconnect()
}
