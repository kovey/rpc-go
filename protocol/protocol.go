package protocol

import (
	"encoding/json"
)

type Request struct {
	Path       string        `json:"p"`
	Method     string        `json:"m"`
	Args       []interface{} `json:"a"`
	TraceId    string        `json:"t"`
	SpanId     string        `json:"s"`
	From       string        `json:"f"`
	Version    string        `json:"v"`
	ClientLang string        `json:"c"`
}

type Response struct {
	Type   string      `json:"type"`
	Err    string      `json:"err"`
	Code   int32       `json:"code"`
	Trace  string      `json:"trace"`
	Result interface{} `json:"result"`
}

func NewRequest() *Request {
	return &Request{}
}

func NewRespose() Response {
	return Response{}
}

type ClientResponse struct {
	Type   string          `json:"type"`
	Err    string          `json:"err"`
	Code   int32           `json:"code"`
	Trace  string          `json:"trace"`
	Result json.RawMessage `json:"result"`
}

func NewClientResponse() *ClientResponse {
	return &ClientResponse{}
}
