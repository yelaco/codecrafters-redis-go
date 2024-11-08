package core

import (
	"fmt"
	"time"

	"github.com/codecrafters-io/redis-starter-go/resp"
)

const (
	PING = "PING"
	ECHO = "ECHO"
	SET  = "SET"
	GET  = "SET"
	INFO = "INFO"
)

type SetKeyExpireType string

const (
	EX   SetKeyExpireType = "EX"
	PX   SetKeyExpireType = "PX"
	EXAT SetKeyExpireType = "EXAT"
	PXAT SetKeyExpireType = "PXAT"
)

func (c *Core) Ping() (resp.RespData, error) {
	return resp.RespData{
		Value: "PONG",
		Type:  resp.SimpleString,
	}, nil
}

func (c *Core) Echo(s string) (resp.RespData, error) {
	return resp.RespData{
		Value: s,
		Type:  resp.BulkString,
	}, nil
}

type SetOptions struct {
	expireType SetKeyExpireType
	expireTime int64
	setIfExist bool
	get        bool
}

func (c *Core) Set(key string, value string, options SetOptions) (resp.RespData, error) {
	oldVal, ok := c.db[key]
	c.db[key] = value
	if options.setIfExist == ok {
		switch options.expireType {
		case EX:
			time.AfterFunc(time.Duration(options.expireTime)*time.Second, func() {
				if c.db[key] == value {
					delete(c.db, key)
				}
			})
		case PX:
			time.AfterFunc(time.Duration(options.expireTime)*time.Millisecond, func() {
				if c.db[key] == value {
					delete(c.db, key)
				}
			})
		case EXAT:
		case PXAT:
		}
	}

	if options.get {
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

func (c *Core) Get(key string) (resp.RespData, error) {
	result := resp.RespData{
		Type: resp.BulkString,
	}
	v, ok := c.db[key]
	if ok {
		result.Value = v
	}
	return result, nil
}

func (c *Core) Info(section string) (resp.RespData, error) {
	switch section {
	case "replication":
		return resp.RespData{
			Value: fmt.Sprintf("# Replication\n%v", c.serverInfo.Replication.String()),
			Type:  resp.BulkString,
		}, nil
	default:
		return resp.RespData{
			Value: "ERROR: unknown section",
			Type:  resp.SimpleError,
		}, nil
	}
}

type ReplConfigOptions struct {
	listeningPort string
	capabilities  []string
}

func (c *Core) ReplConfig(options ReplConfigOptions) (resp.RespData, error) {
	return resp.RespData{Value: "OK", Type: resp.SimpleString}, nil
}
