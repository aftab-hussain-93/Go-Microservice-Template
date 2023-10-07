package main

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	err "github.com/aftab-hussain-93/crypto-price-finder-microservice/error"
)

func TimeoutAndPanicMW(fn http.HandlerFunc, timeout time.Duration) http.HandlerFunc {
	timeoutErr := err.Err{
		Code:    err.ErrTimeout,
		Message: "Request timeout",
	}
	// timeout middleware
	return func(w http.ResponseWriter, r *http.Request) {
		ww := &responseWriter{w: w}
		ctx, cancel := context.WithTimeout(r.Context(), timeout)
		defer cancel()
		panicChan := make(chan error)
		doneChan := make(chan struct{})
		r = r.WithContext(ctx)
		go func() {
			defer func() {
				err := recover()
				if err != nil {
					e := err.(error)
					panicChan <- fmt.Errorf("panic %w", e)
				}
			}()
			fn(ww, r)
			doneChan <- struct{}{}
		}()

		select {
		case <-ctx.Done():
			status, errr := convertErrToErrResponse(ctx, &timeoutErr)
			// write timeout response
			_ = writeJSON(ww, status, errr)
		case e := <-panicChan:
			status, errr := convertErrToErrResponse(ctx, &err.Err{
				Err:  e,
				Code: err.ErrInternal,
			})
			// write timeout response
			_ = writeJSON(ww, status, errr)
		case <-doneChan:
			// already written
		}
	}
}

type responseWriter struct {
	w  http.ResponseWriter
	mu sync.Mutex

	bodyWritten   bool
	headerWritten bool
}

func (r *responseWriter) Header() http.Header {
	return r.w.Header()
}

func (r *responseWriter) Write(ip []byte) (int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if !r.bodyWritten {
		n, err := r.w.Write(ip)
		r.bodyWritten = true
		return n, err
	}
	return 0, nil
}

func (r *responseWriter) WriteHeader(statusCode int) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if !r.headerWritten {
		r.w.WriteHeader(statusCode)
		r.headerWritten = true
	}
}
