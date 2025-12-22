package middleware

import (
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
	"strconv"
	"time"
)

type MetricsObserverHTTP interface {
	IncTotalRequests()
	NewRequest(method, path, status string, durationSec float64)
}

func WithMetrics(obs MetricsObserverHTTP) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {

		fn := func(w http.ResponseWriter, r *http.Request) {

			obs.IncTotalRequests()

			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			start := time.Now()
			defer func() {

				obs.NewRequest(r.Method, r.URL.Path, strconv.Itoa(ww.Status()), time.Since(start).Seconds())

			}()

			next.ServeHTTP(ww, r)
		}

		return http.HandlerFunc(fn)
	}
}
