package core

import (
	"github.com/codecrafters-io/redis-starter-go/core/replication"
	"github.com/codecrafters-io/redis-starter-go/util"
)

type ServerInfo struct {
	Replication replication.InfoSection
}

func NewServerInfo(config util.Config) ServerInfo {
	return ServerInfo{
		Replication: replication.NewInfoSection(config),
	}
}
