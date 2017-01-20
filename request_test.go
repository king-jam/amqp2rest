package mq2http_test

import (
	. "github.com/king-jam/mq2http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("HTTP Request Interface", func() {
	Context("MSG Decoding", func() {
		It("UNIT should properly decode a valid msg", func() {
			msg := []byte(`{"jsonrpc":"2.0","method":"POST /v1/test/route/idthing","params":{"body":"<encoded JSON body>","headers":{"Content-Type":"application/json","Accept":"application/json"}},"id":"1238814hnfasdf1afdf"}`)
			_, err := NewJSONRPC(msg)
			Expect(err).ToNot(HaveOccurred())
		})
		It("UNIT should properly create the Headers", func() {
			msg := []byte(`{"jsonrpc":"2.0","method":"POST /v1/test/route/idthing","params":{"body":"<encoded JSON body>","headers":{"Content-Type":"application/json","Accept":"application/json"}},"id":"1238814hnfasdf1afdf"}`)
			expected := map[string]string{
				"Content-Type": "application/json",
				"Accept":       "application/json",
			}
			r, err := NewJSONRPC(msg)
			Expect(err).ToNot(HaveOccurred())
			Expect(r.Headers()).To(Equal(expected))
		})
		It("UNIT should properly create the Method", func() {
			msg := []byte(`{"jsonrpc":"2.0","method":"POST /v1/test/route/idthing","params":{"body":"<encoded JSON body>","headers":{"Content-Type":"application/json","Accept":"application/json"}},"id":"1238814hnfasdf1afdf"}`)
			expected := "POST"
			r, err := NewJSONRPC(msg)
			Expect(err).ToNot(HaveOccurred())
			Expect(r.Method()).To(Equal(expected))
		})
		It("UNIT should properly create the URL", func() {
			msg := []byte(`{"jsonrpc":"2.0","method":"POST /v1/test/route/idthing","params":{"body":"<encoded JSON body>","headers":{"Content-Type":"application/json","Accept":"application/json"}},"id":"1238814hnfasdf1afdf"}`)
			expected := "/v1/test/route/idthing"
			r, err := NewJSONRPC(msg)
			Expect(err).ToNot(HaveOccurred())
			Expect(r.URL()).To(Equal(expected))
		})
		It("UNIT should properly create the Body", func() {
			msg := []byte(`{"jsonrpc":"2.0","method":"POST /v1/test/route/idthing","params":{"body":"<encoded JSON body>","headers":{"Content-Type":"application/json","Accept":"application/json"}},"id":"1238814hnfasdf1afdf"}`)
			r, err := NewJSONRPC(msg)
			Expect(err).ToNot(HaveOccurred())
			Expect(r.Body()).To(BeAssignableToTypeOf(r.Body()))
		})
	})
})
