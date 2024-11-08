package config

import (
	"flag"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	DEFAULT_HOST = "0.0.0.0"
	DEFAULT_PORT = 6379
)

type Config struct {
	Host string
	Port int
}

func init() {
	flag.String("host", DEFAULT_HOST, "Redis server host")
	flag.Int("port", DEFAULT_PORT, "Redis server port")
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)
}

func NewConfig() Config {
	return Config{
		Host: "0.0.0.0",
		Port: viper.GetInt("port"),
	}
}
