package main

import (
	"net/http"

	_ "github.com/mtdx/keyc/settings"

	"github.com/mtdx/keyc/db"
	"github.com/mtdx/keyc/rest"
)

func main() {
	dbconn := db.Open()
	db.RunMigrations(dbconn)
	defer dbconn.Close()

	r := rest.Router(dbconn)
	http.ListenAndServe(":8080", r)
}
