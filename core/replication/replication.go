package replication

import (
	"fmt"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/util"
)

const (
	DEFAULT_ROLE               = "master"
	DEFAULT_REPLICATION_OFFSET = 0

	REPL_ID_LEN = 40
)

type InfoSection struct {
	Role              string
	ConnectedSlaves   int
	ReplicationId     string
	ReplicationOffset int
}

func NewInfoSection() InfoSection {
	return defaultInfoSection()
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

func (sec *InfoSection) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("role:%v\n", sec.Role))
	sb.WriteString(fmt.Sprintf("master_replid:%v\n", sec.ReplicationId))
	sb.WriteString(fmt.Sprintf("master_repl_offset:%v\n", sec.ReplicationOffset))
	return sb.String()
}
