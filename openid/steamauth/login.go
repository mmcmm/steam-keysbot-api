package steamauth

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/jwtauth"
	"github.com/go-chi/render"
	"github.com/mtdx/keyc/common"
	"github.com/mtdx/keyc/config"
)

type authResponse struct {
	Token string `json:"token"`
}

func (rd *authResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// LoginHandler ...
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	opId := newOpenId(r)
	switch opId.mode() {
	case "":
		http.Redirect(w, r, opId.authUrl(), 301)
		return
	case "cancel":
		render.Render(w, r, common.ErrInvalidRequest(errors.New("Authorization cancelled")))
		return
	default:
		steamId, err := opId.validateAndGetId()
		if err != nil {
			render.Render(w, r, common.ErrInvalidRequest(err))
			return
		}
		player, err := getPlayerSummaries(steamId, config.SteamAPIKey())
		if err != nil {
			render.Render(w, r, common.ErrInternalServer(err))
			return
		}

		dbconn := r.Context().Value("DBCONN").(*sql.DB)
		_, err = dbconn.Exec(`INSERT INTO users (steam_id, display_name, avatar) VALUES ($1, $2, $3) 
		ON CONFLICT (steam_id) DO UPDATE SET display_name = $2, avatar = $3`,
			steamId, player.PersonaName, player.Avatar)
		if err != nil {
			render.Render(w, r, common.ErrInternalServer(err))
			return
		}

		tokenAuth := jwtauth.New("HS256", []byte(config.JwtKey()), nil)
		exp := time.Now().Add(time.Hour * time.Duration(12)).Unix()
		claims := jwtauth.Claims{"id": steamId, "exp": exp}
		_, tokenString, _ := tokenAuth.Encode(claims)
		resp := &authResponse{Token: tokenString}

		render.Status(r, http.StatusOK)
		render.Render(w, r, resp)
	}
}
