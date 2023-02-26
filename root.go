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

type settings struct {
	expr             string
	columns          []string
	debug            bool
	noTitles         bool
	titlesOnly       bool
	titlesNormalized bool
}

func rootCommand(r io.Reader) *cobra.Command {
	var opts settings

	cmd := &cobra.Command{
		Use:           "tabloid",
		Short:         helpShort,
		Long:          helpLong,
		SilenceUsage:  true,
		SilenceErrors: true,
		Version:       version,
		Example:       sliceToTabulated(examples),
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(r, os.Stdout, opts)
		},
	}

	cmd.Flags().StringVarP(&opts.expr, "expr", "e", "", "expression to filter the output")
	cmd.Flags().StringSliceVarP(&opts.columns, "column", "c", []string{}, "columns to display")
	cmd.Flags().BoolVar(&opts.debug, "debug", false, "enable debug mode")
	cmd.Flags().BoolVar(&opts.noTitles, "no-titles", false, "remove column titles from the output")
	cmd.Flags().BoolVar(&opts.titlesOnly, "titles-only", false, "only display column titles")
	cmd.Flags().BoolVar(&opts.titlesNormalized, "titles-normalized", false, "normalize column titles")

	return cmd
}

func run(r io.Reader, w io.Writer, opts settings) error {
	var b bytes.Buffer

	if _, err := io.Copy(&b, r); err != nil {
		return err
	}

	tab := tabloid.New(&b)
	tab.EnableDebug(opts.debug)

	cols, err := tab.ParseColumns()
	if err != nil {
		return err
	}

	if opts.titlesOnly {
		if opts.expr != "" {
			return fmt.Errorf("cannot use --expr with --titles-only")
		}

		if len(opts.columns) > 0 {
			return fmt.Errorf("cannot use --column with --titles-only")
		}

		for _, v := range cols {
			if opts.titlesNormalized {
				fmt.Fprintln(w, v.ExprTitle)
				continue
			}

			fmt.Fprintln(w, v.Title)
		}
		return nil
	}

	filtered, err := tab.Filter(cols, opts.expr)
	if err != nil {
		return err
	}

	output, err := tab.Select(filtered, opts.columns)
	if err != nil {
		return err
	}

	if len(output) == 0 {
		return fmt.Errorf("input had no columns to handle")
	}

	t := tabwriter.NewWriter(w, 0, 0, 3, ' ', 0)

	if !opts.noTitles {
		for _, v := range output {
			if opts.titlesNormalized {
				fmt.Fprintf(t, "%s\t", v.ExprTitle)
				continue
			}

			fmt.Fprintf(t, "%s\t", v.Title)
		}
		fmt.Fprintln(t, "")
	}

	for i := 0; i < len(output[0].Values); i++ {
		for _, v := range output {
			fmt.Fprintf(t, "%s\t", v.Values[i])
		}
		fmt.Fprintln(t, "")
	}

	if err := t.Flush(); err != nil {
		return fmt.Errorf("unable to flush table contents to screen: %w", err)
	}

	return nil
}
