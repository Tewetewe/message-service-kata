package response

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type (
	// HTTPError represents an error that occurred while handling a request.
	HTTPError struct {
		Code     int     `json:"status"`
		Message  Message `json:"message"`
		Internal error   `json:"-"` // Stores the error returned by an external dependency
	}
)

// HTTP Errors
//
//nolint:lll
var (
	ErrBadRequest          = NewHTTPError(http.StatusBadRequest, DefaultErrorMessage)                         // HTTP 400 Bad Request.
	ErrNotFound            = NewHTTPError(http.StatusNotFound, ResponseMessageNotFound)                       // HTTP 404 Not Found.
	ErrMethodNotAllowed    = NewHTTPError(http.StatusMethodNotAllowed, ResponseMessageMethodNotAllowed)       // HTTP 405 Method Not Allowed.
	ErrUnprocessableEntity = NewHTTPError(http.StatusUnprocessableEntity, ResponseMessageUnprocessableEntity) // HTTP 422 Unprocessable Entity.
	ErrInternalServerError = NewHTTPError(http.StatusInternalServerError, ResponseMessageInternalServerError) // HTTP 500 Internal Server Error.
	ErrServiceUnavailable  = NewHTTPError(http.StatusServiceUnavailable, ResponseMessageServiceUnavailable)   // HTTP 503 Service Unavailable.
)

// NewHTTPError creates a new HTTPError instance.
func NewHTTPError(code int, message Message) *HTTPError {
	he := &HTTPError{Code: code, Message: map[string]string{
		"en": http.StatusText(code),
	}}
	if message != nil {
		he.Message = message
	}
	return he
}

// Error makes it compatible with `error` interface.
func (he *HTTPError) Error() string {
	if he.Internal == nil {
		return fmt.Sprintf("HTTP %d error occurred; %s", he.Code, he.Message)
	}

	return fmt.Sprintf("HTTP %d error occurred; %s. Details: %v", he.Code, he.Message, he.Internal)
}

// WithInternal returns clone of HTTPError with err set to HTTPError.Internal field.
func (he *HTTPError) WithInternal(err error) *HTTPError {
	return &HTTPError{
		Code:     he.Code,
		Message:  he.Message,
		Internal: err,
	}
}

// Unwrap satisfies the Go 1.13 error wrapper interface.
func (he *HTTPError) Unwrap() error {
	return he.Internal
}

// DefaultHTTPErrorHandler is the default HTTP error handler. It sends a JSON response with status code.
//
// NOTE: In case errors happens in middleware call-chain that is returning from handler (which did not return an error).
// When handler has already sent response (ala c.JSON()) and there is error in middleware that is returning from
// handler. Then the error that global error handler received will be ignored because we have already "committed"
// the response and status code header has been sent to the client.
//
//nolint:gocyclo
func DefaultHTTPErrorHandler(err error, c echo.Context) {
	if c.Response().Committed {
		return
	}

	var he *HTTPError
	// Check any possibility kind of errors
	switch {
	case errors.Is(err, echo.ErrNotFound):
		he = ErrNotFound
	case errors.Is(err, echo.ErrBadRequest):
		he = ErrBadRequest
	case errors.Is(err, echo.ErrUnprocessableEntity):
		he = ErrUnprocessableEntity
	case errors.Is(err, echo.ErrMethodNotAllowed):
		he = ErrMethodNotAllowed
	case errors.Is(err, echo.ErrInternalServerError):
		he = ErrInternalServerError
	case errors.Is(err, echo.ErrServiceUnavailable):
		he = ErrServiceUnavailable
	default:
		if _he, ok := err.(*HTTPError); ok {
			he = _he
		} else {
			he = &HTTPError{
				Code: http.StatusInternalServerError,
				Message: map[string]string{
					"en": http.StatusText(http.StatusInternalServerError),
				},
				Internal: err,
			}
		}
	}

	response := HTTPResponse{
		Status:  he.Code,
		Message: he.Message,
	}
	if he.Internal != nil {
		response.Errors = []Error{
			{
				Code:            he.Code,
				UserMessage:     he.Message["en"],
				InternalMessage: fmt.Sprintf("%+v", he.Internal),
			},
		}
	}
	// Overwrite error response if debug enabled
	if c.Echo().Debug {
		response.Errors = []Error{
			{
				Code:            he.Code,
				UserMessage:     he.Message["en"],
				InternalMessage: fmt.Sprintf("%+v", err),
			},
		}
	}

	// Send response
	if c.Request().Method == http.MethodHead {
		// Issue echo#608
		err = c.NoContent(response.Status)
	} else {
		err = c.JSON(response.Status, response)
	}

	if err != nil {
		c.Logger().Error(err)
	}
}
