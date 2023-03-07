package tabloid

import (
	"fmt"
	"strings"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/checker/decls"
	expr "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
)

func (t *Tabloid) Filter(columns []Column, expression string) ([]Column, error) {
	expression = strings.TrimSpace(expression)

	if expression == "" {
		t.logger.Printf("no filter expression provided, returning all rows")
		return columns, nil
	}

	vars := make([]*expr.Decl, 0, len(columns))
	for _, column := range columns {
		t.logger.Printf("adding column %q to CEL environment", column.ExprTitle)
		vars = append(vars, decls.NewVar(column.ExprTitle, decls.Dyn))
	}

	var opts []cel.EnvOption
	opts = append(opts, cel.ClearMacros())
	opts = append(opts, cel.Declarations(vars...))
	opts = append(opts, celFunctions...)

	env, err := cel.NewEnv(opts...)
	if err != nil {
		return nil, fmt.Errorf("unable to create common expression language environment: %w", err)
	}

	ast, issues := env.Parse(expression)
	if issues != nil && issues.Err() != nil {
		return nil, fmt.Errorf("unable to parse expression %q: %w", expression, issues.Err())
	}

	check, issues := env.Check(ast)
	if issues != nil && issues.Err() != nil {
		return nil, fmt.Errorf("unable to check expression %q: %w", expression, issues.Err())
	}

	prg, err := env.Program(check)
	if err != nil {
		return nil, fmt.Errorf("unable to create program from expression %q: %w", expression, err)
	}

	newColumns := make([]Column, 0, len(columns))
	for _, column := range columns {
		for pos := range column.Values {
			// Create a map of information to pass to the evaluation
			row := make(map[string]interface{})
			for _, column := range columns {
				row[column.ExprTitle] = column.Values[pos]
			}

			out, det, err := prg.Eval(row)
			t.logger.Printf("row: %v", row)
			t.logger.Printf("details: %v", det)
			if err != nil {
				t.logger.Printf("error type: %T", err)
				return nil, fmt.Errorf("unable to evaluate expression for row %d: %w", pos+1, err)
			}

			chosen, ok := out.Value().(bool)
			if !ok {
				return nil, fmt.Errorf("expression %q must return a boolean value, but instead it returned type \"%T\"", expression, out.Value())
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
