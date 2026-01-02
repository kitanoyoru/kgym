package middlewares

import (
	"net/http"
	"time"

	"github.com/dromara/carbon/v2"
	"github.com/rs/zerolog/log"
)

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := carbon.Now()

		log.Info().
			Str("method", r.Method).
			Str("path", r.URL.Path).
			Str("request_id", r.Header.Get("X-Request-ID")).
			Msg("request received")

		ww := &responseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		next.ServeHTTP(ww, r)

		log.Info().
			Str("method", r.Method).
			Str("path", r.URL.Path).
			Int("status", ww.statusCode).
			Dur("latency", time.Since(startTime.StdTime())).
			Msg("request completed")
	})
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
