package app

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/yolocs/styles/yolodev/internal/fileutil"
	"github.com/yolocs/styles/yolodev/internal/ghostty"
	"github.com/yolocs/styles/yolodev/internal/theme"
)

func Run(ctx context.Context, stdin io.Reader, stdout, stderr io.Writer, args []string) int {
	_ = ctx
	_ = stdin
	if len(args) == 3 && args[0] == "theme" && args[1] == "validate" {
		return validateTheme(stdout, stderr, args[2])
	}
	if len(args) >= 4 && args[0] == "theme" && args[1] == "export" && args[2] == "ghostty" {
		return exportGhostty(stdout, stderr, args[3:])
	}
	printUsage(stderr)
	return 2
}

func exportGhostty(stdout, stderr io.Writer, args []string) int {
	var outputPath, themePath string
	switch {
	case len(args) == 1:
		themePath = args[0]
	case len(args) == 3 && args[0] == "--output":
		outputPath = args[1]
		themePath = args[2]
	default:
		printUsage(stderr)
		return 2
	}

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
	if err := ghostty.Render(&rendered, value); err != nil {
		fmt.Fprintf(stderr, "%s: ERROR: %v\n", themePath, err)
		return 1
	}
	if outputPath == "" {
		if _, err := io.Copy(stdout, &rendered); err != nil {
			fmt.Fprintf(stderr, "ERROR: write Ghostty theme: %v\n", err)
			return 1
		}
		return 0
	}
	if err := fileutil.WriteAtomic(outputPath, rendered.Bytes(), false); err != nil {
		fmt.Fprintf(stderr, "%s: ERROR: write Ghostty theme: %v\n", outputPath, err)
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

func printUsage(w io.Writer) {
	fmt.Fprintln(w, "usage:")
	fmt.Fprintln(w, "  yolodev theme edit [FILE]")
	fmt.Fprintln(w, "  yolodev theme validate FILE")
	fmt.Fprintln(w, "  yolodev theme export ghostty [--output PATH] FILE")
}
