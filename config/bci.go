package config

import "os"

// BlockchainInfoKey private API key from env
func BlockchainInfoKey() string {
	key := os.Getenv("KEYC_BCI_KEY")
	if key == "" {
		key = "BlockchainInfoKey"
	}
	return key
}
