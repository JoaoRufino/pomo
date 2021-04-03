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
		"logger.level":              "info",
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
		"server.host":           "",
		"server.port":           "8080",
		"server.socket":         "../../../tmp/test/pomo.sock",
		"server.datetimeformat": "2006-01-02 15:04",

		// Database Settings
		"database.username":              "postgres",
		"database.password":              "password",
		"database.host":                  "postgres",
		"database.port":                  5432,
		"database.database":              "gorestapi",
		"database.auto_create":           true,
		"database.retries":               5,
		"database.sleep_between_retries": "7s",
		"database.max_connections":       40,
		"database.log_queries":           false,
		"database.wipe_confirm":          false,
		"database.path":                  "../../../tmp/test/pomo.db",
	}, "."), nil)
}
