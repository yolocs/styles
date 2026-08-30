package app

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/yolocs/styles/yolodev/internal/theme"
)

func Run(ctx context.Context, stdin io.Reader, stdout, stderr io.Writer, args []string) int {
	_ = ctx
	_ = stdin
	if len(args) == 3 && args[0] == "theme" && args[1] == "validate" {
		return validateTheme(stdout, stderr, args[2])
	}
	printUsage(stderr)
	return 2
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
