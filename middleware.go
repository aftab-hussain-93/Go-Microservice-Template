package main

import (
	"context"
	"net/http"
	"time"
)

type myResponseWriter struct {
	w http.ResponseWriter

	bodyWritten   bool
	headerWritten bool
}

func (r *myResponseWriter) Header() http.Header {
	return r.w.Header()
}

func (r *myResponseWriter) Write(ip []byte) (int, error) {
	if !r.bodyWritten {
		n, err := r.w.Write(ip)
		r.bodyWritten = true
		return n, err
	}
	return 0, nil
}

func (r *myResponseWriter) WriteHeader(statusCode int) {
	if !r.headerWritten {
		r.w.WriteHeader(statusCode)
		r.headerWritten = true
	}

}

func TimeoutMW(fn http.HandlerFunc, timeout time.Duration) http.HandlerFunc {
	// timeout middleware
	return func(w http.ResponseWriter, r *http.Request) {
		ww := &myResponseWriter{w: w}
		ctx, cancel := context.WithTimeout(r.Context(), timeout)
		defer cancel()
		panicChan := make(chan struct{})
		doneChan := make(chan struct{})
		r = r.WithContext(ctx)
		go func() {
			defer func() {
				if e := recover(); e != nil {
					panicChan <- struct{}{}
				}
			}()
			fn(ww, r)
			doneChan <- struct{}{}
		}()

		select {
		case <-ctx.Done():
			// write timeout response
			ww.WriteHeader(http.StatusRequestTimeout)
			_, _ = w.Write([]byte("timeout"))
		case <-panicChan:
			// write error response
			ww.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte("panic"))
		case <-doneChan:
			//
		}
	}
}
