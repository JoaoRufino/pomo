package conf

import (
	"encoding/json"
	"net/http"
)

var (
	// Executable name
	Executable = K.String("server.name")
	// Version
	Version = K.String("server.version")
)

// GetVersion returns version as a simple json
func GetVersion() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		v := struct {
			Version string `json:"version"`
		}{
			Version: Version,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(v)
	}
}
