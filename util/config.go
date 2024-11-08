package util

import (
	"errors"
	"flag"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var ErrInvalidMasterAddr = errors.New("invalid master host and port")

const (
	DEFAULT_HOST = "0.0.0.0"
	DEFAULT_PORT = "6379"

	ROLE_MASTER = "master"
	ROLE_SLAVE  = "slave"
)

type Config struct {
	Host       string
	Port       string
	MasterHost string
	MasterPort string
}

func init() {
	flag.String("host", DEFAULT_HOST, "Redis server host")
	flag.String("port", DEFAULT_PORT, "Redis server port")
	flag.String("replicaof", "", "Become a replica of <addr>")
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)
}

func NewConfig() (Config, error) {
	config := Config{
		Host: "0.0.0.0",
		Port: viper.GetString("port"),
	}
	masterAddr := viper.GetString("replicaof")
	if masterAddr != "" {
		m := strings.Split(masterAddr, " ")
		if len(m) != 2 {
			return Config{}, ErrInvalidMasterAddr
		}
		config.MasterHost = strings.TrimSpace(m[0])
		config.MasterPort = strings.TrimSpace(m[1])
	}

	return config, nil
}
