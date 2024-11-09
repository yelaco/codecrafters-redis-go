package commands

import (
	"github.com/codecrafters-io/redis-starter-go/internal/resp"
)

func Echo(ctx CommandCtx, args []string) (resp.RespData, error) {
	return resp.RespData{
		Value: args[0],
		Type:  resp.BulkString,
	}, nil
}
