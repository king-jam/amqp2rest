package mq2http

import (
	"encoding/json"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
	samqp "github.com/streadway/amqp"
)

// AMQPRPCSettings is ...
type AMQPRPCSettings struct {
	Exchange     string
	ExchangeType string
	QueueName    string
	BindingKey   string
	ConsumerTag  string
}

// AMQPWriter is ...
type AMQPWriter struct {
	ReplyTo       string
	CorrelationID string
	Exchange      string
	ExchangeType  string
	AMQPConn      *samqp.Connection
	Headers       map[string][]string
	Status        int
}

// Header is ...
func (w AMQPWriter) Header() http.Header {
	return w.Headers
}

// Write is ...
func (w AMQPWriter) Write(body []byte) (int, error) {
	exchangeType := w.ExchangeType
	exchange := w.Exchange
	correlationID := w.CorrelationID
	replyTo := w.ReplyTo
	channel, err := w.AMQPConn.Channel()
	if err != nil {
		return 0, fmt.Errorf("failed to create channel: %s", err)
	}

	log.Debugf("declaring %q Exchange (%q)", exchangeType, exchange)
	/*	err = channel.ExchangeDeclare(
			exchange,     // name
			exchangeType, // type
			true,         // durable
			false,        // auto-deleted
			false,        // internal
			false,        // noWait
			nil,          // arguments
		)
		if err != nil {
			return 0, fmt.Errorf("failed to declare exhange: %s", err)
		}*/
	respParamsStruct := RespParamsStruct{
		Status:        http.StatusText(w.Status),
		StatusCode:    w.Status,
		Header:        w.Headers,
		Body:          string(body),
		ContentLength: int64(len(body)),
		Close:         false,
		Uncompressed:  true,
		Request:       nil,
	}

	jsonRPC := JSONRPCResponse{
		Version: "2.0",
		Result:  respParamsStruct,
		Error:   string(w.Status),
		ID:      correlationID,
	}

	jsonRPCBody, err := json.Marshal(jsonRPC)
	if err != nil {
		return 0, err
	}

	log.Debugf("publishing %dB body (%q)", len(body), body)
	err = channel.Publish(
		"",      // publish to an exchange
		replyTo, // routing to 0 or more queues
		false,   // mandatory
		false,   // immediate
		samqp.Publishing{
			Headers:         samqp.Table{},
			ContentType:     "text/plain",
			ContentEncoding: "",
			CorrelationId:   correlationID,
			ReplyTo:         replyTo,
			Body:            jsonRPCBody,
			DeliveryMode:    samqp.Transient, // 1=non-persistent, 2=persistent
			Priority:        0,               // 0-9
		},
	)
	if err != nil {
		return 0, fmt.Errorf("failed to publish: %s", err)
	}
	return len(body), nil
}

// WriteHeader is ...
func (w AMQPWriter) WriteHeader(status int) {
	w.Status = status
}

// AMQPFactory is ...
func AMQPFactory(ReplyTo, Exchange, ExchangeType, CorrelationID string, AMQPConn *samqp.Connection) (AMQPWriter, error) {
	return AMQPWriter{
		ReplyTo:       ReplyTo,
		CorrelationID: CorrelationID,
		Exchange:      Exchange,
		ExchangeType:  ExchangeType,
		AMQPConn:      AMQPConn,
		// Headers       map[string][]string
		// Status        int
	}, nil
}
