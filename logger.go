package dise

import (
	"context"
	"log/slog"
	"net/http"
	"time"
)

type loggerResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *loggerResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func CanonicalLogger(logger *slog.Logger) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			startedAt := time.Now()
			attributes := map[string]interface{}{}
			wrappedWriter := &loggerResponseWriter{ResponseWriter: w}
			defer func() {
				var logAttrs []any
				for key, value := range attributes {
					logAttrs = append(logAttrs, slog.Any(key, value))
				}

				logger.Info("http.request", append(
					logAttrs,
					slog.String("http.url", r.RequestURI),
					slog.String("http.method", r.Method),
					slog.Int("http.status", wrappedWriter.statusCode),
					slog.Duration("duration", time.Now().Sub(startedAt)),
					slog.String("ip", r.RemoteAddr),
				)...)
			}()

			r = r.WithContext(context.WithValue(r.Context(), "dise.logger", attributes))
			h.ServeHTTP(wrappedWriter, r)
		})
	}
}

func LogAttr(ctx context.Context, key string, value interface{}) {
	attributes, ok := ctx.Value("dise.logger").(map[string]interface{})
	if !ok || attributes == nil {
		return
	}
	attributes[key] = value
}

func LogError(ctx context.Context, err error) {
	LogAttr(ctx, "error", err)
}
