package test

import (
	"github.com/joaorufino/pomo/pkg/conf"
	"github.com/spf13/viper"
)

// ConfFromDefaults loads the default config using Viper
func ConfFromDefaults() *conf.Config {
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
	viper.SetDefault("server.unix.socket", "../../test/pomo.sock")
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
	viper.SetDefault("database.path", "../../test/pomo.db")

	var config conf.Config
	viper.Unmarshal(&config)
	return &config
}
