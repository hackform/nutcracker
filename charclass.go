package nutcracker

import (
	"strings"
)

const (
	spaceCharSet = " \t\r\n"
)

func isSpace(c byte) bool {
	switch c {
	case ' ', '\t', '\r', '\n':
		return true
	default:
		return false
	}
}

func trimLSpace(s string) string {
	return strings.TrimLeft(s, spaceCharSet)
}

func isSpecialStrI(c byte) bool {
	switch c {
	case '$', '`', '"', '\\', '\n':
		return true
	default:
		return false
	}
}

func isNewline(c byte) bool {
	return c == '\n'
}

func unquoteArg(text string) (string, error) {
	s := strings.Builder{}
	for len(text) > 0 {
		k := strings.Index(text, "\\")
		if k < 0 {
			k = len(text)
			s.WriteString(text[0:k])
			text = text[k:]
			break
		}
		s.WriteString(text[0:k])
		text = text[k+1:]
		if len(text) < 1 {
			return "", errInvalidEscape
		}
		ch := text[0]
		if !isNewline(ch) {
			s.WriteByte(ch)
		}
		text = text[1:]
	}
	return s.String(), nil
}

func unquoteStrI(text string) (string, error) {
	s := strings.Builder{}
	for len(text) > 0 {
		k := strings.Index(text, "\\")
		if k < 0 {
			k = len(text)
			s.WriteString(text[0:k])
			text = text[k:]
			break
		}
		s.WriteString(text[0:k])
		text = text[k+1:]
		if len(text) < 1 {
			return "", errInvalidEscape
		}
		ch := text[0]
		if isNewline(ch) {
		} else if isSpecialStrI(ch) {
			s.WriteByte(ch)
		} else {
			s.Write([]byte{'\\', ch})
		}
		text = text[1:]
	}
	return s.String(), nil
}
