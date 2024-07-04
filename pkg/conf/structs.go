package conf

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/fatih/color"
)

// Config represents the application's configuration
type Config struct {
	Logger   LoggerConfig
	Pidfile  string
	Server   ServerConfig
	Database DatabaseConfig
	Runner   RunnerConfig
	CORS     CORSConfig
	Client   ClientConfig
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
	Name            string
	Version         string
	Type            string
	RestHost        string
	RestPort        string
	RestPath        string
	UnixSocket      string
	DatetimeFormat  string
	LogRequests     bool
	LogRequestsBody bool
	LogDuration     bool
	TLS             bool
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

type ClientConfig struct {
	Type     string
	HostType string
	HostPort string
	Host     string
}

// CORSConfig represents the CORS configuration
type CORSConfig struct {
	AllowedOrigins   []string
	AllowedMethods   []string
	AllowedHeaders   []string
	AllowCredentials bool
	MaxAge           int
}

// RunnerConfig represents user preferences
type RunnerConfig struct {
	Colors      *ColorMap `json:"colors"`
	DateTimeFmt string    `json:"dateTimeFmt"`
	BasePath    string    `json:"basePath"`
	DBPath      string    `json:"dbPath"`
	SocketPath  string    `json:"socketPath"`
	IconPath    string    `json:"iconPath"`
}

// ColorMap holds the color mapping
type ColorMap struct {
	colors map[string]*color.Color
	tags   map[string]string
}

func (c *ColorMap) Get(name string) *color.Color {
	if color, ok := c.colors[name]; ok {
		return color
	}
	return nil
}

func (c *ColorMap) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.tags)
}

func (c *ColorMap) UnmarshalJSON(raw []byte) error {
	lookup := map[string]*color.Color{
		"black":     color.New(color.FgBlack),
		"hiblack":   color.New(color.FgHiBlack),
		"blue":      color.New(color.FgBlue),
		"hiblue":    color.New(color.FgHiBlue),
		"cyan":      color.New(color.FgCyan),
		"hicyan":    color.New(color.FgHiCyan),
		"green":     color.New(color.FgGreen),
		"higreen":   color.New(color.FgHiGreen),
		"magenta":   color.New(color.FgMagenta),
		"himagenta": color.New(color.FgHiMagenta),
		"red":       color.New(color.FgRed),
		"hired":     color.New(color.FgHiRed),
		"white":     color.New(color.FgWhite),
		"hiwrite":   color.New(color.FgHiWhite),
		"yellow":    color.New(color.FgYellow),
		"hiyellow":  color.New(color.FgHiYellow),
	}
	cm := &ColorMap{
		colors: map[string]*color.Color{},
		tags:   map[string]string{},
	}
	err := json.Unmarshal(raw, &cm.tags)
	if err != nil {
		return err
	}
	for tag, colorName := range cm.tags {
		if color, ok := lookup[colorName]; ok {
			cm.colors[tag] = color
		}
	}
	*c = *cm
	return nil
}

// Print outputs the configuration details to the console
func (c *Config) Print() {
	fmt.Println("Configuration:")
	c.printStruct(reflect.ValueOf(c).Elem(), "")
}

func (c *Config) printStruct(v reflect.Value, indent string) {
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := v.Type().Field(i)
		fieldName := fieldType.Name

		switch field.Kind() {
		case reflect.Struct:
			fmt.Printf("%s%s:\n", indent, fieldName)
			c.printStruct(field, indent+"  ")
		case reflect.Slice:
			fmt.Printf("%s%s: [", indent, fieldName)
			for j := 0; j < field.Len(); j++ {
				if j > 0 {
					fmt.Printf(", ")
				}
				fmt.Printf("%v", field.Index(j))
			}
			fmt.Println("]")
		default:
			fmt.Printf("%s%s: %v\n", indent, fieldName, field.Interface())
		}
	}
}
