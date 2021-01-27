package httputils

import (
	"bufio"
	"net"
	"net/http"
)

type CallbackResponseWriter struct {
	http.ResponseWriter
	beforeWriteHeader func()

	beforeWriteHeaderCalled bool
}

func NewCallbackResponseWriter(w http.ResponseWriter, beforeWriteHeader func()) *CallbackResponseWriter {
	return &CallbackResponseWriter{
		ResponseWriter:    w,
		beforeWriteHeader: beforeWriteHeader,

		beforeWriteHeaderCalled: false,
	}
}

func (w *CallbackResponseWriter) WriteHeader(statusCode int) {
	w.beforeWriteHeader()
	w.beforeWriteHeaderCalled = true

	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *CallbackResponseWriter) Write(b []byte) (int, error) {
	if !w.beforeWriteHeaderCalled {
		w.beforeWriteHeader()
		w.beforeWriteHeaderCalled = true
	}

	return w.ResponseWriter.Write(b)
}

func (w *CallbackResponseWriter) Flush() {
	if flusher, ok := w.ResponseWriter.(http.Flusher); ok {
		flusher.Flush()
	}
}

func (w *CallbackResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if hijacker, ok := w.ResponseWriter.(http.Hijacker); ok {
		return hijacker.Hijack()
	}

	return nil, nil, ErrResponseNotHijacker
}

func (w *CallbackResponseWriter) Done() {
	if !w.beforeWriteHeaderCalled {
		w.beforeWriteHeader()
		w.beforeWriteHeaderCalled = true
	}
}

var _ http.ResponseWriter = (*CallbackResponseWriter)(nil)
