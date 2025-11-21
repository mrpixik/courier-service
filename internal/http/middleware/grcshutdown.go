package middleware

import (
	"context"
	"net/http"
)

// WithGracefulShutdown depreciated. Use BaseContext instead
func WithGracefulShutdown(shutdownCtx context.Context) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithCancel(r.Context())

			go func() {
				select {
				case <-shutdownCtx.Done():
					cancel()
				case <-ctx.Done():
				}
			}()

			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}
