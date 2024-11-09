package config

import (
	"flag"
	"strings"
	"sync"

	"github.com/codecrafters-io/redis-starter-go/pkg/util"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	REPL_ID_LEN = 40
)

type ServerConfig struct {
	Host          string
	Port          string
	MasterHost    string
	MasterPort    string
	Role          string
	ReplicationId string
}

var (
	srvCfg *ServerConfig
	mu     sync.RWMutex
)

func init() {
	flag.String("host", DEFAULT_SERVER_HOST, "Redis server host")
	flag.String("port", DEFAULT_SERVER_PORT, "Redis server port")
	flag.String("replicaof", "", "Become a replica of <addr>")
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)
}

func newServerConfig() *ServerConfig {
	replId, err := util.RandomAlphanumericString(REPL_ID_LEN)
	if err != nil {
		panic(err)
	}
	cfg := &ServerConfig{
		Host:          "0.0.0.0",
		Port:          viper.GetString("port"),
		Role:          DEFAULT_REPL_ROLE,
		ReplicationId: replId,
	}
	masterAddr := viper.GetString("replicaof")
	if masterAddr != "" {
		m := strings.Split(masterAddr, " ")
		if len(m) != 2 {
			panic(ErrInvalidMasterAddr)
		}
		cfg.MasterHost = m[0]
		cfg.MasterPort = m[1]
		cfg.Role = ROLE_SLAVE
	}

	return cfg
}

func GetServerConfig() ServerConfig {
	mu.RLock()
	defer mu.RUnlock()
	if srvCfg == nil {
		srvCfg = newServerConfig()
	}

	return *srvCfg
}
