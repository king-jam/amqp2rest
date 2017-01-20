package mq2http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

/*
This is a psuedo custom JSON-RPC struct used for this
{
        "jsonrpc": "2.0",
        "method": "POST /v1/test/route/idthing",
        "params": {
                "body": "<encoded JSON body>",
                "headers": {
                        "Content-Type": "application/json",
                        "Accept": "application/json"
		}
        },
        "id": "1238814hnfasdf1afdf"
}
*/
type JSONRPCRequest struct {
	Version string       `json: "jsonrpc"`
	Method  string       `json: "method"`
	Params  ParamsStruct `json: "params,omitempty"`
	ID      string       `json: "id"`
}

type ParamsStruct struct {
	Body    string            `json: "body"`
	Headers map[string]string `json: "headers"`
}

type JSONRPC struct {
	headers map[string]string
	method  string
	url     string
	reader  io.Reader
}

func NewJSONRPC(b []byte) (*JSONRPC, error) {
	var req JSONRPCRequest
	json.Unmarshal(b, &req)
	methodURL := strings.Split(req.Method, " ")
	if len(methodURL) != 2 {
		return &JSONRPC{}, fmt.Errorf("JSONRPC: Request decode failed: Method Invalid")
	}
	r := bytes.NewReader([]byte(req.Params.Body))
	return &JSONRPC{
		headers: req.Params.Headers,
		method:  methodURL[0],
		url:     methodURL[1],
		reader:  r,
	}, nil
}

func (j *JSONRPC) Headers() map[string]string {
	return j.headers
}
func (j *JSONRPC) Method() string {
	return j.method
}

func (j *JSONRPC) URL() string {
	return j.url
}

func (j *JSONRPC) Body() io.Reader {
	return j.reader
}
