package account

import (
	"database/sql"
	"net/http"

	"github.com/go-chi/jwtauth"
	"github.com/go-chi/render"
	"github.com/mtdx/keyc/common"
)

// PurchasesResponse ...
type PurchasesResponse struct {
	Type      string `json:"type" validate:"nonzero"`
	Status    string `json:"status" validate:"nonzero"`
	Amount    uint32 `json:"amount" validate:"min=1"`
	UnitPrice uint32 `json:"unit_price" validate:"nonzero"`

	CreatedAt string `json:"created_at" validate:"nonzero"`
}

func (rd *PurchasesResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// PurchasesHandler rest route handler
func PurchasesHandler(w http.ResponseWriter, r *http.Request) {
	_, claims, _ := jwtauth.FromContext(r.Context())
	dbconn := r.Context().Value("DBCONN").(*sql.DB)
	rows, err := dbconn.Query(`SELECT type, status, failure_details, amount, created_at 
		FROM Purchases WHERE user_steam_id = $1`, claims["id"])
	if err != nil {
		render.Render(w, r, common.ErrInternalServer(err))
		return
	}
	defer rows.Close()
	list := []render.Renderer{}
	for rows.Next() {
		resp := &PurchasesResponse{}
		if err := rows.Scan(&resp.Type, &resp.Status, &resp.FailureDetails, &resp.Amount, &resp.CreatedAt); err != nil {
			render.Render(w, r, common.ErrInternalServer(err))
			return
		}
		list = append(list, resp)
	}
	if err := rows.Err(); err != nil {
		render.Render(w, r, common.ErrInternalServer(err))
		return
	}

	render.Status(r, http.StatusOK)
	if err := render.RenderList(w, r, list); err != nil {
		render.Render(w, r, common.ErrRender(err))
	}
}
