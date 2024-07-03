package conf

import (
	"encoding/json"
	"net/http"

	"github.com/spf13/viper"
)

// GetVersion returns version as a simple json
func GetVersion() http.HandlerFunc {
	// Version
	version := viper.GetString("server.version")

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
