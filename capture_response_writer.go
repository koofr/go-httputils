package httputils

import (
	"bufio"
	"errors"
	"net"
	"net/http"
	"time"
)

var ErrResponseNotHijacker = errors.New("response does not implement http.Hijacker")

type CaptureResponseWriter struct {
	http.ResponseWriter
	StatusCode     int
	ResponseLength int64
	Start          time.Time
	WriteStart     time.Time
	Hijacked       bool
	HeaderWritten  bool
}

func NewCaptureResponseWriter(w http.ResponseWriter) *CaptureResponseWriter {
	return &CaptureResponseWriter{
		ResponseWriter: w,
		StatusCode:     http.StatusOK,
		ResponseLength: 0,
		Start:          time.Now().UTC(),
		Hijacked:       false,
		HeaderWritten:  false,
	}
}

func (w *CaptureResponseWriter) Write(buf []byte) (int, error) {
	if w.WriteStart.IsZero() {
		w.WriteStart = time.Now()
	}
	if !w.HeaderWritten {
		w.HeaderWritten = true
	}
	if w.Hijacked {
		panic("Write on hijacked CaptureResponseWriter")
	}
	n, err := w.ResponseWriter.Write(buf)
	w.ResponseLength += int64(n)
	return n, err
}

func (w *CaptureResponseWriter) WriteHeader(statusCode int) {
	if w.HeaderWritten {
		panic("header already written")
	}
	if w.Hijacked {
		panic("WriteHeader on hijacked CaptureResponseWriter")
	}
	w.StatusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
	w.HeaderWritten = true
}

func (w *CaptureResponseWriter) Duration() time.Duration {
	return time.Since(w.Start)
}

func (w *CaptureResponseWriter) WriteDuration() time.Duration {
	if w.WriteStart.IsZero() {
		return 0
	}
	return time.Since(w.WriteStart)
}

func (w *CaptureResponseWriter) Flush() {
	if flusher, ok := w.ResponseWriter.(http.Flusher); ok {
		flusher.Flush()
	}
}

func (w *CaptureResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if hijacker, ok := w.ResponseWriter.(http.Hijacker); ok {
		w.Hijacked = true
		return hijacker.Hijack()
	}

	return nil, nil, ErrResponseNotHijacker
}
