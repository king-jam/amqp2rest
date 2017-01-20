package mq2http_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestMq2http(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Mq2http Suite")
}
