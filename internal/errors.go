package internal

import (
	"log"
)

// E ... errors
type E struct {
	Func    string `json:"source" validate:"nonzero"`
	Message string `json:"message" validate:"nonzero"`
}

// SaveErr ...
func SaveErr(e E) {
	_, err := dbconn.Exec("INSERT INTO errors (source, message) VALUES($1, $2)", e.Func, e.Message)
	if err != nil {
		log.Printf("Error saving error: %s", err.Error())
	}
}
