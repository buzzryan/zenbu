package httputil

import (
	"io"
	"log/slog"
	"net/http"
	"runtime/debug"
	"slices"
	"time"

	"github.com/google/uuid"

	"github.com/buzzryan/zenbu/internal/logutil"
)

var contentTypesAllowedToLog = []string{
	MIMETypeApplicationForm, MIMETypeApplicationJSON, MIMETypeTextPlain,
}

type reqLogger struct {
	http.Handler
}

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

	if slices.Contains(contentTypesAllowedToLog, req.Header.Get(ContentType)) {
		body, err = io.ReadAll(req.Body)
		if err != nil {
			logutil.From(req.Context()).With("err", err).Error(
				"failed to read request body",
			)
		}
	}

	writer := &responseWriter{ResponseWriter: w}
	rl.Handler.ServeHTTP(writer, req)

	logutil.From(req.Context()).With(
		slog.Any("header", req.Header),
		slog.Any("body", string(body)),
		slog.Any("status_code", writer.statusCode),
	).Info("request")
}

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
