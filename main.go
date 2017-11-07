package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/mtdx/keyc/db"
	"github.com/mtdx/keyc/internal"
	"github.com/mtdx/keyc/rest"
)

func main() {
	dbconn := db.Open()
	db.RunMigrations(dbconn)
	defer dbconn.Close()

	if err := internal.InitLiveBtc(); err != nil {
		fmt.Fprintf(os.Stderr, "Unable to start BTC live price: %v\n", err)
		return
	}

	r := rest.Router(dbconn)
	http.ListenAndServe(":8080", r)
}
