package util

import (
	"github.com/spf13/viper"
)

// Config store all configuration of the application
// the values read by viper from file or enviroment variables
type Config struct {
	Enviroment string `mapstructure:"ENVIROMENT"`
	MongoURL   string `mapstructure:"MONGO_URL"`
	WebPort    string `mapstructure:"WEB_PORT"`
	RpcPort    string `mapstructure:"RPC_PORT"`
	GRpcPort   string `mapstructure:"GRPC_PORT"`
}

// LoadConfig read configuration from file conf.env or enviroment variables
func LoadConfig(configPath string) (config Config, err error) {
	v := viper.New()
	v.SetConfigName("conf")
	v.SetConfigType("env")
	v.AddConfigPath(configPath)
	err = v.ReadInConfig()
	if err != nil {
		return
	}
	v.AutomaticEnv()
	err = v.Unmarshal(&config)
	return
}
