package commands

import (
	"github.com/codecrafters-io/redis-starter-go/internal/resp"
)

func Get(ctx CommandCtx, args []string) (resp.RespData, error) {
	result := resp.RespData{
		Type: resp.BulkString,
	}
	v, ok := ctx.dict[args[0]]
	if ok {
		result.Value = v
	}
	return result, nil
}
