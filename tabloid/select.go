package tabloid

import (
	"fmt"
	"strings"
)

func (t *Tabloid) Select(requestedColumns []string) ([]Column, error) {
	selected := make([]Column, 0, len(t.columns))

	if len(requestedColumns) == 0 {
		selected = append(selected, t.columns...)
	} else {
		for _, requested := range requestedColumns {
			var column Column

			for _, c := range t.columns {
				if c.Title == requested || strings.ToLower(c.Title) == requested || c.ExprTitle == requested {
					column = c
				}
			}

			if column.Title == "" {
				return nil, fmt.Errorf("column %q not available from the input", requested)
			}

			selected = append(selected, Column{
				Title:      column.Title,
				ExprTitle:  column.ExprTitle,
				StartIndex: column.StartIndex,
				EndIndex:   column.EndIndex,
				Values:     column.Values,
			})
		}
	}

	for pos, column := range selected {
		for _, row := range t.filtered {
			value, ok := row[column.ExprTitle]
			if ok {
				selected[pos].Values = append(selected[pos].Values, value.(string))
			}
		}
	}

	t.logger.Printf("columns after select: %#v", selected)
	return selected, nil
}
