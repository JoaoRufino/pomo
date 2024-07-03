package conf

// Config represents the application's configuration
type Config struct {
	Logger   LoggerConfig
	Pidfile  string
	Server   ServerConfig
	Database DatabaseConfig
}

// LoggerConfig represents the logger's configuration
type LoggerConfig struct {
	Level             string
	Encoding          string
	Color             bool
	DevMode           bool
	DisableCaller     bool
	DisableStacktrace bool
}

// ServerConfig represents the server's configuration
type ServerConfig struct {
	Name           string
	Version        string
	Type           string
	RestHost       string
	RestPort       string
	UnixSocket     string
	DatetimeFormat string
	LogRequests    bool
}

// DatabaseConfig represents the database's configuration
type DatabaseConfig struct {
	Username            string
	Password            string
	Host                string
	Port                int
	Database            string
	Type                string
	AutoCreate          bool
	SearchPath          string
	SSLMode             string
	Retries             int
	SleepBetweenRetries string
	MaxConnections      int
	LogQueries          bool
	Path                string
}
