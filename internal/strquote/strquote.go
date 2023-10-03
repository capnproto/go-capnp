// Package strquote provides a function for formatting a string as a
// Cap'n Proto string literal.
package strquote

// Append appends a Cap'n Proto string literal of s to buf.
func Append(buf []byte, s []byte) []byte {
	buf = append(buf, '"')
	for _, b := range s {
		switch b {
		case '\a':
			buf = append(buf, '\\', 'a')
		case '\b':
			buf = append(buf, '\\', 'b')
		case '\f':
			buf = append(buf, '\\', 'f')
		case '\n':
			buf = append(buf, '\\', 'n')
		case '\r':
			buf = append(buf, '\\', 'r')
		case '\t':
			buf = append(buf, '\\', 't')
		case '\v':
			buf = append(buf, '\\', 'v')
		case '\'':
			buf = append(buf, '\\', '\'')
		case '"':
			buf = append(buf, '\\', '"')
		case '\\':
			buf = append(buf, '\\', '\\')
		default:
			if needsEscape(b) {
				buf = append(buf, '\\', 'x', hexDigit(b/16), hexDigit(b%16))
			} else {
				buf = append(buf, b)
			}
		}
	}
	buf = append(buf, '"')
	return buf
}

func needsEscape(b byte) bool {
	return b < 0x20 || b >= 0x7f
}

func hexDigit(b byte) byte {
	const digits = "0123456789abcdef"
	return digits[b]
}
