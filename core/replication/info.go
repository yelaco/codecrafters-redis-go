package replication

import (
	"fmt"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/util"
)

const (
	REPL_ID_LEN = 40
	ROLE_MASTER = "master"
	ROLE_SLAVE  = "slave"

	DEFAULT_ROLE               = ROLE_MASTER
	DEFAULT_REPLICATION_OFFSET = 0
)

type InfoSection struct {
	Role              string
	ConnectedSlaves   int
	ReplicationId     string
	ReplicationOffset int
	MasterHost        *string
	MasterPort        *string
}

func NewInfoSection(c util.Config) InfoSection {
	sec := defaultInfoSection()
	if c.MasterHost != "" && c.MasterPort != "" {
		sec.Role = ROLE_SLAVE
		sec.MasterHost = &c.MasterHost
		sec.MasterPort = &c.MasterPort
	}
	return sec
}

func defaultInfoSection() InfoSection {
	replicationId, err := util.RandomAlphanumericString(REPL_ID_LEN)
	if err != nil {
		panic(err)
	}
	return InfoSection{
		Role:              DEFAULT_ROLE,
		ConnectedSlaves:   0,
		ReplicationId:     replicationId,
		ReplicationOffset: DEFAULT_REPLICATION_OFFSET,
	}
}

func (sec InfoSection) IsReplica() bool {
	return sec.Role == ROLE_SLAVE &&
		sec.MasterHost != nil &&
		sec.MasterPort != nil
}

func (sec *InfoSection) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("role:%v\n", sec.Role))
	sb.WriteString(fmt.Sprintf("master_replid:%v\n", sec.ReplicationId))
	sb.WriteString(fmt.Sprintf("master_repl_offset:%v\n", sec.ReplicationOffset))
	return sb.String()
}
