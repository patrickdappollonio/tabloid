package tabloid

import (
	"bytes"
	"io"
	"log"
	"os"
)

type Logger interface {
	Println(v ...interface{})
	Printf(format string, v ...interface{})
	SetOutput(w io.Writer)
}

type Tabloid struct {
	input  *bytes.Buffer
	logger Logger
}

type Column struct {
	VisualPosition int
	Title          string
	ExprTitle      string
	StartIndex     int
	EndIndex       int
	Values         []string
}

func New(input *bytes.Buffer) *Tabloid {
	return &Tabloid{
		input:  input,
		logger: log.New(io.Discard, "ℹ️ --> ", log.Lshortfile),
	}
}

func (t *Tabloid) EnableDebug(debug bool) {
	if debug {
		t.logger.SetOutput(os.Stderr)
	} else {
		t.logger.SetOutput(io.Discard)
	}
}
