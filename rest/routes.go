package rest

import (
	"database/sql"
	"time"

	"github.com/mtdx/keyc/account"
	"github.com/mtdx/keyc/keys"
	"github.com/mtdx/keyc/openid/steamauth"
	"github.com/mtdx/keyc/vault"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/jwtauth"
	"github.com/go-chi/render"
	"github.com/mtdx/keyc/config"
)

var r *chi.Mux

func routes() {
	tokenAuth := jwtauth.New("HS256", []byte(config.JwtKey()), nil)
	r.Get("/login", steamauth.LoginHandler)

	r.Route("/api/v1", func(r chi.Router) {
		// Protected routes
		r.Group(func(r chi.Router) {
			r.Use(jwtauth.Verifier(tokenAuth))
			r.Use(jwtauth.Authenticator)

			r.Get("/account", account.InfoHandler)

			r.Get("/keys-transactions", keys.TransactionsHandler)

			r.Get("/withdrawals", vault.WithdrawalsHandler)
			r.Post("/withdrawals", vault.WithdrawalRequestHandler)
		})

		r.Put("/keys-transactions/{tradeofferID:[1-9]+}", keys.TransactionUpdateHandler)
		r.Post("/keys-transactions", keys.TransactionCreateHandler)
	})
}

// Router create chi router & add the routes
func Router(dbconn *sql.DB) *chi.Mux {
	r = chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Use(middleware.DefaultCompress)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Use(middleware.WithValue("DBCONN", dbconn))

	routes()

	return r
}
