package resp

import "errors"

var (
	ErrNilPayload                      = errors.New("nil payload")
	ErrInvalidPayloadLength            = errors.New("invalid payload length")
	ErrInvalidRespDataType             = errors.New("invalid RESP data type")
	ErrTerminatorNotFound              = errors.New("terminator not found")
	ErrInvalidBulkStringParts          = errors.New("invalid bulk string parts")
	ErrBulkStringLengthMismatch        = errors.New("bulk string length mismatch")
	ErrArrayLengtAndDataAmountMismatch = errors.New("array length and amount of data mismatch")
	ErrItemDataExpected                = errors.New("expect item data")
	ErrItemTerminatorExpected          = errors.New("expect item terminator")
)
