package tabloid

import (
	"strings"
	"unicode"
)

func fnKey(s string) string {
	s = strings.ToLower(s)

	out := make([]rune, 0, len(s))

	for _, v := range s {
		if unicode.IsLetter(v) || unicode.IsDigit(v) || v == ' ' {
			if v == ' ' {
				out = append(out, '_')
			} else {
				out = append(out, v)
			}
		}
	}

	return string(out)
}
