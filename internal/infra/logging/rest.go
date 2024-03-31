package logging

import (
	"log/slog"
	"net/http"
	"time"
)

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}

// Middleware creates a middleware for logging HTTP requests.
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &statusRecorder{ResponseWriter: w}
		next.ServeHTTP(rec, r)
		slog.Debug("",
			"statusCode", rec.status,
			"latency", time.Since(start),
			"clientIP", r.RemoteAddr,
			"httpMethod", r.Method,
			"path", r.RequestURI,
		)
	})
}
