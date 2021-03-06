package common

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/mtdx/keyc/internal"
)

// ErrResponse ...
type ErrResponse struct {
	Err            error `json:"-"`    // low-level runtime error
	HTTPStatusCode int   `json:"code"` // http response status code

	StatusText string `json:"status"`            // user-level status message
	AppCode    int64  `json:"appcode,omitempty"` // application-specific error code
	ErrorText  string `json:"error,omitempty"`   // application-level error message, for debugging
}

// Render ...
func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

// ErrInvalidRequest ...
func ErrInvalidRequest(err error) render.Renderer {
	internal.SaveErr(internal.E{Func: "400", Message: err.Error()})
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 400,
		StatusText:     "Invalid request",
	}
}

// ErrInternalServer ...
func ErrInternalServer(err error) render.Renderer {
	internal.SaveErr(internal.E{Func: "500", Message: err.Error()})
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 500,
		StatusText:     "Internal server error",
	}
}

// ErrRender ...
func ErrRender(err error) render.Renderer {
	internal.SaveErr(internal.E{Func: "422", Message: err.Error()})
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 422,
		StatusText:     "Error rendering response",
	}
}

// ErrNotFound ...
var ErrNotFound = &ErrResponse{HTTPStatusCode: 404, StatusText: "Resource not found"}
