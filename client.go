package mq2http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/project-voyager/utils/amqp"
	"github.com/project-voyager/utils/random"
	"github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	samqp "github.com/streadway/amqp"
)

// Transport is a ....
type Transport struct {
	AMQPConn        *samqp.Connection
	ReplyToQueue    string
	AMQPRPCSettings AMQPRPCSettings
	Timeout         time.Duration
}

// RoundTrip is a ...
func (t Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	rpc := t.AMQPRPCSettings
	// Listen
	_, deliveries, err := amqp.Listen(t.AMQPConn, rpc.Exchange, rpc.ExchangeType, t.ReplyToQueue, rpc.BindingKey, rpc.ConsumerTag)
	if err != nil {
		return nil, err
	}
	// Grab request and send over amqp using "json rpc"?
	// take body and wrap json rpc and marshall to string

	correlationID := uuid.NewV4().String()
	if err := amqp.Send(t.AMQPConn, rpc.Exchange, rpc.ExchangeType, rpc.BindingKey, "", correlationID, t.ReplyToQueue); err != nil {
		return nil, err //return some sort of response that is useful
	}
	//	response := JSONRPCResponse{}
	for {
		select {
		case d := <-deliveries:
			log.Infof("got replies correctionID 1: %s, correlationID2: %s", d.CorrelationId, correlationID)
			if d.CorrelationId == correlationID {
				// err := json.Unmarshal(d.Body, &response)
				// if err != nil {
				// 	log.Warnf("%s", err)
				// 	return nil, err
				// }

				response, err := NewJSONRPCResponseReader(d.Body)
				if err != nil {
					log.Warnf("%s", err)
					return nil, err
				}

				resp, err := NewResponse(response)
				if err != nil {
					log.Warnf("%s", err)
					return nil, err
				}
				return resp, nil
			}
		case <-time.After(t.Timeout):
			return nil, fmt.Errorf("Error: Message response timeout")
		}
	}
}

// NewClient returns a http.Client which implements our Transport
func NewClient(amqpConn *samqp.Connection, amqpRPCSettings AMQPRPCSettings) *http.Client {
	queue := random.RandQueue()

	return &http.Client{
		Transport: Transport{
			AMQPConn:        amqpConn,
			ReplyToQueue:    queue,
			AMQPRPCSettings: amqpRPCSettings,
			Timeout:         30 * time.Second,
		},
	}
}
