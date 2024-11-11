package commands

import (
	"fmt"

	"github.com/codecrafters-io/redis-starter-go/internal/replication"
	"github.com/codecrafters-io/redis-starter-go/internal/resp"
)

func Psync(ctx CommandCtx, args []string) (resp.RespData, error) {
	go func() {
		ctx.mu.Lock()
		defer ctx.mu.Unlock()
		replication.StartReplication(ctx.conn)
	}()
	return resp.RespData{
		Value: fmt.Sprintf("FULLRESYNC %v 0", ctx.serverConfig.ReplicationId),
		Type:  resp.SimpleString,
	}, nil
}
