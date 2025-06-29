package config

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

// ConsoleLoggerConfig contains configuration for console logging
type ConsoleLoggerConfig struct {
	Enabled    bool `yaml:"enabled" json:"enabled"`
	Colors     bool `yaml:"colors" json:"colors"`
	JsonFormat bool `yaml:"json_format" json:"json_format"`
}

// FileLoggerConfig contains configuration for file logging with rotation
type FileLoggerConfig struct {
	Enabled    bool   `yaml:"enabled" json:"enabled"`
	Filename   string `yaml:"filename" json:"filename"`
	Directory  string `yaml:"directory" json:"directory"` // Directory where log files will be created
	JsonFormat bool   `yaml:"json_format" json:"json_format"`
	MaxSize    int    `yaml:"max_size" json:"max_size"`
	MaxBackups int    `yaml:"max_backups" json:"max_backups"`
	MaxAge     int    `yaml:"max_age" json:"max_age"`
	Compress   bool   `yaml:"compress" json:"compress"`
}

// LogConfig represents the complete logger configuration
// Design is extensible - new logger configs can be added here
type LogConfig struct {
	Level         string              `yaml:"level" json:"level"`
	ConsoleLogger ConsoleLoggerConfig `yaml:"console_logger" json:"console_logger"`
	FileLogger    FileLoggerConfig    `yaml:"file_logger" json:"file_logger"`
	// Future logger configs can be added here when implementations exist:
	// DatadogLogger DatadogLoggerConfig  `yaml:"datadog_logger" json:"datadog_logger"`
	// LogDNALogger  LogDNALoggerConfig   `yaml:"logdna_logger" json:"logdna_logger"`
}

func NewConfig() *viper.Viper {
	envConf := os.Getenv("APP_CONF")
	var configPath string

	if envConf == "" {
		// Support both --config and --conf flags for backwards compatibility
		flag.StringVar(&configPath, "config", "", "config path, eg: --config @/local.yml or --config configs/local.yml")
		flag.StringVar(&envConf, "conf", "", "config path (deprecated, use --config), eg: --conf configs/local.yml")

		// Add validation flag for config files
		var validateConfig bool
		flag.BoolVar(&validateConfig, "validate-config", false, "validate config file and exit")

		flag.Parse()

		// Prefer --config over --conf
		if configPath != "" {
			envConf = configPath
		}
	}

	// Handle @/configs path alias
	if strings.HasPrefix(envConf, "@/") {
		envConf = strings.Replace(envConf, "@/", "configs/", 1)
	}

	// Set default if no config specified
	if envConf == "" {
		envConf = "configs/local.yml"
	}

	conf := getConfig(envConf)
	fmt.Printf("Loaded config file: %s\n", envConf)

	// Handle validation flag
	if len(os.Args) > 1 {
		for _, arg := range os.Args[1:] {
			if arg == "--validate-config" {
				fmt.Printf("✓ Config file %s is valid\n", envConf)
				os.Exit(0)
			}
		}
	}

	return conf
}

func getConfig(path string) *viper.Viper {
	conf := viper.New()
	conf.SetConfigFile(path)
	if err := conf.ReadInConfig(); err != nil {
		panic(fmt.Errorf("failed to read config file %s: %w", path, err))
	}
	return conf
}
