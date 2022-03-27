package main

import (
	"fmt"
	"os"
)

func main() {
	if err := rootCommand(os.Stdin).Execute(); err != nil {
		errfn("Error: %s", err)
		os.Exit(1)
	}
}

func errfn(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
}
