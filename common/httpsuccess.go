package common

import (
	"net/http"

	"github.com/go-chi/render"
)

// SuccessResponse ...
type SuccessResponse struct {
	HTTPStatusCode int `json:"-"` // http response status code

	StatusText  string `json:"status"`          // user-level status message
	AppCode     int64  `json:"code,omitempty"`  // application-specific code
	SuccessText string `json:"error,omitempty"` // application-level message, for debugging
}

// Render ...
func (s *SuccessResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, s.HTTPStatusCode)
	return nil
}

// SuccessCreatedResponse ...
func SuccessCreatedResponse(message string) render.Renderer {
	return &SuccessResponse{
		HTTPStatusCode: 201,
		StatusText:     message,
	}
}
