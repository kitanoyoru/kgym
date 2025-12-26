package middlewares

import (
	"net/http"
	"time"

	"github.com/dromara/carbon/v2"
	pkgLogger "github.com/kitanoyoru/kgym/pkg/logger"
	"go.uber.org/zap"
)

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger, err := pkgLogger.New(pkgLogger.WithDev())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		startTime := carbon.Now()

		logger.Info("request received",
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
			zap.String("request_id", r.Header.Get("X-Request-ID")),
		)

		r = r.WithContext(pkgLogger.Inject(r.Context(), logger))

		next.ServeHTTP(w, r)

		logger.Info("request completed",
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
			zap.Int("status", http.StatusOK),
			zap.Duration("latency", time.Since(startTime.StdTime())),
		)
	})
}
