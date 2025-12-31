package middlewares

import (
	"log"
	"net/http"
	"time"

	"github.com/dromara/carbon/v2"
)

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := carbon.Now()

		log.Printf("request received: method=%s path=%s request_id=%s",
			r.Method,
			r.URL.Path,
			r.Header.Get("X-Request-ID"),
		)

		next.ServeHTTP(w, r)

		log.Printf("request completed: method=%s path=%s status=%d latency=%v",
			r.Method,
			r.URL.Path,
			http.StatusOK,
			time.Since(startTime.StdTime()),
		)
	})
}
