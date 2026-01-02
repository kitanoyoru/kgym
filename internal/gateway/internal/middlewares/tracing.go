package middlewares

import (
	"net/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"
)

const (
	GatewayTracerName = "kgym.gateway.http"
)

func Tracing(next http.Handler) http.Handler {
	tracer := otel.Tracer(GatewayTracerName)
	propagator := otel.GetTextMapPropagator()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := propagator.Extract(r.Context(), propagation.HeaderCarrier(r.Header))

		spanName := r.Method + " " + r.URL.Path
		ctx, span := tracer.Start(ctx, spanName, trace.WithSpanKind(trace.SpanKindServer))
		defer span.End()

		span.SetAttributes(
			semconv.HTTPMethodKey.String(r.Method),
			semconv.HTTPURLKey.String(r.URL.String()),
			semconv.HTTPRouteKey.String(r.URL.Path),
		)

		if r.UserAgent() != "" {
			span.SetAttributes(attribute.String("http.user_agent", r.UserAgent()))
		}

		if r.RemoteAddr != "" {
			span.SetAttributes(attribute.String("http.client_ip", r.RemoteAddr))
		}

		requestID := r.Header.Get("X-Request-ID")
		if requestID != "" {
			span.SetAttributes(attribute.String("http.request_id", requestID))
		}

		platform := r.Header.Get("X-Platform")
		if platform != "" {
			span.SetAttributes(attribute.String("http.platform", platform))
		}

		appVersion := r.Header.Get("X-App-Version")
		if appVersion != "" {
			span.SetAttributes(attribute.String("http.app_version", appVersion))
		}

		propagator.Inject(ctx, propagation.HeaderCarrier(r.Header))

		r = r.WithContext(ctx)

		ww := &responseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		next.ServeHTTP(ww, r)

		span.SetAttributes(semconv.HTTPStatusCodeKey.Int(ww.statusCode))

		if ww.statusCode >= http.StatusBadRequest {
			span.SetStatus(codes.Error, http.StatusText(ww.statusCode))
		} else {
			span.SetStatus(codes.Ok, "")
		}
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
