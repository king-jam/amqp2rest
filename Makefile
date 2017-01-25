PROJECT = MQ-2-REST
BINARYNAME = mq2rest
GOOUT = ./bin

deps:
	go get github.com/onsi/ginkgo/ginkgo
	go get github.com/onsi/gomega
	go get ./...

integration-test:
	ginkgo -r -race -trace -cover -randomizeAllSpecs --slowSpecThreshold=30 --focus="\bINTEGRATION\b" -v

unit-test:
	ginkgo -r -race -trace -cover -randomizeAllSpecs --slowSpecThreshold=30 --focus="\bUNIT\b" -v

build:
	go build -o $(GOOUT)/$(BINARYNAME)
