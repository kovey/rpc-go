package rpc

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"time"

	"github.com/kovey/logger-go/logger"
	"github.com/kovey/logger-go/monitor"
	"github.com/kovey/rpc-go/protocol"
	"github.com/kovey/server-go/server"
	"github.com/kovey/server-go/util"
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
	monLog := getMonitor(request)
	spanId := util.SpanId()
	monLog.SpanId = spanId
	defer func(monLog *monitor.Monitor) {
		err := recover()
		monLog.Trace = getTrace(err)
		monLog.Err = ""
		if err != nil {
			s.Close(ev.Fd())
			monLog.Err = fmt.Sprintf("%s", err)
		}

		monLog.End = time.Now().UnixNano() / 1e6
		monLog.Delay = float64(monLog.End-monLog.RequestTime) / 1e6
		monitor.Write(*monLog)
	}(monLog)

	if request.Path == "" || request.Method == "" {
		e.Error(s, ev.Fd(), "protocol_exception", "path or method is empty.", "")
		s.Close(ev.Fd())
		logger.Debug("path[%s] or method[%s] is empty", request.Path, request.Method)
		monLog.Type = "protocol_exception"
		monLog.Err = "path or method is empty."
		monLog.Trace = ""
		return
	}

	router, err := Get(request.Path, request.Method)
	if err != nil {
		logger.Debug("router[%s] is not exists", request.Path)
		e.Error(s, ev.Fd(), "exception", "router is not exists", "")
		s.Close(ev.Fd())
		monLog.Type = "exception"
		monLog.Err = "router is not exists"
		monLog.Trace = ""
		return
	}

	result, err := router.Call(request, spanId)
	if err != nil {
		e.Error(s, ev.Fd(), "busi_exception", err.Error(), "")
		monLog.Type = "busi_exception"
		monLog.Err = err.Error()
		monLog.Trace = ""
		return
	}

	monLog.Response = result
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
