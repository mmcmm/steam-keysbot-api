package config

import "os"

// JwtKey private key from env or insecure default
func JwtKey() string {
	key := os.Getenv("KEYC_JWT_KEY")
	if key == "" {
		key = "InsecurePrivateDefaultKey" // default insecure private key
	}
	return key
}
