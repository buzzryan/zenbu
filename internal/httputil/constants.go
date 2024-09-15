package httputil

import "errors"

/* Header Keys */
const (
	ContentType   = "Content-Type"
	CorrelationID = "Correlation-Id"
)

/* MIME Types */
const (
	MIMETypeApplicationJSON = "application/json"
	MIMETypeApplicationForm = "application/x-www-form-urlencoded"
	MIMETypeTextPlain       = "text/plain"
)

/* Common Error Codes. 1000 - 1999 is reserved for general errors. */
const (
	CodeInvalidRequestParams = 1000 // Invalid request parameters (body, query, path param, etc.)
	CodeInvalidContentType   = 1001
	CodeInvalidJSONBody      = 1002
)

/* Common Errors */
var (
	ErrInvalidContentType = errors.New("invalid content type")
	ErrInvalidJSONBody    = errors.New("invalid JSON body")
)
