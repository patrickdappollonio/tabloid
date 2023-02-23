package tabloid

import (
	"fmt"
	"strings"
	"unicode"
)

type DuplicateColumnTitleError struct {
	Title string
}

func (e *DuplicateColumnTitleError) Error() string {
	return fmt.Sprintf("duplicate column title found: %q -- unable to work with non-unique column titles", e.Title)
}

func fnKey(s string) string {
	s = strings.ToLower(s)

	out := make([]rune, 0, len(s))

	for _, v := range s {
		if unicode.IsLetter(v) || unicode.IsDigit(v) || v == ' ' || v == '-' {
			switch v {
			case ' ', '-':
				out = append(out, '_')
			default:
				out = append(out, v)
			}
		}
	}

	return string(out)
}
