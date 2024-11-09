package commands

import (
	"fmt"

	"github.com/codecrafters-io/redis-starter-go/internal/resp"
)

func Psync(ctx CommandCtx, args []string) (resp.RespData, error) {
	return resp.RespData{
		Value: fmt.Sprintf("+FULLRESYNC %v 0\r\n", ctx.serverConfig.ReplicationId),
		Type:  resp.SimpleString,
	}, nil
}
