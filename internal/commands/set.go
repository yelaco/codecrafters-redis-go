package commands

import (
	"strconv"
	"time"

	"github.com/codecrafters-io/redis-starter-go/internal/resp"
	"github.com/codecrafters-io/redis-starter-go/internal/store"
)

type SetKeyExpireType string

const (
	EX   SetKeyExpireType = "ex"
	PX   SetKeyExpireType = "px"
	EXAT SetKeyExpireType = "exat"
	PXAT SetKeyExpireType = "pxat"
)

type SetOptions struct {
	expireType SetKeyExpireType
	expireTime int64
	setIfExist bool
	get        bool
}

func Set(ctx CommandCtx, args []string) (resp.RespData, error) {
	key := args[0]
	value := args[1]
	oldVal, ok := store.Get(key)
	store.Set(key, value)
	if len(args) >= 4 {
		expireTime, err := strconv.Atoi(args[3])
		if err != nil {
			return resp.RespData{}, err
		}
		switch SetKeyExpireType(args[2]) {
		case EX:
			time.AfterFunc(time.Duration(expireTime)*time.Second, func() {
				curVal, _ := store.Get(key)
				if curVal == value {
					store.Delete(key)
				}
			})
		case PX:
			time.AfterFunc(time.Duration(expireTime)*time.Millisecond, func() {
				curVal, _ := store.Get(key)
				if curVal == value {
					store.Delete(key)
				}
			})
		case EXAT:
		case PXAT:
		}
	}

	if false {
		result := resp.RespData{
			Type: resp.BulkString,
		}
		if ok {
			result.Value = oldVal
		}
		return result, nil
	}

	return resp.RespData{Value: "OK", Type: resp.SimpleString}, nil
}
