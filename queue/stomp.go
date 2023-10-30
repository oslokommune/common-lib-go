package queue

import (
	"crypto/tls"

	stomp "github.com/go-stomp/stomp/v3"
	"github.com/go-stomp/stomp/v3/frame"
)

type StompClient struct {
	conn *stomp.Conn
}

func NewStompClient(broker, username, password string) (*StompClient, error) {
	options := []func(*stomp.Conn) error{
		stomp.ConnOpt.Login(username, password),
	}

	tlsConfig := &tls.Config{
		InsecureSkipVerify: true, // Set to false for production use with a valid certificate
	}

	// Establish a TLS connection
	tlsConnection, err := tls.Dial("tcp", broker, tlsConfig)
	if err != nil {
		return nil, err
	}

	conn, err := stomp.Connect(tlsConnection, options...)
	if err != nil {
		return nil, err
	}

	return &StompClient{conn}, nil
}

func (s *StompClient) Publish(destination string, msg string) error {
	return s.conn.Send(
		destination,                     // destination
		"application/xml;charset=utf-8", // content-type
		[]byte(msg),                     // body
		func(f *frame.Frame) error {
			f.Header.Del(frame.ContentLength)
			return nil
		})
}

func (s *StompClient) Disconnect() {
	s.conn.Disconnect()
}
