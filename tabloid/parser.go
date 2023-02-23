package tabloid

import (
	"bufio"
	"fmt"
	"strings"
)

const endOfLine = -1

type DuplicateColumnTitleError struct {
	Title string
}

func (e *DuplicateColumnTitleError) Error() string {
	return fmt.Sprintf("duplicate column title found: %q -- unable to work with non-unique column titles", e.Title)
}

// ParseHeading parses the heading of a tabloid table and returns a list of
// columns with their respective start and end indexes. If it's the last column,
// the end index is -1. It also returns an error if there are duplicate column
// titles.
func (t *Tabloid) ParseHeading(heading string) ([]Column, error) {
	var columns []Column
	uniques := make(map[string]struct{})

	prevIndex := 0
	spaceCount := 0

	for i := 0; i < len(heading); i++ {
		if heading[i] == ' ' {
			spaceCount++
			continue
		}

		if spaceCount > 1 {
			titleTrimmed := strings.TrimSpace(heading[prevIndex:i])
			if _, ok := uniques[titleTrimmed]; ok {
				return nil, &DuplicateColumnTitleError{Title: titleTrimmed}
			}
			uniques[titleTrimmed] = struct{}{}
			columns = append(columns, Column{
				VisualPosition: len(columns) + 1,
				Title:          titleTrimmed,
				ExprTitle:      fnKey(titleTrimmed),
				StartIndex:     prevIndex,
				EndIndex:       i,
			})
			prevIndex = i
			spaceCount = 0
		}

		if len(heading)-1 == i {
			titleTrimmed := strings.TrimSpace(heading[prevIndex:])
			if _, ok := uniques[titleTrimmed]; ok {
				return nil, &DuplicateColumnTitleError{Title: titleTrimmed}
			}
			uniques[titleTrimmed] = struct{}{}
			columns = append(columns, Column{
				VisualPosition: len(columns) + 1,
				Title:          titleTrimmed,
				ExprTitle:      fnKey(titleTrimmed),
				StartIndex:     prevIndex,
				EndIndex:       endOfLine,
			})
		}
	}

	return columns, nil
}

func (t *Tabloid) ParseColumns() ([]Column, error) {
	scanner := bufio.NewScanner(t.input)

	var columns []Column

	for rowNumber := 1; scanner.Scan(); rowNumber++ {
		line := scanner.Text()

		// The first line is the header, so we use it to find the column titles
		// the assumption here is that both target apps, kubectl and docker use
		// a Go tabwriter with a padding of 3 spaces. So we use a regular expression
		// to find any empty space 3+ times and we split on them. Then, we use
		// strings.Index() to find where each column starts, we'll need it later
		// to get the column values.
		if rowNumber == 1 {
			// Find the column titles
			local, err := t.ParseHeading(line)
			if err != nil {
				return nil, err
			}

			// The first line is the header, so it doesn't need any processing
			t.logger.Printf("finished parsing columns, found: %d", len(local))
			columns = local
			continue
		}

		// Skip empty lines
		if strings.TrimSpace(line) == "" {
			t.logger.Printf("omitting empty row found in line %d", rowNumber)
			continue
		}

		// Parse each column's content
		for pos := 0; pos < len(columns); pos++ {
			// Calculate end index if it's the last column
			endIdx := columns[pos].EndIndex
			if endIdx == endOfLine {
				endIdx = len(line)
			}

			value := strings.TrimSpace(line[columns[pos].StartIndex:endIdx])

			// Store the value in the local copy of the metadata
			columns[pos].Values = append(columns[pos].Values, value)
		}
	}

	t.logger.Printf("finished parsing contents, found: %#v", columns)

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error while scanning input: %w", err)
	}

	if len(columns) == 0 {
		return nil, fmt.Errorf("no data found in input")
	}

	return columns, nil
}
