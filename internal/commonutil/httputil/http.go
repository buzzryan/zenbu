package httputil

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/buzzryan/zenbu/internal/commonutil/logutil"
)

// ParseJSONBody is a helper function to parse JSON request body.
func ParseJSONBody(r *http.Request, v interface{}) error {
	if !strings.Contains(r.Header.Get(ContentType), MIMETypeApplicationJSON) {
		return ErrInvalidContentType
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("failed to read body when parseBody: %w", err)
	}

	if err = json.Unmarshal(body, v); err != nil {
		return fmt.Errorf("%w: %w", ErrInvalidJSONBody, err)
	}

	return nil
}

// HandleParseJSONBodyError is a helper function to handle error when parsing JSON request body.
func HandleParseJSONBodyError(ctx context.Context, w http.ResponseWriter, err error) error {
	if errors.Is(err, ErrInvalidContentType) {
		return ResponseError(
			w, http.StatusBadRequest, CodeInvalidContentType, err.Error(),
		)
	}
	if errors.Is(err, ErrInvalidJSONBody) {
		return ResponseError(
			w, http.StatusBadRequest, CodeInvalidJSONBody, err.Error(),
		)
	}

	logutil.From(ctx).With(slog.Any("err", err)).Error("failed to parse json body")
	return ResponseError(w, http.StatusInternalServerError, 0, "internal server error")
}

// ResponseJSON is a helper function to response JSON with status code.
func ResponseJSON(w http.ResponseWriter, code int, v interface{}) error {
	w.Header().Set(ContentType, MIMETypeApplicationJSON)
	w.WriteHeader(code)

	body, err := json.Marshal(v)
	if err != nil {
		return fmt.Errorf("failed to marshal response: %w", err)
	}

	contentLength, err := w.Write(body)
	if err != nil {
		return fmt.Errorf("failed to write response: %w", err)
	}
	w.Header().Set("Content-Length", strconv.Itoa(contentLength))

	return nil
}

// errorResponse is unified DTO for handling error
type errorResponse struct {
	ErrorMessage string `json:"error_message"`
	// ErrorCode is detailed error code. Sometimes HTTP status code is not enough to represent error.
	ErrorCode int `json:"error_code"`
}

// ResponseError is a helper function to response error with error code and message.
// It helps to response error in unified format.
func ResponseError(w http.ResponseWriter, statusCode, errCode int, errMsg string) error {
	return ResponseJSON(w, statusCode, errorResponse{ErrorCode: errCode, ErrorMessage: errMsg})
}

// HandlerFuncWithErr is similar with http.HandleFunc, but it returns unhandled error.
// In most cases, HandlerFuncWithErr should handle error, response to client and not return error.
type HandlerFuncWithErr func(w http.ResponseWriter, r *http.Request) error

// RegisterHandler is a helper function to register handler with error handling.
// It covers error such as fail to write response, etc.
func RegisterHandler(mux *http.ServeMux, method string, pattern string, handler HandlerFuncWithErr) {
	mux.HandleFunc(method+" "+pattern, func(w http.ResponseWriter, r *http.Request) {
		if err := handler(w, r); err != nil {
			logutil.From(r.Context()).Error("failed to handle request", slog.Any("err", err))
			w.WriteHeader(http.StatusInternalServerError)
		}
	})
}

// GetBearerToken is a helper function to get bearer token from Authorization header.
// If Authorization header is not found or invalid, it returns error.
func GetBearerToken(req *http.Request) (string, error) {
	authHeader := req.Header.Get(Authorization)
	if authHeader == "" {
		return "", errors.New("authorization header not found")
	}

	authParts := strings.Split(authHeader, " ")
	if len(authParts) != 2 || authParts[0] != Bearer {
		return "", errors.New("invalid bearer authorization header")
	}

	return authParts[1], nil
}
