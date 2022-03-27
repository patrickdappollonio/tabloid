package tabloid

import (
	"fmt"
	"strings"

	"github.com/Knetic/govaluate"
)

func (t *Tabloid) Filter(expression string) error {
	expression = strings.TrimSpace(expression)

	if expression == "" {
		t.filtered = append(t.filtered, t.contents...)
		return nil
	}

	expr, err := govaluate.NewEvaluableExpression(expression)
	if err != nil {
		return fmt.Errorf("unable to process expression %q: %w", expression, err)
	}

	for pos, row := range t.contents {
		result, err := expr.Evaluate(row)
		if err != nil {
			return fmt.Errorf("unable to parse expression for row %d: %w", pos+1, err)
		}

		chosen, ok := result.(bool)
		if !ok {
			return fmt.Errorf("%q must be a boolean expression", expression)
		}

		if chosen {
			t.filtered = append(t.filtered, row)
		}
	}

	t.logger.Printf("columns after filter: %#v", t.filtered)

	return nil
}
