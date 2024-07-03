package conf

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/viper"
)

const ConfDefaultPath = "configs/defaults.yaml"

// LoadConfig loads the configuration from the given file and merges it with defaults
func LoadConfig(configFile string) (*Config, error) {
	// Load default configuration
	if err := loadDefaultConfig(); err != nil {
		return nil, err
	}

	// Load user configuration
	if configFile != "" {
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

		if err := viper.MergeInConfig(); err != nil {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("unable to decode into struct: %w", err)
	}

	return &config, nil
}

// loadDefaultConfig loads the default configuration from a YAML file
func loadDefaultConfig() error {
	viper.SetConfigFile(ConfDefaultPath)
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("error reading default config file: %w", err)
	}
	return nil
}
