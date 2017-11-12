package steam

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"
	"github.com/go-chi/render"
	"github.com/mtdx/keyc/common"
	"github.com/mtdx/keyc/config"
)

// TradeoffersHandler rest route handler GET /
func TradeoffersHandler(w http.ResponseWriter, r *http.Request) {
	_, claims, _ := jwtauth.FromContext(r.Context())
	dbconn := r.Context().Value("DBCONN").(*sql.DB)
	tradeoffersresp, err := findAllTradeoffers(dbconn, claims["id"])
	common.ValidateRenderResults(w, r, tradeoffersresp, err)
}

// TradeoffersCreateHandler POST /tradeoffers
func TradeoffersCreateHandler(w http.ResponseWriter, r *http.Request) {
	queryValues := r.URL.Query()
	dbconn := r.Context().Value("DBCONN").(*sql.DB)
	if queryValues.Get("key") != config.SteamBotsAPIKey() || !IsOurSteamBot(dbconn, r.RemoteAddr) {
		render.Render(w, r, common.ErrInvalidRequest(errors.New("Unauthorized")))
		return
	}

	tradeoffer := &TradeoffersRequest{}
	if err := render.Bind(r, tradeoffer); err != nil {
		render.Render(w, r, common.ErrInvalidRequest(err))
		return
	}

	if err := saveTradeoffer(dbconn, tradeoffer); err != nil {
		render.Render(w, r, common.ErrInternalServer(err))
		return
	}
	render.Render(w, r, common.SuccessCreatedResponse("Tradeoffer has been created"))
}

// TradeoffersUpdateHandler PUT /tradeoffers/:id
func TradeoffersUpdateHandler(w http.ResponseWriter, r *http.Request) {
	queryValues := r.URL.Query()
	dbconn := r.Context().Value("DBCONN").(*sql.DB)
	if queryValues.Get("key") != config.SteamBotsAPIKey() || !IsOurSteamBot(dbconn, r.RemoteAddr) {
		render.Render(w, r, common.ErrInvalidRequest(errors.New("Unauthorized")))
		return
	}

	tradeoffer := &TradeoffersUpdateRequest{}
	if err := render.Bind(r, tradeoffer); err != nil {
		render.Render(w, r, common.ErrInvalidRequest(err))
		return
	}

	if err := updateStatus(dbconn, tradeoffer, chi.URLParam(r, "tradeofferID")); err != nil {
		render.Render(w, r, common.ErrInternalServer(err))
		return
	}
	render.Render(w, r, common.SuccessCreatedResponse("Tradeoffer has been updated"))
}
