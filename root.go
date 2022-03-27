package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"text/tabwriter"

	"github.com/patrickdappollonio/tabloid/tabloid"
	"github.com/spf13/cobra"
)

var version = "development"

const (
	helpShort = "tabloid is a simple command line tool to parse and filter column-based CLI outputs from commands like kubectl or docker"
	helpLong  = `tabloid is a simple command line tool to parse and filter column-based CLI outputs from commands like kubectl or docker.
For documentation, see https://github.com/patrickdappollonio/tabloid`
)

var examples = []string{
	`kubectl api-resources | tabloid --expr 'kind == "Namespace"'`,
	`kubectl api-resources | tabloid --expr 'apiversion =~ "networking"'`,
	`kubectl api-resources | tabloid --expr 'shortnames == "sa"' --column name,shortnames`,
	`kubectl get pods --all-namespaces | tabloid --expr 'name =~ "^frontend" || name =~ "redis$"'`,
}

func rootCommand(r io.Reader) *cobra.Command {
	var (
		expr     string
		columns  []string
		debug    bool
		noTitles bool
	)

	cmd := &cobra.Command{
		Use:           "tabloid",
		Short:         helpShort,
		Long:          helpLong,
		SilenceUsage:  true,
		SilenceErrors: true,
		Version:       version,
		Example:       sliceToTabulated(examples),
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(r, os.Stdout, expr, columns, debug, noTitles)
		},
	}

	cmd.Flags().StringVarP(&expr, "expr", "e", "", "expression to filter the output")
	cmd.Flags().StringSliceVarP(&columns, "column", "c", []string{}, "columns to display")
	cmd.Flags().BoolVar(&debug, "debug", false, "enable debug mode")
	cmd.Flags().BoolVar(&noTitles, "no-titles", false, "remove column titles from the output")

	return cmd
}

func run(r io.Reader, w io.Writer, expr string, columns []string, debug bool, noTitles bool) error {
	var b bytes.Buffer

	if _, err := io.Copy(&b, r); err != nil {
		return err
	}

	tab := tabloid.New(&b)
	tab.EnableDebug(debug)

	if err := tab.ParseColumns(); err != nil {
		return err
	}

	if err := tab.Filter(expr); err != nil {
		return err
	}

	data, err := tab.Select(columns)
	if err != nil {
		return err
	}

	t := tabwriter.NewWriter(w, 0, 0, 3, ' ', 0)

	if !noTitles {
		for _, v := range data {
			fmt.Fprintf(t, "%s\t", v.Title)
		}
		fmt.Fprintln(t, "")
	}

	for i := 0; i < len(data[0].Values); i++ {
		for _, v := range data {
			fmt.Fprintf(t, "%s\t", v.Values[i])
		}
		fmt.Fprintln(t, "")
	}

	if err := t.Flush(); err != nil {
		return fmt.Errorf("unable to flush table contents to screen: %w", err)
	}

	return nil
}
