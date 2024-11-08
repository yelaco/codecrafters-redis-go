package resp

import (
	"fmt"
)

type (
	RespDataType  byte
	RespDataGroup int
)

const (
	SimpleString RespDataType = '+'
	SimpleError  RespDataType = '-'
	Integer      RespDataType = ':'
	BulkString   RespDataType = '$'
	Array        RespDataType = '*'

	Simple RespDataGroup = iota
	Bulk
	Aggregate
)

type RespParser interface {
	Parse() (RespData, error)
}

type RespData struct {
	Value interface{}
	Type  RespDataType
}

type RespDataArray []RespData

func (d RespData) String() string {
	switch d.Type {
	case SimpleString:
		return fmt.Sprintf("+%v\r\n", d.Value)
	case SimpleError:
		return fmt.Sprintf("-%v\r\n", d.Value)
	case BulkString:
		if d.Value == nil {
			return "$-1\r\n"
		}
		s, _ := d.Value.(string)
		return fmt.Sprintf("$%d\r\n%s\r\n", len(s), s)
	case Array:
		return ""
	default:
		return ""
	}
}

func (d RespData) Group() RespDataGroup {
	switch d.Type {
	case SimpleString, SimpleError, Integer:
		return Simple
	case BulkString:
		return Bulk
	case Array:
		return Aggregate
	default:
		return -1
	}
}
