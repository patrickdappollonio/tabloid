package tabloid

import (
	"fmt"
	"strings"

	"github.com/antonmedv/expr"
	"github.com/antonmedv/expr/vm"
)

func (t *Tabloid) Filter(columns []Column, expression string) ([]Column, error) {
	expression = strings.TrimSpace(expression)

	if expression == "" {
		t.logger.Printf("no filter expression provided, returning all rows")
		return columns, nil
	}
	// Configure options
	options := []expr.Option{
		expr.AsBool(),
		expr.Function("ready", isready),
		expr.Function("restarted", hasrestarts),
		expr.Function("olderthan", olderThan),
		expr.Function("olderthaneq", olderThanEq),
		expr.Function("newerthan", newerThan),
		expr.Function("newerthaneq", newerThanEq),
		expr.Function("eqduration", eqduration),
	}

	// Compile the expression against an empty set of values
	program, err := expr.Compile(expression, options...)
	if err != nil {
		return nil, fmt.Errorf("unable to compile expression %q: %w", expression, err)
	}

	// Create a VM to reuse
	vm := vm.VM{}

	newColumns := make([]Column, 0, len(columns))
	for _, column := range columns {
		for pos := range column.Values {
			// Create a single row of key/value pairs
			row := make(map[string]interface{})
			for _, column := range columns {
				row[column.ExprTitle] = column.Values[pos]
			}

			// Add additional utility functions
			row = addFunctions(row)

			// Run the program against this row only
			result, err := vm.Run(program, row)
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

func addFunctions(rows map[string]interface{}) map[string]interface{} {
	if rows == nil {
		rows = make(map[string]interface{})
	}

	rows["has_restarts"] = wrapCall[bool](hasrestarts, rows["restarts"])
	rows["is_ready"] = wrapCall[bool](isready, rows["ready"])
	rows["is_running"] = wrapCall[bool](isrunning, rows["status"])
	rows["is_pending"] = wrapCall[bool](ispending, rows["status"])
	rows["is_completed"] = wrapCall[bool](iscompleted, rows["status"])
	rows["is_failed"] = wrapCall[bool](isfailed, rows["status"])
	rows["is_unknown"] = wrapCall[bool](isunknown, rows["status"])
	rows["is_succeeded"] = wrapCall[bool](issucceeded, rows["status"])
	rows["is_waiting"] = wrapCall[bool](iswaiting, rows["status"])
	rows["is_crashloopbackoff"] = wrapCall[bool](iscrashloopbackoff, rows["status"])
	rows["is_imagepullbackoff"] = wrapCall[bool](isimagepullbackoff, rows["status"])
	rows["is_errimagepull"] = wrapCall[bool](iserrimagepull, rows["status"])
	rows["is_terminated"] = wrapCall[bool](isterminated, rows["status"])

	return rows
}

// wrapCall wraps a function that returns an interface{} and an error
// it uses generics so that the return type can be changed
// and it's limited to functions with no arguments, but based on the values
// feed into the function via "args"
func wrapCall[T any](fn func(...interface{}) (interface{}, error), args interface{}) func() (T, error) {
	return func() (T, error) {
		var t T

		if args == nil {
			return t, fmt.Errorf("no column found to use as argument for the given function")
		}

		res, err := fn(args)
		if err != nil {
			return t, err
		}

		return res.(T), nil
	}
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
