package main

import (
	"context"
	"net/http"
	"time"
)

func TimeoutMW(fn http.HandlerFunc, timeout time.Duration) http.HandlerFunc {
	// timeout middleware
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), timeout)
		defer cancel()
		r = r.WithContext(ctx)
		go func() {
			fn(w, r)
		}()

	}
}
