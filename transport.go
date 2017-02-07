package mq2http

import (
  "fmt"
  "net/http"
  "time"

  uuid "github.com/satori/go.uuid"
  samqp "github.com/streadway/amqp"
)

// ExchangeTypes are the valid exchange types supported by this transport
var ExchangeTypes = []string{"topic", "direct", "fanout", "headers"}

const (
  // ReplyQueuePrefix is the mq2http reply to queue for the broker to identify
  // clients using this library
  ReplyQueuePrefix string = "mq2http"
)

// Transport is a ....
type Transport struct {
  AMQPConn     *samqp.Connection
  ReplyToQueue string
  ExchangeName string
  ExchangeType string
  BindingKey   string
  ConsumerTag  string
  Timeout      time.Duration
}

// RoundTrip is a ...
func (t Transport) RoundTrip(req *http.Request) (*http.Response, error) {
  // rpc := t.AMQPRPCSettings
  // // Listen
  // _, deliveries, err := amqp.Listen(t.AMQPConn, rpc.Exchange, rpc.ExchangeType, t.ReplyToQueue, rpc.BindingKey, rpc.ConsumerTag)
  // if err != nil {
  //   return nil, err
  // }
  // // Grab request and send over amqp using "json rpc"?
  // // take body and wrap json rpc and marshall to string
  //
  // correlationID := uuid.NewV4().String()
  // if err := amqp.Send(t.AMQPConn, rpc.Exchange, rpc.ExchangeType, rpc.BindingKey, "", correlationID, t.ReplyToQueue); err != nil {
  //   return nil, err //return some sort of response that is useful
  // }
  // //	response := JSONRPCResponse{}
  // for {
  //   select {
  //   case d := <-deliveries:
  //     log.Infof("got replies correctionID 1: %s, correlationID2: %s", d.CorrelationId, correlationID)
  //     if d.CorrelationId == correlationID {
  //       // err := json.Unmarshal(d.Body, &response)
  //       // if err != nil {
  //       // 	log.Warnf("%s", err)
  //       // 	return nil, err
  //       // }
  //
  //       response, err := NewJSONRPCResponseReader(d.Body)
  //       if err != nil {
  //         log.Warnf("%s", err)
  //         return nil, err
  //       }
  //
  //       resp, err := NewResponse(response)
  //       if err != nil {
  //         log.Warnf("%s", err)
  //         return nil, err
  //       }
  //       return resp, nil
  //     }
  //   case <-time.After(t.Timeout):
  //     return nil, fmt.Errorf("Error: Message response timeout")
  //   }
  // }
  return &http.Response{}, nil
}

// // NewClient returns a http.Client which implements our Transport
// func NewClient(amqpConn *samqp.Connection, amqpRPCSettings AMQPRPCSettings) *http.Client {
//   randUUID := uuid.NewV4().String()
//   queue := ReplyQueuePrefix + "-" + randUUID
//
//   return &http.Client{
//     Transport: Transport{
//       AMQPConn:        amqpConn,
//       ReplyToQueue:    queue,
//       AMQPRPCSettings: amqpRPCSettings,
//       Timeout:         30 * time.Second,
//     },
//   }
// }

// NewTransport returns an http.Transport that can be used to create an http.Client
func NewTransport(amqpConn *samqp.Connection, exchangeName, exchangeType, bindingKey, consumerID string, timeout time.Duration) (*Transport, error) {

  randUUID := uuid.NewV4().String()
  replyToQueue := ReplyQueuePrefix + "-" + randUUID

  // verify we have a declared exchangeName
  if exchangeName == "" {
    return &Transport{}, fmt.Errorf("MQ2HTTP: Transport: Exchange Name not valid")
  }

  // verify we have a valid exchangeType
  if exchangeType == "" || !contains(ExchangeTypes, exchangeType) {
    return &Transport{}, fmt.Errorf("MQ2HTTP: Transport: Invalid Exchange Type")
  }

  // default to a wildcard bindingKey if not specified
  if bindingKey == "" {
    bindingKey = "#"
  }

  // if this is left empty, create a random ID that matches our random queue
  if consumerID == "" {
    consumerID = "cid-" + ReplyQueuePrefix + "-" + randUUID
  }

  if timeout < 0*time.Second {
    return &Transport{}, fmt.Errorf("MQ2HTTP: Transport: Invalid Timeout")
  }

  return &Transport{
    AMQPConn:     amqpConn,
    ReplyToQueue: replyToQueue,
    ExchangeName: exchangeName,
    ExchangeType: exchangeType,
    BindingKey:   bindingKey,
    ConsumerTag:  consumerID,
    Timeout:      timeout,
  }, nil
}

// contains is a slice check helper
func contains(s []string, e string) bool {
  for _, a := range s {
    if a == e {
      return true
    }
  }
  return false
}
