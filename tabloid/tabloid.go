package tabloid

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
)

type Tabloid struct {
	input    *bytes.Buffer
	logger   *log.Logger
	columns  []Column
	contents []map[string]interface{}
	filtered []map[string]interface{}
}

type Column struct {
	Title      string
	ExprTitle  string
	StartIndex int
	EndIndex   int
	Values     []string
}

func New(input *bytes.Buffer) *Tabloid {
	return &Tabloid{
		input:    input,
		logger:   log.New(ioutil.Discard, "[debug] ", log.LstdFlags),
		columns:  make([]Column, 0),
		contents: make([]map[string]interface{}, 0),
		filtered: make([]map[string]interface{}, 0),
	}
}

func (t *Tabloid) EnableDebug(debug bool) {
	if debug {
		t.logger.SetOutput(os.Stderr)
	} else {
		t.logger.SetOutput(ioutil.Discard)
	}
}

func (t *Tabloid) Columns() []Column {
	return t.columns
}

func (t *Tabloid) Meta() []map[string]interface{} {
	return t.contents
}
