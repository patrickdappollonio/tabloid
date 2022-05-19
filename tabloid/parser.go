package tabloid

import (
	"bufio"
	"fmt"
	"regexp"
	"strings"
)

var rePadding = regexp.MustCompile(` {3,}\b`)

const endOfLine = -1

func (t *Tabloid) ParseColumns() error {
	scanner := bufio.NewScanner(t.input)

	for rowNumber := 1; scanner.Scan(); rowNumber++ {
		line := scanner.Text()

		// The first line is the header, so we use it to find the column titles
		// the assumption here is that both target apps, kubectl and docker use
		// a Go tabwriter with a padding of 3 spaces. So we use a regular expression
		// to find any empty space 3+ times and we split on them. Then, we use
		// strings.Index() to find where each column starts, we'll need it later
		// to get the column values.
		if rowNumber == 1 {
			// Holder used to find if column titles do repeat
			titles := make(map[string]string)

			// Parse all the columns and find their start index
			cols := rePadding.Split(line, -1)
			t.logger.Printf("found columns: %v", cols)

			for _, columnTitle := range cols {
				data := Column{
					Title:     columnTitle,
					ExprTitle: fnKey(columnTitle),
				}

				// Generate a regexp to find location based on the column title
				expr := `(^|\b)` + regexp.QuoteMeta(columnTitle) + `(\b| {3,}|$)`

				t.logger.Printf("compiled regexp expression: %s", expr)

				re, err := regexp.Compile(expr)
				if err != nil {
					return fmt.Errorf("unable to find starting index of column %q, unable to compile regexp: %q: %w", columnTitle, expr, err)
				}

				// Find the start based off the last used position
				pos := re.FindStringIndex(line)
				if len(pos) != 2 {
					return fmt.Errorf("unable to find starting index of column %q, unable to find column title boundaries", columnTitle)
				}
				data.StartIndex = pos[0]

				// Add the column metadata to the list of columns
				t.columns = append(t.columns, data)

				// Find if the title already exists (check if dupe)
				for _, v := range titles {
					if v == columnTitle {
						return fmt.Errorf("duplicated column title %q -- unable to work with non-unique column titles", columnTitle)
					}
				}

				titles[fnKey(columnTitle)] = columnTitle
			}

			// For each column, calculate the end index based on the previous
			// column start index
			for pos := 0; pos < len(t.columns); pos++ {
				if pos == 0 {
					t.logger.Printf("found column named %q, located between %d and %d", t.columns[pos].Title, t.columns[pos].StartIndex, t.columns[pos].EndIndex)
					continue
				}

				if pos == len(t.columns)-1 {
					t.columns[pos].EndIndex = endOfLine
				}

				t.columns[pos-1].EndIndex = t.columns[pos].StartIndex - 1
				t.logger.Printf("found column named %q, located between %d and %d", t.columns[pos].Title, t.columns[pos].StartIndex, t.columns[pos].EndIndex)
			}

			// The first column does not need any more processing
			t.logger.Printf("finished parsing columns, found: %d", len(t.columns))
			continue
		}

		if strings.TrimSpace(line) == "" {
			t.logger.Printf("omitting empty row found in line %d", rowNumber)
			continue
		}

		// Create a local copy of stored metadata
		local := make(map[string]interface{})

		// Parse each column's content
		for pos := 0; pos < len(t.columns); pos++ {
			var value string

			// check if start and end indexes can be used in the string
			if t.columns[pos].StartIndex > len(line)-1 || (t.columns[pos].EndIndex != endOfLine && t.columns[pos].EndIndex > len(line)-1) {
				start, end := fmt.Sprintf("%d", t.columns[pos].StartIndex+1), fmt.Sprintf("%d", t.columns[pos].EndIndex+1)

				if t.columns[pos].EndIndex == endOfLine {
					end = "the end of the line"
				}

				return fmt.Errorf(
					"input line %d does not contain enough data to fill column %q: the line is %d characters long, column expects data between characters %s and %s",
					pos+1, t.columns[pos].Title, len(line), start, end,
				)
			} else {
				if t.columns[pos].EndIndex == endOfLine {
					value = line[t.columns[pos].StartIndex:]
				} else {
					value = line[t.columns[pos].StartIndex:t.columns[pos].EndIndex]
				}
			}

			local[fnKey(t.columns[pos].Title)] = strings.TrimSpace(value)
		}

		t.contents = append(t.contents, local)
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error while scanning input: %w", err)
	}

	if len(t.contents) == 0 {
		return fmt.Errorf("no data found in input")
	}

	t.logger.Printf("contents parsed: %#v", t.contents)

	return nil
}
