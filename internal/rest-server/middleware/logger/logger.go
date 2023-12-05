package logger

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
)

func New(log zerolog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		log = log.With().Str("component", "middleware/logger").Logger()

		log.Info().Msg("logger middleware enabled")

		fn := func(w http.ResponseWriter, r *http.Request) {
			entry := log.With().
				Str("method", r.Method).
				Str("path", r.URL.Path).
				Str("remote_addr", r.RemoteAddr).
				Str("request_id", middleware.GetReqID(r.Context())).
				Logger()

			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			t1 := time.Now()
			defer func() {
				entry.Info().
					Int("status", ww.Status()).
					Int("bytes", ww.BytesWritten()).
					Str("duration", time.Since(t1).String()).
					Msg("request completed")
			}()

			next.ServeHTTP(ww, r)
		}

		return http.HandlerFunc(fn)
	}
}
