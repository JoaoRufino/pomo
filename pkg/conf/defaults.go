package conf

import (
	"fmt"
	"os"
	"os/user"
	"path"
	"path/filepath"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/providers/file"
)

// File loads configuration from a file
// https://github.com/knadh/koanf/blob/master/examples/read-file/main.go
func ConfFromFile(k *koanf.Koanf, configFile string) error {
	ext := filepath.Ext(configFile)
	f := file.Provider(configFile)
	switch ext {
	case ".yaml":
		return k.Load(f, yaml.Parser())
	case ".json":
		return k.Load(f, json.Parser())
	}
	return fmt.Errorf("unknown config extension %s", ext)
}

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
		"server.type":           "unix",
		"server.host":           "",
		"server.port":           "8080",
		"server.socket":         defaultConfigPath() + "/pomo.sock",
		"server.datetimeformat": "2006-01-02 15:04",

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
		"database.sslcert":               "",
		"database.sslkey":                "",
		"database.sslrootcert":           "",
		"database.retries":               5,
		"database.sleep_between_retries": "7s",
		"database.max_connections":       40,
		"database.log_queries":           false,
		"database.wipe_confirm":          false,
		"database.path":                  defaultConfigPath() + "/pomo.db",
	}, "."), nil)
}

func defaultConfigPath() string {
	u, err := user.Current()
	if err != nil {
		fmt.Printf("Error:\n%s\n", err)
		os.Exit(1)
	}
	return path.Join(u.HomeDir, "/.pomo")
}
