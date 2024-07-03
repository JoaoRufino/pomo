package test

import (
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/providers/confmap"
)

// Defaults loads the default config for the
// https://github.com/knadh/koanf/blob/master/examples/default-values/main.go
func ConfFromDefaults(k *koanf.Koanf) error {
	return k.Load(confmap.Provider(map[string]interface{}{
		// Logger Defaults
		"logger.level":              "debug",
		"logger.encoding":           "console",
		"logger.color":              true,
		"logger.dev_mode":           true,
		"logger.disable_caller":     false,
		"logger.disable_stacktrace": true,

		// Pidfile
		"pidfile": "",

		// Server Configuration
		"server.name":           "pomo",
		"server.version":        "debug",
		"server.type":           "rest",
		"server.rest.host":      "",
		"server.rest.port":      "8080",
		"server.unix.socket":    "../../test/pomo.sock",
		"server.datetimeformat": "2006-01-02 15:04",
		"server.log_requests":   true,

		// Database Settings
		"database.username":              "postgres",
		"database.password":              "password",
		"database.host":                  "postgres",
		"database.port":                  5432,
		"database.database":              "pomo",
		"database.type":                  "sqlite",
		"database.auto_create":           true,
		"database.search_path":           "",
		"database.sslmode":               "disable",
		"database.retries":               5,
		"database.sleep_between_retries": "7s",
		"database.max_connections":       40,
		"database.log_queries":           true,
		"database.path":                  "../../test/pomo.db",
	}, "."), nil)
}
