package fileutil

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

var ErrExists = errors.New("destination exists")

func WriteAtomic(path string, data []byte, overwrite bool) (returnErr error) {
	if !overwrite {
		_, err := os.Lstat(path)
		switch {
		case err == nil:
			return fmt.Errorf("%s: %w", path, ErrExists)
		case !errors.Is(err, os.ErrNotExist):
			return fmt.Errorf("inspect destination %s: %w", path, err)
		}
	}

	directory := filepath.Dir(path)
	temporary, err := os.CreateTemp(directory, ".yolodev-*")
	if err != nil {
		return fmt.Errorf("create temporary file for %s: %w", path, err)
	}
	temporaryPath := temporary.Name()
	defer func() {
		if err := os.Remove(temporaryPath); err != nil && !errors.Is(err, os.ErrNotExist) && returnErr == nil {
			returnErr = fmt.Errorf("remove temporary file %s: %w", temporaryPath, err)
		}
	}()

	if err := temporary.Chmod(0o600); err != nil {
		temporary.Close()
		return fmt.Errorf("set temporary file permissions: %w", err)
	}
	if _, err := temporary.Write(data); err != nil {
		temporary.Close()
		return fmt.Errorf("write temporary file for %s: %w", path, err)
	}
	if err := temporary.Sync(); err != nil {
		temporary.Close()
		return fmt.Errorf("sync temporary file for %s: %w", path, err)
	}
	if err := temporary.Close(); err != nil {
		return fmt.Errorf("close temporary file for %s: %w", path, err)
	}
	if err := os.Rename(temporaryPath, path); err != nil {
		return fmt.Errorf("replace %s: %w", path, err)
	}
	return nil
}
