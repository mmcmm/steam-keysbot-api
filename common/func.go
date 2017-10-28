package common

import (
	"database/sql"
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

// Transact execute a db transaction
func Transact(db *sql.DB, txFunc func(*sql.Tx) error) (err error) {
	tx, err := db.Begin()
	if err != nil {
		return
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p) // re-throw panic after Rollback
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()
	err = txFunc(tx)
	return err
}
