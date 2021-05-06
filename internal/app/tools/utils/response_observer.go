package utils

import (
	"bufio"
	"errors"
	"net"
	"net/http"
)

type ResponseObserver struct {
	http.ResponseWriter
	Status      int
	Written     int64
	WroteHeader bool
}

func (o *ResponseObserver) Write(p []byte) (n int, err error) {
	if !o.WroteHeader {
		o.WriteHeader(http.StatusOK)
	}

	n, err = o.ResponseWriter.Write(p)
	o.Written += int64(n)

	return
}

func (o *ResponseObserver) WriteHeader(code int) {
	o.ResponseWriter.WriteHeader(code)
	if o.WroteHeader {
		return
	}

	o.WroteHeader = true
	o.Status = code
}

func (o *ResponseObserver) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	h, ok := o.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, errors.New("hijack not supported")
	}

	return h.Hijack()
}
