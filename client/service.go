package client

import (
	"encoding/json"
	"errors"

	"github.com/kovey/rpc-go/protocol"
	"github.com/kovey/server-go/server"
)

type Service struct {
	cli        *Client
	traceId    string
	spanId     string
	from       string
	version    string
	clientLang string
	path       string
}

func NewService(host string, port int, path string, traceId string, spanId string, from string) Service {
	return Service{
		cli:  NewClient(host, port, server.NewConfig(true, 0, 4, server.INT_32)),
		path: path, traceId: traceId, spanId: spanId, from: from,
		version: "1.0.0", clientLang: "golang",
	}
}

func (s Service) Call(method string, args ...interface{}) (json.RawMessage, error) {
	if !s.cli.isConnected {
		err := s.cli.Connect()
		if err != nil {
			return nil, err
		}
	}

	request := protocol.NewRequest()
	request.Path = s.path
	request.Method = method
	request.Args = args
	request.TraceId = s.traceId
	request.SpanId = s.spanId
	request.From = s.from
	request.Version = s.version
	request.ClientLang = s.clientLang
	err := s.cli.Send(*request)
	if err != nil {
		return nil, err
	}

	response, e := s.cli.Recv()
	if e != nil {
		return nil, e
	}

	if response.Type != "success" {
		return nil, errors.New(response.Err)
	}

	return response.Result, nil
}

func (s Service) Close() {
	s.cli.Close()
}
