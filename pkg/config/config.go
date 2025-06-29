package config

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

// NewConfig creates a new Viper config instance.
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
