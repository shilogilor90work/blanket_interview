package natsHandler

import (
	"encoding/json"
	"time"

	"github.com/nats-io/nats.go"
)

var natsConn *nats.Conn

// InitNATS initializes the NATS connection
func InitNATS(natsURL string) error {
	var err error
	natsConn, err = nats.Connect(natsURL)
	if err != nil {
		return err
	}
	return nil
}

// Publish sends a message to a NATS subject
func Publish(subject string, msg interface{}) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return natsConn.Publish(subject, data)
}

func Request(subject string, requestMsg interface{}, timeout time.Duration) (*nats.Msg, error) {
	// Marshal the request message to JSON
	data, err := json.Marshal(requestMsg)
	if err != nil {
		return nil, err
	}

	// Send the request and wait for a response
	responseMsg, err := natsConn.Request(subject, data, timeout)
	if err != nil {
		return nil, err
	}

	return responseMsg, nil
}
