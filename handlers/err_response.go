package handlers

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/minheq/kedul_server_main/errors"
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

	if e, ok := err.(*errors.Error); ok && e.Kind != "" {
		switch e.Kind {
		case errors.KindInvalid:
			return http.StatusBadRequest
		case errors.KindUnauthorized:
			return http.StatusUnauthorized
		case errors.KindNotFound:
			return http.StatusNotFound
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
		Message:        errors.ErrorMessage(err),
	}
}
