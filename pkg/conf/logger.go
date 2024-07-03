package conf

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// InitLogger loads a global logger based on Viper configuration
func InitLogger() {
	logConfig := zap.NewDevelopmentConfig()
	logConfig.Sampling = nil

	// Log Level
	var logLevel zapcore.Level
	if err := logLevel.Set(viper.GetString("logger.level")); err != nil {
		zap.S().Fatalw("Could not determine logger.level", "error", err)
	}
	logConfig.Level.SetLevel(logLevel)

	// Handle different logger encodings
	loggerEncoding := viper.GetString("logger.encoding")
	logConfig.Encoding = loggerEncoding
	// Enable Color
	if viper.GetBool("logger.color") {
		logConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}
	logConfig.DisableStacktrace = viper.GetBool("logger.disable_stacktrace")
	// Use sane timestamp when logging to console
	if logConfig.Encoding == "console" {
		logConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	}

	// JSON Fields
	logConfig.EncoderConfig.MessageKey = "msg"
	logConfig.EncoderConfig.LevelKey = "level"
	logConfig.EncoderConfig.CallerKey = "caller"

	// Settings
	logConfig.Development = viper.GetBool("logger.dev_mode")
	logConfig.DisableCaller = viper.GetBool("logger.disable_caller")

	// Build the logger
	globalLogger, _ := logConfig.Build()
	zap.ReplaceGlobals(globalLogger)
}
