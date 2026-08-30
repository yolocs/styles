package fileutil

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestWriteAtomic(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		initial   string
		data      string
		overwrite bool
		wantErr   error
		want      string
	}{
		{name: "creates destination", data: "new", want: "new"},
		{name: "refuses existing destination", initial: "original", data: "replacement", wantErr: ErrExists, want: "original"},
		{name: "replaces when allowed", initial: "original", data: "replacement", overwrite: true, want: "replacement"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			path := filepath.Join(t.TempDir(), "theme")
			if test.initial != "" {
				if err := os.WriteFile(path, []byte(test.initial), 0o600); err != nil {
					t.Fatal(err)
				}
			}
			err := WriteAtomic(path, []byte(test.data), test.overwrite)
			if !errors.Is(err, test.wantErr) {
				t.Fatalf("WriteAtomic() error = %v, want %v", err, test.wantErr)
			}
			got, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			if string(got) != test.want {
				t.Fatalf("destination = %q, want %q", got, test.want)
			}
		})
	}
}
