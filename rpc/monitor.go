package rpc

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/kovey/logger-go/logger"
	"github.com/kovey/logger-go/monitor"
	"github.com/kovey/rpc-go/protocol"
)

func getMonitor(request *protocol.Request) *monitor.Monitor {
	monLog := &monitor.Monitor{}
	monLog.Path = request.Path
	monLog.Params = request.Args
	monLog.RequestTime = time.Now().UnixNano() / 1e6
	monLog.ServiceType = "rpc"
	monLog.Class = request.Path
	monLog.Method = request.Method
	monLog.Args = request.Args
	monLog.Ip = ""
	monLog.Time = time.Now().Unix()
	monLog.Timestamp = time.Now().Format("2006-01-02 15:04:05")
	minute, err := strconv.ParseInt(time.Now().Format("200601021504"), 10, 64)
	if err != nil {
		monLog.Minute = minute
	}
	monLog.HttpCode = 200
	monLog.TraceId = request.TraceId
	monLog.SpanId = ""
	monLog.ParentId = request.SpanId
	monLog.ClientVersion = request.Version
	monLog.ServerVersion = "1.0.0"
	monLog.ClientLang = request.ClientLang
	monLog.ServerLang = "golang"
	monLog.From = request.From

	return monLog
}

func getTrace(err interface{}) string {
	if err == nil {
		logger.Debug("err is nil")
		return ""
	}

	logger.Error("panic error[%s]", err)

	traces := make([]string, 1)
	traces[0] = fmt.Sprintf("panic error[%s]", err)

	for i := 3; ; i++ {
		_, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		traces = append(traces, fmt.Sprintf("%s(%d)", file, line))
		logger.Error("%s(%d)", file, line)
	}

	return strings.Join(traces, "#")
}
