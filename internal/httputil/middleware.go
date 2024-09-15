package httputil

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"
	"runtime/debug"
	"slices"
	"time"

	"github.com/google/uuid"

	"github.com/buzzryan/zenbu/internal/logutil"
)

// WithGlobalMiddlewares wraps http.Handler with global(common) middlewares.
func WithGlobalMiddlewares(handler http.Handler) http.Handler {
	return &tracer{
		Handler: &reqLogger{
			Handler: &recovery{
				Handler: handler,
			},
		},
	}
}

type reqLogger struct {
	http.Handler
}

// loggableContentTypes is a list of content types that can be logged.
// It is recommended to log only text-based content types. (e.g. application/json, text/plain)
// Media types like a File should not be logged. It can be issues with performance and security.
var loggableContentTypes = []string{
	MIMETypeApplicationForm, MIMETypeApplicationJSON, MIMETypeTextPlain,
}

// responseWriter is middleware for log request and response.
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *responseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func (rl *reqLogger) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var (
		body []byte
		err  error
	)

	// It logs request body if it is loggable.
	if slices.Contains(loggableContentTypes, req.Header.Get(ContentType)) {
		body, err = io.ReadAll(req.Body)
		if err != nil {
			logutil.From(req.Context()).With("err", err).Error(
				"failed to read request body",
			)
		}
		// body is read and closed, so it should be set again to req.Body.
		req.Body = io.NopCloser(bytes.NewBuffer(body))
	}

	writer := &responseWriter{ResponseWriter: w}
	rl.Handler.ServeHTTP(writer, req)

	logutil.From(req.Context()).With(
		slog.Any("header", req.Header),
		slog.Any("body", string(body)),
		slog.Any("status_code", writer.statusCode),
	).Info("request")
}

// tracer is middleware for tracing request.
// It injects logger to context with correlation ID and request scoped information for tracing.
// Correlation ID will be created if request header doesn't include that.
type tracer struct {
	http.Handler
}

func (t *tracer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	correlationID := req.Header.Get(CorrelationID)
	if correlationID == "" {
		correlationID = uuid.NewString()
	}

	t.Handler.ServeHTTP(w, req.WithContext(
		logutil.ContextWithLogger(
			req.Context(),
			slog.Default().With(
				slog.String("correlation_id", correlationID),
				slog.Time("request_time", time.Now()),
				slog.String("method", req.Method),
				slog.String("url", req.URL.String()),
			),
		),
	))
}

// recovery is middleware for recover from unexpected panic so that prevents the server from fault.
type recovery struct {
	http.Handler
}

func (t *recovery) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			logutil.From(req.Context()).With(
				slog.Any("stack_trace", string(debug.Stack())),
			).Error("unexpected panic has occurred")
		}
	}()
	t.Handler.ServeHTTP(w, req)
}
