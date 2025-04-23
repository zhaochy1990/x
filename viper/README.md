# Configuration Loader with Viper

Usage:

```
go get -u "github.com/zhaochy1990/x/viper"
```

```golang
import "github.com/zhaochy1990/x/viper"

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
