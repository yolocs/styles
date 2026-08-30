package fileutil

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestWriteAtomicCreatesDestination(t *testing.T) {
	path := filepath.Join(t.TempDir(), "theme")
	if err := WriteAtomic(path, []byte("new"), false); err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != "new" {
		t.Fatalf("destination = %q, want new", got)
	}
}

func TestWriteAtomicRefusesExistingDestination(t *testing.T) {
	path := filepath.Join(t.TempDir(), "theme")
	if err := os.WriteFile(path, []byte("original"), 0o600); err != nil {
		t.Fatal(err)
	}
	err := WriteAtomic(path, []byte("replacement"), false)
	if !errors.Is(err, ErrExists) {
		t.Fatalf("WriteAtomic() error = %v, want ErrExists", err)
	}
	got, _ := os.ReadFile(path)
	if string(got) != "original" {
		t.Fatalf("destination changed to %q", got)
	}
}

func TestWriteAtomicReplacesWhenAllowed(t *testing.T) {
	path := filepath.Join(t.TempDir(), "theme")
	if err := os.WriteFile(path, []byte("original"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := WriteAtomic(path, []byte("replacement"), true); err != nil {
		t.Fatal(err)
	}
	got, _ := os.ReadFile(path)
	if string(got) != "replacement" {
		t.Fatalf("destination = %q, want replacement", got)
	}
}
