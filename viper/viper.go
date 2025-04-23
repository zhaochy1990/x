package viper

import (
	"fmt"
	"os"
	"strings"

	"github.com/qinrundev/x/consts"
	"github.com/spf13/viper"
	"github.com/go-playground/validator/v10"
)

// MustLoadConfig load configurations from yml configuration file. 
// All configurations in the yaml file can be over write by corresponding environment variable.
// It’s important to recognize that the ENV variables are case-sensitive.
func MustLoadConfig(envPrefix string, configPath *string, configurations any) *viper.Viper {
	var config string
	if configPath == nil {
		if configEnv := os.Getenv(consts.ConfigPath); configEnv == "" {
			config = consts.DefaultConfigPath
			fmt.Printf("can not find viper path or viper path environment variable, using default viper path: %s\n", config)
		} else {
			config = configEnv
			fmt.Printf("get viper path from env:%s, value:%s\n", consts.ConfigPath, config)
		}
	} else {
		config = *configPath
		fmt.Printf("you are using viper file from: %v\n", config)
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
		panic(fmt.Errorf("Fatal error viper file: %s \n", err))
	}
	fmt.Println("Config loaded successfully...")

	err = v.Unmarshal(&configurations)
	if err != nil {
		panic(fmt.Errorf("Fatal unmarshal configurations: %s \n", err))
	}

	validate := validator.New()
	err = validate.Struct(configurations)
	if err != nil {
		panic(fmt.Errorf("Failed to validate configurations: %s \n", err))
	}

	return v
}
