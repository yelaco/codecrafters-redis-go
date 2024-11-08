package respv2

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/codecrafters-io/redis-starter-go/internal/resp"
)

var ErrUnknownItemType = errors.New("unknown item type")

type parser struct {
	lexer   *lexer
	payload []byte
}

func NewParser(payload []byte) resp.RespParser {
	return &parser{
		lexer: lex(payload),
	}
}

func (p *parser) Parse() (resp.RespData, error) {
	var result resp.RespData
	for {
		item := p.lexer.nextItem()

		switch item.typ {
		case itemDataType:
			result.Type = resp.RespDataType(item.val[0])
		case itemData:
			switch result.Group() {
			case resp.Simple:
				result.Value = item.val
			case resp.Bulk:
				dataLength, err := strconv.Atoi(string(item.val))
				if err != nil {
					return resp.RespData{}, err
				}
				if p.lexer.nextItem().typ != itemTerminator {
					return resp.RespData{}, fmt.Errorf("%w: got %v", resp.ErrItemTerminatorExpected, item)
				}
				item = p.lexer.nextItem()
				if item.typ != itemData {
					return resp.RespData{}, fmt.Errorf("%w: got %v", resp.ErrItemDataExpected, item)
				}
				if len(item.val) != dataLength {
					return resp.RespData{}, fmt.Errorf("len(%s) != %d: %w", string(item.val), dataLength, resp.ErrArrayLengtAndDataAmountMismatch)
				}
				result.Value = item.val
			case resp.Aggregate:
				size, err := strconv.Atoi(string(item.val))
				if err != nil {
					return resp.RespData{}, err
				}
				dataArr := make([]resp.RespData, 0, size)
				for range size {
					data, err := p.Parse()
					if err != nil {
						return resp.RespData{}, err
					}
					dataArr = append(dataArr, data)
				}
				result.Value = resp.RespDataArray(dataArr)
			}
			return result, nil
		case itemTerminator:
			continue
		case itemEOF:
			return result, nil
		default:
			return resp.RespData{}, fmt.Errorf("item: %v:  %w", item, ErrUnknownItemType)
		}
	}
}
