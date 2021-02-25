package config

import (
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const envPrefix string = "LOXPROM"

// ReadConfigErr is returned if something goes wrong while reading the config
// We use this error to return a meaningful string
// instaead of a viper error object
type ReadConfigErr struct {
	err string
}

func (e *ReadConfigErr) Error() string {
	return e.err
}

// Config holds our config values
type Config struct {
	Host     string
	Port     int
	User     string
	Password string
}

// NewConfig reads the config into a new Config object
func NewConfig() (*Config, error) {

	cfg := new(Config)

	// Flags
	pflag.String("configFile", "", "Path and name of Config")
	pflag.String("host", "", "URL of the Miniserver")
	pflag.Int("port", 8080, "Port of the Miniserver")
	pflag.String("user", "", "Username for Miniserver")
	pflag.String("password", "", "Password for Miniserver")
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)

	// Config file
	if file := viper.GetString("configFile"); file != "" {
		viper.SetConfigFile(file)
	} else {
		viper.SetConfigName("loxone-prometheus-exporter")
		viper.SetConfigType("yml")
		viper.AddConfigPath(defaultConfigPath)
	}
	err := viper.ReadInConfig()
	if err != nil {
		switch err.(type) {
		case viper.ConfigFileNotFoundError:
			// It's not mandatory to have a config file
		default:
			return nil, &ReadConfigErr{fmt.Sprintf("Unable to read config file %s: %v", viper.GetViper().ConfigFileUsed(), err)}
		}
	}

	// Environment Variables
	viper.SetEnvPrefix(envPrefix)
	viper.BindEnv("Host")
	viper.BindEnv("User")
	viper.BindEnv("Password")

	err = viper.Unmarshal(cfg)
	if err != nil {
		return nil, &ReadConfigErr{fmt.Sprintf("Unable to marshal config: %v", err)}
	}
	return cfg, nil
}
