package tabloid

import (
	"fmt"
	"strings"
)

func (t *Tabloid) Select(columns []Column, requestedColumnNames []string) ([]Column, error) {
	// If there are no requested columns, we return them all
	if len(requestedColumnNames) == 0 {
		return columns, nil
	}

	returnedColumns := make([]Column, 0, len(requestedColumnNames))
	for _, v := range requestedColumnNames {
		var column Column

		for _, c := range columns {
			if c.Title == v || strings.ToLower(c.Title) == v || c.ExprTitle == v {
				column = c
				break
			}
		}

		if column.ExprTitle == "" {
			return nil, fmt.Errorf("column %q does not exist in the input dataset", v)
		}

		returnedColumns = append(returnedColumns, column)
	}

	return returnedColumns, nil
}

// func (t *Tabloid) Select(columns []Column, data []map[string]interface{}, requestedColumns []string) ([]map[string]interface{}, error) {
// 	foundColumnNames := make([]string, 0, len(requestedColumns))

// 	// If there are no requested columns, we return them all
// 	for _, v := range requestedColumns {
// 		var columnExpr string

// 		for _, c := range columns {
// 			if c.Title == v || strings.ToLower(c.Title) == v || c.ExprTitle == v {
// 				columnExpr = c.ExprTitle
// 				break
// 			}
// 		}

// 		if columnExpr == "" {
// 			return nil, fmt.Errorf("column %q does not exist in the input dataset", v)
// 		}
// 	}

// 	for pos, column := range selectedColumns {
// 		for _, row := range data {
// 			value, ok := row[column.ExprTitle]
// 			if ok {
// 				selectedColumns[pos].Values = append(selectedColumns[pos].Values, value.(string))
// 			}
// 		}
// 	}

// 	t.logger.Printf("columns after select: %#v", selectedColumns)
// 	return selectedColumns, nil
// }
