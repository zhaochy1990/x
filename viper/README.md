# Configuration Loader with Viper

Usage:

```golang
// define your config here
type Config struct {
  xxxx

}

// load config

var config Config
var configPath = "config.yml"
viper.MustLoadConfig("SERVICE_X_", &configPath, &config)

fmt.Println(config)
```
