package viper

import (
	"fmt"
	"os"
	"strings"

	"github.com/qinrundev/x/consts"
	"github.com/spf13/viper"
)

// MustLoadConfig load configurations from yml viper file
// All services should provide a configuration file along with its docker images.
// And all configurations inside the viper file can be over write by corresponding environment variable
// It’s important to recognize that Viper treats ENV variables as case-sensitive.
func MustLoadConfig(envPrefix string, configPath *string, configurations interface{}) *viper.Viper {
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

	return v
}
