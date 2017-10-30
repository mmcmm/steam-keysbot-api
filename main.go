package main

import (
	"net/http"

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
