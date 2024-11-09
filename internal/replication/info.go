package replication

import (
	"fmt"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/internal/config"
)

func GetInfoSection(cfg config.ServerConfig) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("role:%v\n", cfg.Role))
	sb.WriteString(fmt.Sprintf("master_replid:%v\n", cfg.ReplicationId))
	sb.WriteString(fmt.Sprintf("master_repl_offset:%v\n", config.DEFAULT_REPL_OFFSET))
	return sb.String()
}
