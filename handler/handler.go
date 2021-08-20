package handler

import "github.com/kovey/logger-go/logger"

type Handler struct {
	traceId string
	spandId string
}

type HandlerInterface interface {
	SetTraceId(traceId string)
	SetSpandId(spanId string)
	TraceId() string
	SpandId() string
}

func NewHandler(traceId string, spandId string) Handler {
	return Handler{traceId: traceId, spandId: spandId}
}

func (h Handler) SetTraceId(traceId string) {
	logger.Debug("SetTraceId(%s)", traceId)
	h.traceId = traceId
}

func (h Handler) SetSpandId(spandId string) {
	logger.Debug("SetSpandId(%s)", spandId)
	h.spandId = spandId
}

func (h Handler) TraceId() string {
	return h.traceId
}

func (h Handler) SpandId() string {
	return h.spandId
}
