package mq2http

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	samqp "github.com/streadway/amqp"
	"net/http"
)

type AMQPWriter struct {
	ReplyTo       string
	CorrelationID string
	Exchange      string
	ExchangeType  string
	AMQPConn      *samqp.Connection
}

func (w AMQPWriter) Header() http.Header {
	header := make(map[string][]string)
	return header
}

func (w AMQPWriter) Write(body []byte) (int, error) {
	exchangeType := w.ExchangeType
	exchange := w.Exchange
	correlationId := w.CorrelationID
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
			CorrelationId:   correlationId,
			ReplyTo:         replyTo,
			Body:            []byte(body),
			DeliveryMode:    samqp.Transient, // 1=non-persistent, 2=persistent
			Priority:        0,               // 0-9
		},
	)
	if err != nil {
		return 0, fmt.Errorf("failed to publish: %s", err)
	}
	return len(body), nil
}

func (w AMQPWriter) WriteHeader(status int) {
}

func AMQPFactory(ReplyTo, Exchange, ExchangeType, CorrelationID string, AMQPConn *samqp.Connection) (AMQPWriter, error) {
	return AMQPWriter{
		ReplyTo:       ReplyTo,
		CorrelationID: CorrelationID,
		Exchange:      Exchange,
		ExchangeType:  ExchangeType,
		AMQPConn:      AMQPConn,
	}, nil
}
