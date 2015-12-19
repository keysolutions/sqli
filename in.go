package sqli

import (
	"bytes"
	"fmt"
	"reflect"
	"unicode"
)

// In expands placeholder arguments in the supplied query when the
// associated value is a slice. To be used when query contains one
// placeholer for an IN query.
func In(query string, args ...interface{}) (string, []interface{}) {
	var buf argBuffer
	scn := argScanner{query: query}

loop:
	for {
		tok := scn.scan()
		switch tok {
		case argEOF:
			break loop
		case argQuery:
			buf.WriteString(scn.String())
		default:
			if len(args) < 1 {
				break
			}
			buf.WriteArg(scn.String(), args[0])
			args = args[1:]
		}
	}
	return buf.String(), buf.args
}

// argBuffer is a buffer which holds the written query and expanded
// arguments.
type argBuffer struct {
	bytes.Buffer
	args []interface{}
}

// WriteArg writes the placeholder and argument to the buffer, expanding
// where necessary.
func (b *argBuffer) WriteArg(placeholder string, arg interface{}) {
	if reflect.TypeOf(arg).Kind() == reflect.Slice {
		s := reflect.ValueOf(arg)
		for i := 0; i < s.Len(); i++ {
			if i > 0 {
				b.WriteString(", ")
			}
			b.writeArg(placeholder, s.Index(i).Interface())
		}
		return
	}
	b.writeArg(placeholder, arg)
}

// writeArg adds the placeholder and argument to the buffer. No expansion is
// performed here.
func (b *argBuffer) writeArg(placeholder string, arg interface{}) {
	b.args = append(b.args, arg)
	if placeholder[0] == '$' {
		b.WriteString(fmt.Sprintf("$%d", len(b.args)))
		return
	}
	b.WriteByte('?')
}

const (
	seof byte = 0

	// argScanner tokens.
	argQuery = iota
	argPlaceholder
	argEOF
)

// argScanner scans over the query to find placeholder arguments.
type argScanner struct {
	query string
	start int
	pos   int
}

// String returns the range of the string that corresponds with the
// token returned by scan(). The value of String is only valid after
// calling scan() and will change with each scan iteration.
func (s *argScanner) String() string {
	return s.query[s.start:s.pos]
}

// scan returns the next token in the scanning phase. It may be called
// until argEOF is reached.
func (s *argScanner) scan() int {
	s.start = s.pos
	switch s.next() {
	case seof:
		return argEOF
	case '?':
		return argPlaceholder
	case '$':
		return s.scanPlaceholder()
	default:
		return s.scanQuery()
	}
}

// scanPlaceholder continues looking for characters in the
// placeholder string before returning the token.
func (s *argScanner) scanPlaceholder() int {
	for unicode.IsDigit(rune(s.peek())) {
		s.next()
	}
	return argPlaceholder
}

// scanQuery continues looking for cahracters in the query string
// before returning the token.
func (s *argScanner) scanQuery() int {
	for isQuery(s.peek()) {
		s.next()
	}
	return argQuery
}

// isQuery validates that the supplied character is part of the
// query string.
func isQuery(ch byte) bool {
	return ch != '$' && ch != '?' && ch != seof
}

// peek returns the character in the query at the current pointer position.
func (s *argScanner) peek() byte {
	if s.pos >= len(s.query) {
		return seof
	}
	return s.query[s.pos]
}

// next peeks at the character in the query and moves the pointer ahead if
// not already at the end of the string.
func (s *argScanner) next() byte {
	ch := s.peek()
	if ch != seof {
		s.pos++
	}
	return ch
}
