package conf

import (
	"fmt"
	"os"
	"os/user"
	"path"
	"path/filepath"

	"github.com/spf13/viper"
)

// LoadConfig loads the configuration from the given file
func LoadConfig(configFile string) (*Config, error) {
	viper.SetConfigFile(configFile)
	ext := filepath.Ext(configFile)
	switch ext {
	case ".yaml", ".yml":
		viper.SetConfigType("yaml")
	case ".json":
		viper.SetConfigType("json")
	default:
		return nil, fmt.Errorf("unknown config extension %s", ext)
	}

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("unable to decode into struct: %w", err)
	}

	return &config, nil
}

// LoadDefaultConfig loads the default configuration
func LoadDefaultConfig() *Config {
	viper.SetDefault("logger.level", "debug")
	viper.SetDefault("logger.encoding", "console")
	viper.SetDefault("logger.color", true)
	viper.SetDefault("logger.dev_mode", true)
	viper.SetDefault("logger.disable_caller", false)
	viper.SetDefault("logger.disable_stacktrace", true)

	viper.SetDefault("pidfile", "")

	viper.SetDefault("server.name", "pomo")
	viper.SetDefault("server.version", "debug")
	viper.SetDefault("server.type", "rest")
	viper.SetDefault("server.rest.host", "")
	viper.SetDefault("server.rest.port", "8080")
	viper.SetDefault("server.unix.socket", defaultConfigPath()+"/pomo.sock")
	viper.SetDefault("server.datetimeformat", "2006-01-02 15:04")
	viper.SetDefault("server.log_requests", true)

	viper.SetDefault("database.username", "postgres")
	viper.SetDefault("database.password", "password")
	viper.SetDefault("database.host", "postgres")
	viper.SetDefault("database.port", 5432)
	viper.SetDefault("database.database", "pomo")
	viper.SetDefault("database.type", "sqlite")
	viper.SetDefault("database.auto_create", true)
	viper.SetDefault("database.search_path", "")
	viper.SetDefault("database.sslmode", "disable")
	viper.SetDefault("database.retries", 5)
	viper.SetDefault("database.sleep_between_retries", "7s")
	viper.SetDefault("database.max_connections", 40)
	viper.SetDefault("database.log_queries", true)
	viper.SetDefault("database.path", defaultConfigPath()+"/pomo.db")

	var config Config
	viper.Unmarshal(&config)
	return &config
}

func defaultConfigPath() string {
	u, err := user.Current()
	if err != nil {
		fmt.Printf("Error:\n%s\n", err)
		os.Exit(1)
	}
	return path.Join(u.HomeDir, "/.pomo")
}
