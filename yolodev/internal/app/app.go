// Package app defines the yolodev command-line application.
package app

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/yolocs/styles/yolodev/internal/exporter"
	"github.com/yolocs/styles/yolodev/internal/fileutil"
	"github.com/yolocs/styles/yolodev/internal/theme"
	"github.com/yolocs/styles/yolodev/internal/tui"
)

var runEditor = tui.Run

// Run executes the yolodev command and returns a process exit code.
func Run(ctx context.Context, stdin io.Reader, stdout, stderr io.Writer, args []string) int {
	exitCode := 0
	registry := exporter.NewDefaultRegistry()
	command := newRootCommand(stdin, stdout, stderr, registry, &exitCode)
	command.SetArgs(args)

	executed, err := command.ExecuteContextC(ctx)
	if err == nil {
		return exitCode
	}
	if executed == nil {
		executed = command
	}
	fmt.Fprintf(stderr, "ERROR: %v\n\n%s", err, executed.UsageString())
	return 2
}

func newRootCommand(stdin io.Reader, stdout, stderr io.Writer, registry *exporter.Registry, exitCode *int) *cobra.Command {
	root := &cobra.Command{
		Use:           "yolodev",
		Short:         "Create and edit personal style assets",
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(_ *cobra.Command, _ []string) error {
			return errors.New("a command is required")
		},
	}
	root.SetIn(stdin)
	root.SetOut(stdout)
	root.SetErr(stderr)

	themeCommand := &cobra.Command{
		Use:   "theme",
		Short: "Work with terminal themes",
		RunE: func(_ *cobra.Command, _ []string) error {
			return errors.New("a theme command is required")
		},
	}
	themeCommand.AddCommand(
		newEditCommand(stdout, stderr, exitCode),
		newValidateCommand(stdout, stderr, exitCode),
		newExportCommand(stdout, stderr, registry, exitCode),
	)
	root.AddCommand(themeCommand)
	return root
}

func newEditCommand(stdout, stderr io.Writer, exitCode *int) *cobra.Command {
	return &cobra.Command{
		Use:   "edit [FILE]",
		Short: "Edit a canonical terminal theme",
		Args:  cobra.MaximumNArgs(1),
		Run: func(_ *cobra.Command, args []string) {
			path := "themes/yolodev/placeholder.toml"
			if len(args) == 1 {
				path = args[0]
			}
			*exitCode = editTheme(stdout, stderr, path)
		},
	}
}

func newValidateCommand(stdout, stderr io.Writer, exitCode *int) *cobra.Command {
	return &cobra.Command{
		Use:   "validate FILE",
		Short: "Validate a canonical terminal theme",
		Args:  cobra.ExactArgs(1),
		Run: func(_ *cobra.Command, args []string) {
			*exitCode = validateTheme(stdout, stderr, args[0])
		},
	}
}

func newExportCommand(stdout, stderr io.Writer, registry *exporter.Registry, exitCode *int) *cobra.Command {
	var format, outputPath string
	command := &cobra.Command{
		Use:   "export FILE",
		Short: "Export a canonical terminal theme",
		Args:  cobra.ExactArgs(1),
		PreRunE: func(_ *cobra.Command, _ []string) error {
			if format == "" {
				return errors.New("--format is required")
			}
			return nil
		},
		Run: func(_ *cobra.Command, args []string) {
			*exitCode = exportTheme(stdout, stderr, registry, format, outputPath, args[0])
		},
	}
	command.Flags().StringVar(&format, "format", "", "export format")
	command.Flags().StringVarP(&outputPath, "output", "o", "", "write output to path instead of stdout")
	return command
}

func editTheme(stdout, stderr io.Writer, path string) int {
	file, err := os.Open(path)
	if err != nil {
		fmt.Fprintf(stderr, "%s: ERROR: open terminal theme: %v\n", path, err)
		return 1
	}
	value, decodeErr := theme.Decode(file)
	closeErr := file.Close()
	if decodeErr != nil {
		fmt.Fprintf(stderr, "%s: ERROR: %v\n", path, decodeErr)
		return 1
	}
	if closeErr != nil {
		fmt.Fprintf(stderr, "%s: ERROR: close terminal theme: %v\n", path, closeErr)
		return 1
	}

	diagnostics := theme.Validate(value)
	for _, diagnostic := range diagnostics {
		output := stdout
		if diagnostic.Severity == theme.Error {
			output = stderr
		}
		fmt.Fprintf(output, "%s:%s: %s: %s\n", path, diagnostic.Field, strings.ToUpper(string(diagnostic.Severity)), diagnostic.Message)
	}
	if theme.HasErrors(diagnostics) {
		return 1
	}
	if err := runEditor(theme.Normalize(value), path); err != nil {
		fmt.Fprintf(stderr, "%s: ERROR: %v\n", path, err)
		return 1
	}
	return 0
}

func exportTheme(stdout, stderr io.Writer, registry *exporter.Registry, format, outputPath, themePath string) int {
	file, err := os.Open(themePath)
	if err != nil {
		fmt.Fprintf(stderr, "%s: ERROR: open terminal theme: %v\n", themePath, err)
		return 1
	}
	value, decodeErr := theme.Decode(file)
	closeErr := file.Close()
	if decodeErr != nil {
		fmt.Fprintf(stderr, "%s: ERROR: %v\n", themePath, decodeErr)
		return 1
	}
	if closeErr != nil {
		fmt.Fprintf(stderr, "%s: ERROR: close terminal theme: %v\n", themePath, closeErr)
		return 1
	}

	diagnostics := theme.Validate(value)
	for _, diagnostic := range diagnostics {
		fmt.Fprintf(stderr, "%s:%s: %s: %s\n", themePath, diagnostic.Field, strings.ToUpper(string(diagnostic.Severity)), diagnostic.Message)
	}
	if theme.HasErrors(diagnostics) {
		return 1
	}

	var rendered bytes.Buffer
	if err := registry.Export(format, &rendered, value); err != nil {
		fmt.Fprintf(stderr, "%s: ERROR: %v\n", themePath, err)
		return 1
	}
	if outputPath == "" {
		if _, err := io.Copy(stdout, &rendered); err != nil {
			fmt.Fprintf(stderr, "ERROR: write %s theme: %v\n", format, err)
			return 1
		}
		return 0
	}
	if err := fileutil.WriteAtomic(outputPath, rendered.Bytes(), false); err != nil {
		fmt.Fprintf(stderr, "%s: ERROR: write %s theme: %v\n", outputPath, format, err)
		return 1
	}
	return 0
}

func validateTheme(stdout, stderr io.Writer, path string) int {
	file, err := os.Open(path)
	if err != nil {
		fmt.Fprintf(stderr, "%s: ERROR: open terminal theme: %v\n", path, err)
		return 1
	}
	defer file.Close()

	value, err := theme.Decode(file)
	if err != nil {
		fmt.Fprintf(stderr, "%s: ERROR: %v\n", path, err)
		return 1
	}

	diagnostics := theme.Validate(value)
	for _, diagnostic := range diagnostics {
		output := stdout
		if diagnostic.Severity == theme.Error {
			output = stderr
		}
		fmt.Fprintf(output, "%s:%s: %s: %s\n", path, diagnostic.Field, strings.ToUpper(string(diagnostic.Severity)), diagnostic.Message)
	}
	if theme.HasErrors(diagnostics) {
		return 1
	}
	return 0
}
