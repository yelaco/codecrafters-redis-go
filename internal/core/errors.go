package core

import "errors"

var (
	ErrInvalidRespDataValueType   = errors.New("invalid resp data value type")
	ErrInvalidCommandRespDataType = errors.New("invalid resp data type for command")
)
