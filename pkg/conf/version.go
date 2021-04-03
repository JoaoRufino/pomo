package conf

import (
	"encoding/json"
	"net/http"

	"github.com/knadh/koanf"
)

// GetVersion returns version as a simple json
func GetVersion(K *koanf.Koanf) http.HandlerFunc {
	// Version
	version := K.String("server.version")

	return func(w http.ResponseWriter, r *http.Request) {
		v := struct {
			Version string `json:"version"`
		}{
			Version: version,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(v)
	}
}
