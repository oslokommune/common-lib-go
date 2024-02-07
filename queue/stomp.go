package queue

import (
	"crypto/tls"
	"time"

	stomp "github.com/go-stomp/stomp/v3"
	"github.com/go-stomp/stomp/v3/frame"
)

type StompClient struct {
	conn *stomp.Conn
}

func NewStompClient(broker, username, password string) (*StompClient, error) {
	options := []func(*stomp.Conn) error{
		stomp.ConnOpt.Login(username, password),
		stomp.ConnOpt.HeartBeat(3*time.Second, 3*time.Second),
		stomp.ConnOpt.UseStomp,
	}

	tlsConfig := &tls.Config{
		InsecureSkipVerify: false, // Set to false for production use with a valid certificate
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

func (s *StompClient) Disconnect() {
	s.conn.Disconnect()
}
