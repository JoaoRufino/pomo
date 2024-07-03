package rest

import (
	"bytes"
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/go-chi/chi/middleware"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func loggerHTTPMiddlewareDefault(logRequestBody bool, logDuration bool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			writerWrap := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			var response *bytes.Buffer
			if logRequestBody {
				response = new(bytes.Buffer)
				writerWrap.Tee(response)
			}

			next.ServeHTTP(w, r)

			fields := []zapcore.Field{
				zap.Int("status", writerWrap.Status()),
				zap.String("path", r.RequestURI),
				zap.String("method", r.Method),
				zap.String("package", "rest_server"),
			}

			if reqID := r.Context().Value(middleware.RequestIDKey); reqID != nil {
				fields = append(fields, zap.String("request-id", reqID.(string)))
			}

			if logRequestBody {
				if req, err := httputil.DumpRequest(r, true); err == nil {
					fields = append(fields, zap.ByteString("request", req))
				}
				fields = append(fields, zap.ByteString("response", response.Bytes()))
			}
			if logDuration {
				fields = append(fields, zap.Duration("duration", time.Since(start)))
			}
			zap.L().Debug("HTTP Request", fields...)
		})
	}
}
