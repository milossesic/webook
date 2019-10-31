package config

import (
	"log"
	"os"
	"strings"
	"time"

	"git.corp.adobe.com/dc/notifications_load_test/util"
	"github.com/imdario/mergo"
	"github.com/kelseyhightower/envconfig"
	"github.com/spf13/viper"
)

var (
	configInstance *AppConfig
	once           util.Once
)

// AppConfig is application config
type AppConfig struct {
	AppPort        string `split_words:"true"`
	LogLevel       string `split_words:"true"`
	RequestTimeout time.Duration
}

// GetConfig method returns the application config
func GetConfig() *AppConfig {
	once.Do(func() {
		configInstance = generateConfig()
	})
	return configInstance
}

// UnsetConfig method sets value of configInstance to nil
func UnsetConfig() {
	configInstance = nil
	once.Reset()
}

// generateConfig method generates the application config
// using combination of environment variable and config files
// note that environment variable takes precedence over config files
func generateConfig() *AppConfig {
	output := ReadFromEnv()
	config := ReadFromConfig(strings.ToLower(os.Getenv("ENVIRONMENT_NAME")))
	mergo.Merge(output, config)
	return output
}

// ReadFromEnv method reads the environment variables
func ReadFromEnv() *AppConfig {
	var dest AppConfig
	err := envconfig.Process("", &dest)
	if err != nil {
		log.Fatal(err.Error())
		return nil
	}
	return &dest
}

// ReadFromConfig method reads the values from configuration files
// for the given environment
func ReadFromConfig(env string) *AppConfig {
	var dest AppConfig
	// reading configuration for the given environment
	// looking for the env yaml file (like dev.yaml, stage.yaml and prod.yaml)
	if env != "default" && env != "" {
		viper.SetConfigName(env)
		// looks for the yaml file in directory
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")
		// alternate path, looks for the yaml file in config directory
		viper.AddConfigPath("./config")
		// reads the config values
		if err := viper.ReadInConfig(); err != nil {
			log.Fatalf("Error reading config file, %s", err)
		}
		// maps the read values to AppConfig struct
		err := viper.Unmarshal(&dest)
		if err != nil {
			log.Fatalf("unable to decode into struct, %v", err)
		}
	}
	defaultconfiguration := readDefaultConfig()
	// merges the source (default configuration) and
	// environment specific configuration (like dev,stage,prod etc..)
	mergo.Merge(&dest, defaultconfiguration)
	return &dest
}

// readDefaultConfig method reads the default configuration from default.yaml
// located either in root folder or config folder
func readDefaultConfig() *AppConfig {
	var defaultconfiguration AppConfig
	viper.SetConfigName("default")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("../config")
	viper.AddConfigPath("../../config")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}
	// maps the read values to AppConfig struct
	err := viper.Unmarshal(&defaultconfiguration)
	if err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}
	return &defaultconfiguration
}
