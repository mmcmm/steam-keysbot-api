package common

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/mtdx/keyc/validator"
)

// ValidateRenderResults ... validate & renders `multiple` results
func ValidateRenderResults(w http.ResponseWriter, r *http.Request, resp []render.Renderer, err error) {
	if err != nil {
		render.Render(w, r, ErrInternalServer(err))
		return
	}
	for _, entry := range resp {
		if err := validator.Validate(entry); err != nil {
			render.Render(w, r, ErrInternalServer(err))
			return
		}
	}
	render.Status(r, http.StatusOK)
	if err := render.RenderList(w, r, resp); err != nil {
		render.Render(w, r, ErrRender(err))
	}
}

// ValidateRenderResult ... validate & renders `single` results
func ValidateRenderResult(w http.ResponseWriter, r *http.Request, resp render.Renderer, err error) {
	if err != nil {
		render.Render(w, r, ErrInternalServer(err))
		return
	}
	if err := validator.Validate(resp); err != nil {
		render.Render(w, r, ErrInternalServer(err))
		return
	}
	render.Status(r, http.StatusOK)
	if err := render.Render(w, r, resp); err != nil {
		render.Render(w, r, ErrRender(err))
	}
}
