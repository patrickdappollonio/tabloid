package tabloid

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	str2duration "github.com/xhit/go-str2duration/v2"
)

// isready checks if a string is in the form of <current>/<total> and if the
// current value is equal to the total value, false otherwise.
func isready(args ...interface{}) (interface{}, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("isready function requires one argument")
	}

	str, ok := args[0].(string)
	if !ok {
		return nil, fmt.Errorf("isready function requires string arguments")
	}

	pieces := strings.FieldsFunc(str, func(r rune) bool {
		return r == '/'
	})

	if len(pieces) != 2 {
		return nil, fmt.Errorf("isready function requires string arguments in the form of <current>/<total>")
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
		return nil, fmt.Errorf("hasrestarts function requires one argument")
	}

	str, ok := args[0].(string)
	if !ok {
		return nil, fmt.Errorf("hasrestarts function requires string arguments")
	}

	return reRestart.MatchString(str), nil
}

// parseDurations parses two string arguments into time.Duration values.
func parseDurations(args ...interface{}) (time.Duration, time.Duration, error) {
	if len(args) != 2 {
		return time.Duration(0), time.Duration(0), fmt.Errorf("olderthan function requires two arguments")
	}

	str, ok := args[0].(string)
	if !ok {
		return time.Duration(0), time.Duration(0), fmt.Errorf("olderthan function requires string arguments")
	}

	age, ok := args[1].(string)
	if !ok {
		return time.Duration(0), time.Duration(0), fmt.Errorf("olderthan function requires string arguments")
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

// compareFirstArgumentToString compares the first argument to the given string.
func compareFirstArgumentToString(fnName string, str string, args ...interface{}) (bool, error) {
	if len(args) != 1 {
		return false, fmt.Errorf("%s function requires one argument", fnName)
	}

	if args[0] == nil {
		return false, fmt.Errorf("empty argument for %s function", fnName)
	}

	arg, ok := args[0].(string)
	if !ok {
		return false, fmt.Errorf("%s function requires string arguments", fnName)
	}

	return arg == str, nil
}

// isrunning checks if the first argument is equal to the string "Running".
func isrunning(args ...interface{}) (interface{}, error) {
	return compareFirstArgumentToString("isrunning", "Running", args...)
}

func ispending(args ...interface{}) (interface{}, error) {
	return compareFirstArgumentToString("ispending", "Pending", args...)
}

func iscompleted(args ...interface{}) (interface{}, error) {
	return compareFirstArgumentToString("iscompleted", "Completed", args...)
}

func isfailed(args ...interface{}) (interface{}, error) {
	return compareFirstArgumentToString("isfailed", "Failed", args...)
}

func isunknown(args ...interface{}) (interface{}, error) {
	return compareFirstArgumentToString("isunknown", "Unknown", args...)
}

func issucceeded(args ...interface{}) (interface{}, error) {
	return compareFirstArgumentToString("issucceeded", "Succeeded", args...)
}

func iswaiting(args ...interface{}) (interface{}, error) {
	return compareFirstArgumentToString("iswaiting", "Waiting", args...)
}

func iscrashloopbackoff(args ...interface{}) (interface{}, error) {
	return compareFirstArgumentToString("iscrashloopbackoff", "CrashLoopBackOff", args...)
}

func isimagepullbackoff(args ...interface{}) (interface{}, error) {
	return compareFirstArgumentToString("isimagepullbackoff", "ImagePullBackOff", args...)
}

func iserrimagepull(args ...interface{}) (interface{}, error) {
	return compareFirstArgumentToString("iserrimagepull", "ErrImagePull", args...)
}

func isterminated(args ...interface{}) (interface{}, error) {
	return compareFirstArgumentToString("isterminated", "Terminated", args...)
}
