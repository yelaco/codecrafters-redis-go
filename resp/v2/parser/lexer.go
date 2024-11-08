package parser

import (
	"bytes"
	"fmt"
	"strings"
)

type itemType int

const (
	itemError itemType = iota
	itemDataType
	itemData
	itemTerminator
	itemEOF

	terminator = "\r\n"
)

type item struct {
	typ itemType
	val []byte
}

func (i item) String() string {
	switch i.typ {
	case itemEOF:
		return "EOF"
	case itemError:
		return string(i.val)
	}
	if len(i.val) > 10 {
		return fmt.Sprintf("%.10q...", i.val)
	}
	return fmt.Sprintf("%q", i.val)
}

type stateFn func(*lexer) stateFn

// lexer holds the state of the scanner
type lexer struct {
	input []byte // the string being scanned
	start int    // start position of this item
	pos   int    // current position in the input
	state stateFn
	items chan item // channel of scanned items
}

func lex(input []byte) *lexer {
	l := &lexer{
		input: input,
		state: lexDataType,
		items: make(chan item, 2),
	}
	return l
}

// nextItem returns the next item from the input
func (l *lexer) nextItem() item {
	for {
		select {
		case item := <-l.items:
			return item
		default:
			l.state = l.state(l)
		}
	}
}

func (l *lexer) emit(t itemType) {
	l.items <- item{t, l.input[l.start:l.pos]}
	l.start = l.pos
}

func lexData(l *lexer) stateFn {
	n := len(l.input)
	for ; l.pos < n; l.next() {
		if bytes.HasPrefix(l.input[l.pos:], []byte(terminator)) {
			if l.pos >= l.start {
				l.emit(itemData)
			}
			return lexTerminator
		}
	}
	return lexEOF
}

func lexDataType(l *lexer) stateFn {
	if l.accept("*$+-") {
		l.emit(itemDataType)
		return lexData
	} else {
		return l.errorf("invalid data type: byte(%d)", l.input[l.pos])
	}
}

func lexTerminator(l *lexer) stateFn {
	l.pos += len(terminator)
	l.emit(itemTerminator)
	if l.pos >= len(l.input) {
		return lexEOF
	}
	if strings.IndexByte("*$+-", l.peek()) >= 0 {
		return lexDataType
	}
	return lexData
}

func lexEOF(l *lexer) stateFn {
	l.emit(itemEOF)
	return lexEOF
}

func (l *lexer) next() byte {
	if l.pos >= len(l.input) {
		return byte(4)
	}
	b := l.input[l.pos]
	l.pos += 1
	return b
}

func (l *lexer) ignore() {
	l.start = l.pos
}

func (l *lexer) backup() {
	l.pos -= 1
}

func (l *lexer) peek() byte {
	b := l.next()
	l.backup()
	return b
}

func (l *lexer) accept(valid string) bool {
	if strings.IndexByte(valid, l.next()) >= 0 {
		return true
	}
	l.backup()
	return false
}

func (l *lexer) acceptRun(valid string) {
	for strings.IndexByte(valid, l.next()) >= 0 {
	}
	l.backup()
}

func (l *lexer) errorf(format string, args ...interface{}) stateFn {
	l.items <- item{
		itemError,
		[]byte(fmt.Sprintf(format, args...)),
	}
	return nil
}
