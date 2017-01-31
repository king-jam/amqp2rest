package mq2http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

/*
This is a psuedo custom JSON-RPC struct used for this
{
        "jsonrpc": "2.0",
        "method": "GET /v1/nodes",
        "params": {
                "body": "",
                "headers": {
                        "Content-Type": "application/json"
		}
        },
        "id": "1238814hnfasdf1afdf"
}
*/

/*
This is a psudo custom JSON-RPC struct used for a resp
{
		"jsonrpc": "2.0",
		"result": "RespParamsStruct"
		"error": "spec error code"
		"id": "1234567890qwertyuii"
}
*/

// RespParamsStruct is ....
type RespParamsStruct struct {
	Status     string `json:"status"` // Have this in resp
	StatusCode int    `json:"statuscode"`
	// Proto string
	// ProtoMajor int
	// ProtoMinor int
	Header        map[string][]string `json:"header"` // have this in amqp.go mq2http
	Body          string              `json:"body"`
	ContentLength int64               `json:"contentlength"` // = len(body)
	//	TransferEncoding []string
	Close        bool `json:"close"`        // always false
	Uncompressed bool `json:"uncompressed"` // always true
	// 	Trailer      http.Header          `json:"trailer"`      // same as headers
	Request *http.Request `json:"reqpest"` // nil
	//	TLS     *tls.ConnectionState `json:"tls"`     // nil
}

// JSONRPCRequest is ...
type JSONRPCRequest struct {
	Version string          `json:"jsonrpc"`
	Method  string          `json:"method"`
	Params  ReqParamsStruct `json:"params,omitempty"`
	ID      string          `json:"id,omitempty"`
}

// ReqParamsStruct is ...
type ReqParamsStruct struct {
	Body    string              `json:"body"`
	Headers map[string][]string `json:"headers"`
}

// JSONRPCResponse is ...
type JSONRPCResponse struct {
	Version string           `json:"jsonrpc"`
	Result  RespParamsStruct `json:"result,omitempty"`
	Error   string           `json:"error,omitempty"`
	ID      string           `json:"id,emitempty"`
}

// JSONRPCResponseReader is ...
type JSONRPCResponseReader struct {
	status     string
	statusCode int
	headers    map[string][]string
	reader     io.ReadCloser
}

// NewJSONRPCResponseReader is ...
func NewJSONRPCResponseReader(b []byte) (*JSONRPCResponseReader, error) {
	var resp JSONRPCResponse
	json.Unmarshal(b, &resp)

	r := bytes.NewReader([]byte(resp.Result.Body))
	rc := ioutil.NopCloser(r)
	return &JSONRPCResponseReader{
		status:     resp.Result.Status,
		statusCode: resp.Result.StatusCode,
		headers:    resp.Result.Header,
		reader:     rc,
	}, nil
}

// Status is ...
func (j *JSONRPCResponseReader) Status() string {
	return j.status
}

// StatusCode is ...
func (j *JSONRPCResponseReader) StatusCode() int {
	return j.statusCode
}

// Headers is ...
func (j *JSONRPCResponseReader) Headers() http.Header {
	header := http.Header{}
	for k, v := range j.headers {
		header.Set(k, v[0]) // TODO: needs to change for more complex headers
	}
	return header
}

// Reader is ...
func (j *JSONRPCResponseReader) Reader() io.ReadCloser {
	return j.reader
}

// JSONRPCRequestReader is ...
type JSONRPCRequestReader struct {
	headers map[string][]string
	method  string
	url     string
	reader  io.Reader
}

// NewJSONRPCRequestReader is ...
func NewJSONRPCRequestReader(b []byte) (*JSONRPCRequestReader, error) {
	var req JSONRPCRequest
	json.Unmarshal(b, &req)
	methodURL := strings.Split(req.Method, " ")
	if len(methodURL) != 2 {
		return &JSONRPCRequestReader{}, fmt.Errorf("JSONRPC: Request decode failed: Method Invalid")
	}
	r := bytes.NewReader([]byte(req.Params.Body))
	return &JSONRPCRequestReader{
		headers: req.Params.Headers,
		method:  methodURL[0],
		url:     methodURL[1],
		reader:  r,
	}, nil
}

// Headers is ...
func (j *JSONRPCRequestReader) Headers() map[string][]string {
	return j.headers
}

// Method is ...
func (j *JSONRPCRequestReader) Method() string {
	return j.method
}

// URL is ...
func (j *JSONRPCRequestReader) URL() string {
	return j.url
}

// Reader is ...
func (j *JSONRPCRequestReader) Reader() io.Reader {
	return j.reader
}
