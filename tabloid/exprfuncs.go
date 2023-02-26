package tabloid

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/Knetic/govaluate"
	str2duration "github.com/xhit/go-str2duration/v2"
)

// isready checks if a string is in the form of <current>/<total> and if the
// current value is equal to the total value, false otherwise.
func isready(args ...interface{}) (interface{}, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("isready function only accepts one argument")
	}

	str, ok := args[0].(string)
	if !ok {
		return nil, fmt.Errorf("isready function only accepts string arguments")
	}

	pieces := strings.FieldsFunc(str, func(r rune) bool {
		return r == '/'
	})

	if len(pieces) != 2 {
		return nil, fmt.Errorf("isready function only accepts string arguments in the form of <current>/<total>")
	}

	if pieces[0] != pieces[1] {
		return false, nil
	}

	return true, nil
}

var reRestart = regexp.MustCompile(`[1-9]\d*( \([^\)]+\))?`)

// hasrestarts checks if a string contains a restart count, or if it's zero.
func hasrestarts(args ...interface{}) (interface{}, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("hasrestarts function only accepts one argument")
	}

	str, ok := args[0].(string)
	if !ok {
		return nil, fmt.Errorf("hasrestarts function only accepts string arguments")
	}

	return reRestart.MatchString(str), nil
}

// parseDurations parses two string arguments into time.Duration values.
func parseDurations(args ...interface{}) (time.Duration, time.Duration, error) {
	if len(args) != 2 {
		return time.Duration(0), time.Duration(0), fmt.Errorf("olderthan function only accepts two arguments")
	}

	str, ok := args[0].(string)
	if !ok {
		return time.Duration(0), time.Duration(0), fmt.Errorf("olderthan function only accepts string arguments")
	}

	age, ok := args[1].(string)
	if !ok {
		return time.Duration(0), time.Duration(0), fmt.Errorf("olderthan function only accepts string arguments")
	}

	t1, err := str2duration.ParseDuration(str)
	if err != nil {
		return time.Duration(0), time.Duration(0), fmt.Errorf("unable to parse duration: %w", err)
	}

	t2, err := str2duration.ParseDuration(age)
	if err != nil {
		return time.Duration(0), time.Duration(0), fmt.Errorf("unable to parse duration: %w", err)
	}

	return t1, t2, nil
}

// olderthan checks if the first argument is older than the second argument,
// using Go's time.Duration parsing.
func olderThan(args ...interface{}) (interface{}, error) {
	t1, t2, err := parseDurations(args...)
	return t1 > t2, err
}

// olderthaneq checks if the first argument is older than or equal to the second
// argument, using Go's time.Duration parsing.
func olderThanEq(args ...interface{}) (interface{}, error) {
	t1, t2, err := parseDurations(args...)
	return t1 >= t2, err
}

// newerthan checks if the first argument is newer than the second argument,
// using Go's time.Duration parsing.
func newerThan(args ...interface{}) (interface{}, error) {
	t1, t2, err := parseDurations(args...)
	return t1 < t2, err
}

// newerthaneq checks if the first argument is newer than or equal to the
// second argument, using Go's time.Duration parsing.
func newerThanEq(args ...interface{}) (interface{}, error) {
	t1, t2, err := parseDurations(args...)
	return t1 <= t2, err
}

// eqduration checks if the first argument is equal to the second argument,
// using Go's time.Duration parsing.
func eqduration(args ...interface{}) (interface{}, error) {
	t1, t2, err := parseDurations(args...)
	return t1 == t2, err
}

// funcs is a map of functions that can be used in the filter expression.
var funcs = map[string]govaluate.ExpressionFunction{
	"isready": isready,
	"isnotready": func(args ...interface{}) (interface{}, error) {
		ready, err := isready(args...)
		return !ready.(bool), err
	},

	"hasrestarts": hasrestarts,
	"hasnorestarts": func(args ...interface{}) (interface{}, error) {
		restarts, err := hasrestarts(args...)
		return !restarts.(bool), err
	},

	"olderthan":   olderThan,
	"olderthaneq": olderThanEq,
	"newerthan":   newerThan,
	"newerthaneq": newerThanEq,
	"eqduration":  eqduration,
}
