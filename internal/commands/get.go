package commands

import (
	"github.com/codecrafters-io/redis-starter-go/internal/resp"
	"github.com/codecrafters-io/redis-starter-go/internal/store"
)

func Get(ctx CommandCtx, args []string) (resp.RespData, error) {
	result := resp.RespData{
		Type: resp.BulkString,
	}
	v, ok := store.Get(args[0])
	if ok {
		result.Value = v
	}
	return result, nil
}
