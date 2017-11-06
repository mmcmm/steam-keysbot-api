package internal

import (
	"log"
)

// E ... errors
type E struct {
	Func      uint8  `json:"func" validate:"nonzero"`
	Message   uint8  `json:"message" validate:"nonzero"`
	CreatedAt string `json:"created_at" validate:"nonzero"`
}

func saveErr(e E) {
	_, err := dbconn.Exec("INSERT INTO errors (func, message) VALUES($1, $2)", e.Func, e.Message)
	if err != nil {
		log.Printf("Error saving error: %s", err.Error())
	}
}
