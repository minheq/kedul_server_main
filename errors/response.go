package errors

import (
	"net/http"

	"github.com/go-chi/render"
)

// ErrResponse represents standardized error response
type ErrResponse struct {
	HTTPStatusCode int `json:"-"` // http response status code

	Message string `json:"message,omitempty"` // human readable message
	Code    string `json:"code,omitempty"`    // human readable message
	DocURL  string `json:"doc_url,omitempty"` // human readable message
}

// Render error with HTTP status code
func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

// HTTPStatusCode extracts http status code from error
func HTTPStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}

	if e, ok := err.(*Error); ok && e.Kind != 0 {
		switch e.Kind {
		case KindInvalid:
			return http.StatusBadRequest
		case KindUnauthorized:
			return http.StatusUnauthorized
		case KindNotFound:
			return http.StatusNotFound
		case KindUnexpected:
			return http.StatusInternalServerError
		default:
			return http.StatusInternalServerError
		}
	} else if ok && e.Err != nil {
		return HTTPStatusCode(e.Err)
	}

	return http.StatusInternalServerError
}

// NewErrResponse are converted from Error from internal errors package
func NewErrResponse(err error) render.Renderer {
	return &ErrResponse{
		HTTPStatusCode: HTTPStatusCode(err),
		Message:        ErrorMessage(err),
	}
}
