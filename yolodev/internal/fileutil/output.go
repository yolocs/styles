package fileutil

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ResolveOutputPath returns the file path used for an export destination.
func ResolveOutputPath(path, stem, extension string) (string, error) {
	if path == "~" || strings.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("resolve home directory: %w", err)
		}
		if path == "~" {
			path = home
		} else {
			path = filepath.Join(home, filepath.FromSlash(strings.TrimPrefix(path, "~/")))
		}
	}

	info, err := os.Stat(path)
	switch {
	case err == nil && info.IsDir():
		return outputPathInDirectory(path, stem, extension)
	case err == nil:
		return path, nil
	case errors.Is(err, os.ErrNotExist) && filepath.Ext(filepath.Base(path)) == "":
		return outputPathInDirectory(path, stem, extension)
	case errors.Is(err, os.ErrNotExist):
		return path, nil
	default:
		return "", fmt.Errorf("inspect output path %s: %w", path, err)
	}
}

func outputPathInDirectory(directory, stem, extension string) (string, error) {
	if strings.ContainsAny(stem, `/\`) {
		return "", fmt.Errorf("derive output filename: theme variant %q contains a path separator", stem)
	}
	return filepath.Join(directory, stem+extension), nil
}
