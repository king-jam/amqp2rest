package mq2http_test

import (
  //	. "github.com/king-jam/mq2http"

  "time"

  "github.com/king-jam/mq2http"
  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"
  samqp "github.com/streadway/amqp"
)

var _ = Describe("mq2http Client RoundTripper", func() {
  var transport *mq2http.Transport
  var err error
  var mqConnection *samqp.Connection
  var mqTimeout time.Duration
  var mqExchangeName string
  var mqExchangeType string
  var mqBindingKey string
  var mqConsumerID string

  BeforeEach(func() {
    mqConnection = &samqp.Connection{}
    mqTimeout = 500 * time.Millisecond
    mqExchangeName = "test-exchange"
    mqExchangeType = "topic"
    mqBindingKey = "test.bind.key"
    mqConsumerID = ""
  })

  JustBeforeEach(func() {
    transport, err = mq2http.NewTransport(
      mqConnection,
      mqExchangeName,
      mqExchangeType,
      mqBindingKey,
      mqConsumerID,
      mqTimeout,
    )
  })

  Describe("Create a RoundTripper", func() {
    Context("Creating MQ2HTTP Transport with default parameters", func() {

      It("UNIT can successfully create a transport", func() {
        Expect(err).ToNot(HaveOccurred())
        Expect(transport).To(BeAssignableToTypeOf(&mq2http.Transport{}))
        Expect(transport.ReplyToQueue).To(ContainSubstring("mq2http"))
        Expect(transport.Timeout).To(Equal(500 * time.Millisecond))
      })
    })

    Context("Creating MQ2HTTP Transport with bad timeout", func() {

      BeforeEach(func() {
        mqTimeout = -1
      })

      It("UNIT should generate an error for invalid timeout", func() {
        Expect(err).To(HaveOccurred())
      })
    })

    Context("Creating MQ2HTTP Transport with invalid exchange name", func() {
      BeforeEach(func() {
        mqExchangeName = ""
      })

      It("UNIT should generate an error for invalid exchange name", func() {
        Expect(err).To(HaveOccurred())
      })
    })

    Context("Creating MQ2HTTP Transport with invalid exchange type", func() {
      BeforeEach(func() {
        mqExchangeType = "notype"
      })

      It("UNIT should generate an error for invalid exchange type", func() {
        Expect(err).To(HaveOccurred())
      })
    })

    Context("Creating MQ2HTTP Transport with empty binding key", func() {
      BeforeEach(func() {
        mqBindingKey = ""
      })

      It("UNIT should default to wildcard binding key", func() {
        Expect(err).ToNot(HaveOccurred())
        Expect(transport.BindingKey).To(Equal("#"))
      })
    })

    Context("Creating MQ2HTTP Transport with empty consumer ID", func() {
      BeforeEach(func() {
        mqConsumerID = ""
      })

      It("UNIT should generate a consumer ID", func() {
        Expect(err).ToNot(HaveOccurred())
        Expect(transport.ConsumerTag).To(ContainSubstring("cid-mq2http"))
      })
    })

    Context("Creating MQ2HTTP Transport with consumer ID", func() {
      BeforeEach(func() {
        mqConsumerID = "test-consumer"
      })

      It("UNIT should use provided consumer ID", func() {
        Expect(err).ToNot(HaveOccurred())
        Expect(transport.ConsumerTag).To(Equal("test-consumer"))
      })
    })

    Context("Creating MQ2HTTP Transport with invalid connection", func() {

      BeforeEach(func() {
        mqConnection = nil
      })

      It("INTEGRATION should generate an error for no connection", func() {
        Expect(err).To(HaveOccurred())
      })
    })

    Context("Other stuff", func() {
      It("UNIT should take valid connection data", func() {

      })
      It("UNIT should validate MQ settings input", func() {

      })
    })

    Context("Can successfully make HTTP 2 MQ Calls", func() {
      It("INTEGRATION can successfully create a connection to a HTTP server", func() {
      })
      It("INTEGRATION should successfully connect to MQ broker", func() {

      })
      It("INTEGRATION should successfully setup MQ broker requirements", func() {

      })
      It("INTEGRATION should take an HTTP Request", func() {

      })
      It("INTEGRATION should validate the HTTP Request", func() {

      })
      It("INTEGRATION should send the request to the server", func() {

      })
      It("INTEGRATION should provide a valid HTTP Response", func() {

      })
      It("INTEGRATION should receive a encoded response from server", func() {

      })
    })
  })
})
