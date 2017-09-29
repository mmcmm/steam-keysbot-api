package config

import "os"

// Env config from env or dev defaults
func Env() map[string]string {
	m := map[string]string{
		"host":     os.Getenv("KEYC_DB_HOST"),
		"user":     os.Getenv("KEYC_DB_USER"),
		"password": os.Getenv("KEYC_DB_PASSWORD"),
	}

	if m["host"] == "" {
		m["host"] = "localhost"
	}
	if m["user"] == "" {
		m["user"] = "keyc"
	}
	if m["password"] == "" {
		m["password"] = "password"
	}

	return m
}
