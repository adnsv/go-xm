package xm

import (
	"strings"
)

const AttrQuotationMark = '\''

// ScrambleFunc is a generic string scrambler that replaces
// codeunits matched by f with xml character references.
func ScrambleFunc(s string, f func(byte) bool) []byte {
	c, i := find_byte_func(s, f)
	if i < 0 {
		return []byte(s)
	}
	b := strings.Builder{}
	b.Grow(len(s))
	for {
		b.WriteString(s[:i])
		switch c {
		case '&':
			b.WriteString("&amp;")
		case '<':
			b.WriteString("&lt;")
		case '>':
			b.WriteString("&gt;")
		case '\'':
			b.WriteString("&apos;")
		case '"':
			b.WriteString("&quot;")
		default:
			b.WriteByte('&')
			b.WriteByte('#')
			b.WriteByte(hex_chars[c>>4])
			b.WriteByte(hex_chars[c&0b1111])
			b.WriteByte(';')
		}
		s = s[i+1:]
		c, i = find_byte_func(s, f)
		if i < 0 {
			break
		}
	}
	b.WriteString(s)
	return []byte(b.String())
}

// ScrambleAttr is a scrambler for attribute values.
func ScrambleAttr(s string) RawAttr {
	return RawAttr(ScrambleFunc(s, attr_scramble))
}

// ScrambleCont is a scrambler for content.
func ScrambleCont(s string) RawCont {
	return RawCont(ScrambleFunc(s, content_scramble))
}

const hex_chars = "0123456789abcdef"

func find_byte_func(s string, f func(b byte) bool) (byte, int) {
	i, n := 0, len(s)
	for i < n {
		c := s[i]
		if f(c) {
			return c, i
		} else {
			i++
		}
	}
	return 0, -1
}

func attr_scramble(b byte) bool {
	return b < 0x20 || b == '<' || b == '&' || b == '>' || b == AttrQuotationMark
}

func content_scramble(b byte) bool {
	return b == '<' || b == '&' || b == '>' || b == 0
}
