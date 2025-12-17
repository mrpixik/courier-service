package middleware

import (
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
	"strconv"
	"time"
)

type MetricsObserverHTTP interface {
	IncTotalRequests()
	NewRequest(method, path, status string, durationSec float64)
}

func WithMonitoring(log *slog.Logger, obs MetricsObserverHTTP) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		log = log.With(
			slog.String("component", "middleware/logger"),
		)

		log.Info("logger middleware is enabled")
		fn := func(w http.ResponseWriter, r *http.Request) {
			entry := log.With(
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.String("remote_addr", r.RemoteAddr),
				slog.String("user_agent", r.UserAgent()),
			)

			obs.IncTotalRequests()

			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			start := time.Now()
			defer func() {
				entry.Info("request completed",
					slog.Int("status", ww.Status()),
					slog.Int("bytes", ww.BytesWritten()),
					slog.Duration("time", time.Since(start)),
				)

				obs.NewRequest(r.Method, r.URL.Path, strconv.Itoa(ww.Status()), time.Since(start).Seconds())
			}()

			next.ServeHTTP(ww, r)
		}

		return http.HandlerFunc(fn)
	}
}
