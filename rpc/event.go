package rpc

import (
	"bytes"
	"encoding/binary"
	"encoding/json"

	"github.com/kovey/logger-go/logger"
	"github.com/kovey/rpc-go/protocol"
	"github.com/kovey/server-go/server"
)

type RpcEvent struct {
}

func NewRpcEvent() RpcEvent {
	return RpcEvent{}
}

func (e RpcEvent) Close(s *server.Server, fd int64) {
	logger.Debug("connections[%d] disconnect", fd)
}

func (e RpcEvent) Connect(s *server.Server, fd int64) {
	logger.Debug("connections[%d] accept", fd)
}

func (e RpcEvent) Receive(s *server.Server, ev *server.Event) {
	logger.Debug("receive msg: %s", string(ev.Body()))
	request := protocol.NewRequest()
	json.Unmarshal(ev.Body(), request)
	if request.Path == "" || request.Method == "" {
		e.Error(s, ev.Fd(), "protocol_exception", "path or method is empty.", "")
		s.Close(ev.Fd())
		logger.Debug("path[%s] or method[%s] is empty", request.Path, request.Method)
		return
	}

	defer func(req *protocol.Request, begin float64) {
		data := make(map[string]interface{})
		end := time.Now().UnixNano()
		data["delay"] = end - begin
		data["request_time"] = begin
		data["type"] = "success"
		monitor.Write(data)
	}(request, time.Now().UnixNano())

	router, err := Get(request.Path, request.Method)
	if err != nil {
		logger.Debug("router[%s] is not exists", request.Path)
		e.Error(s, ev.Fd(), "exception", "router is not exists", "")
		s.Close(ev.Fd())
		return
	}

	result, err := router.Call(request.Args...)
	if err != nil {
		e.Error(s, ev.Fd(), "busi_exception", err.Error(), "")
		return
	}

	e.Success(s, ev.Fd(), result)
}

func (e RpcEvent) Success(s *server.Server, fd int64, result interface{}) {
	response := protocol.Response{
		Type:   "success",
		Err:    "",
		Code:   0,
		Trace:  "",
		Result: result,
	}

	logger.Debug("response: %v", response)
	send(s, fd, response)
}

func (e RpcEvent) Error(s *server.Server, fd int64, errType string, err string, trace string) {
	response := protocol.Response{
		Type:   errType,
		Err:    err,
		Code:   1000,
		Trace:  trace,
		Result: "",
	}

	send(s, fd, response)
}

func send(s *server.Server, fd int64, response protocol.Response) {
	buf, er := json.Marshal(response)
	if er != nil {
		logger.Debug("json marshal fail: %v, err: %s", response, er)
	}

	logger.Debug("response to client, package: %s", buf)
	length := len(buf)
	lBuf := Int32ToBytes(int32(length))
	buf = append(lBuf, buf...)
	s.Send(buf, fd)
}

func Int32ToBytes(n int32) []byte {
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, n)
	return bytesBuffer.Bytes()
}
