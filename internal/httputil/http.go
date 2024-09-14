package httputil

import (
	"net/http"
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
