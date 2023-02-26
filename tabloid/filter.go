package tabloid

import (
	"fmt"
	"strings"

	"github.com/Knetic/govaluate"
)

func (t *Tabloid) Filter(columns []Column, expression string) ([]Column, error) {
	expression = strings.TrimSpace(expression)

	if expression == "" {
		t.logger.Printf("no filter expression provided, returning all rows")
		return columns, nil
	}

	expr, err := govaluate.NewEvaluableExpressionWithFunctions(expression, funcs)
	if err != nil {
		return nil, fmt.Errorf("unable to process expression %q: %w", expression, err)
	}

	newColumns := make([]Column, 0, len(columns))
	for _, column := range columns {
		for pos := range column.Values {
			row := make(map[string]interface{})
			for _, column := range columns {
				row[column.ExprTitle] = column.Values[pos]
			}

			result, err := expr.Evaluate(row)
			if err != nil {
				t.logger.Printf("error type: %T", err)
				return nil, fmt.Errorf("unable to evaluate expression for row %d: %w", pos+1, err)
			}

			chosen, ok := result.(bool)
			if !ok {
				return nil, fmt.Errorf("expression %q must return a boolean value", expression)
			}

			if chosen {
				newColumns = upsertColumn(newColumns, column, column.Values[pos])
			}
		}
	}

	return newColumns, nil
}

func upsertColumn(columns []Column, column Column, data string) []Column {
	for pos, v := range columns {
		if v.ExprTitle == column.ExprTitle {
			columns[pos].Values = append(columns[pos].Values, data)
			return columns
		}
	}

	return append(columns, Column{
		VisualPosition: column.VisualPosition,
		ExprTitle:      column.ExprTitle,
		Title:          column.Title,
		StartIndex:     column.StartIndex,
		EndIndex:       column.EndIndex,
		Values:         []string{data},
	})
}
