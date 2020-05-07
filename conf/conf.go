package conf

import (
	"fmt"
	"path/filepath"

	"github.com/jacobstr/confer"
)

// TODO:
// - allow command line arguments to be passed
// - eg. --APP_PORT=1234, should be bubbled to the right place

var configuration *confer.Config

func init() {
	LoadDefaultConfig()
}

func LoadDefaultConfig() {
	configuration = confer.NewConfig()
	seek := []string{
		"production.yaml",
		"dev.yaml",
		"config.yaml",
		"docker_config.yaml",
		"./config/production.yaml",
		"./config/dev.yaml",
		"./config/config.yaml",
	}
	var err error
	for _, file := range seek {
		err := configuration.ReadPaths(file)
		if err == nil {
			abs, _ := filepath.Abs(file)
			fmt.Println("Configuration loaded:", abs)
			break
		}
	}
	if err != nil {
		fmt.Println("No configuration file found")
	}
}

func Get(key string, defValue interface{}) interface{} {
	if configuration.IsSet(key) {
		return configuration.Get(key)
	} else {
		return defValue
	}
}

func Int(key string, defValue int) int {
	if Exists(key) {
		return configuration.GetInt(key)
	} else {
		return defValue
	}
}

func Float64(key string, defValue float64) float64 {
	if Exists(key) {
		return configuration.GetFloat64(key)
	} else {
		return defValue
	}
}

func String(key string, defValue string) string {
	if Exists(key) {
		return configuration.GetString(key)
	} else {
		return defValue
	}
}

func StringSlice(key string, defValue []string) []string {
	if Exists(key) {
		return configuration.GetStringSlice(key)
	} else {
		return defValue
	}
}

func Bool(key string, defValue bool) bool {
	if Exists(key) {
		return configuration.GetBool(key)
	} else {
		return defValue
	}
}

func Exists(key string) bool {
	return configuration.IsSet(key)
}
