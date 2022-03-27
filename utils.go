package main

import (
	"bytes"
	"fmt"
)

func sliceToTabulated(slice []string) string {
	var s bytes.Buffer
	for pos, v := range examples {
		s.WriteString(fmt.Sprintf("  %s", v))

		if pos != len(examples)-1 {
			s.WriteString("\n")
		}
	}

	return s.String()
}
