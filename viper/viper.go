package viper

import (
	"fmt"
	"os"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
	"github.com/zhaochy1990/x/consts"
)

// MustLoadConfig load configurations from yml configuration file.
// All configurations in the yaml file can be over write by corresponding environment variable.
// It’s important to recognize that the ENV variables are case-sensitive.
// Using below environment variables to override the default behavior
func MustLoadConfig(envPrefix string, configPath string, configurations any) *viper.Viper {
	var config string
	if configPath == "" {
		if configEnv := os.Getenv(consts.ConfigPath); configEnv == "" {
			config = consts.DefaultConfigPath
			fmt.Printf("can not find configuration file path or the environment variable 'CONFIG_PATH', using default configuration file path: %s\n", config)
		} else {
			config = configEnv
			fmt.Printf("get configuration file from env:%s, value:%s\n", consts.ConfigPath, config)
		}
	} else {
		config = configPath
		fmt.Printf("you are using configuration file from: %v\n", config)
	}

	v := viper.New()
	v.SetConfigFile(config)
	v.SetConfigType("yaml")
	v.SetEnvPrefix(envPrefix)
	v.AutomaticEnv()
	// FOO-BAR => FOO_BAR
	// FOO.BAR => FOO_BAR
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_", ".", "_"))

	// Read Config from disk and env vars
	err := v.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error configuration file: %s \n", err))
	}
	fmt.Println("Config loaded successfully...")

	err = v.Unmarshal(configurations)
	if err != nil {
		panic(fmt.Errorf("Fatal unmarshal configurations: %s \n", err))
	}

	err = validator.New().Struct(configurations)
	if err != nil {
		panic(fmt.Errorf("Failed to validate configurations: %s \n", err))
	}

	return v
}
