package count_test

import (
	"io/fs"
	"os"
	"testing"

	"github.com/namikaze-dev/clines/internal/service/count"
)

// creates a temp file and write CONTENT into it
func setupTempFile(dir, content string) (*os.File, error) {
	file, err := os.CreateTemp(dir, "*.txt")
	if err != nil {
		return nil, err
	}

	_, err = file.WriteString(content)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func TestDir(t *testing.T) {
	// should error on invalid in path
	_, err := count.Dir("invalid/dir/path", &count.Options{})
	if err == nil {
		t.Errorf("want error %v, got nil", fs.ErrNotExist)
	}

	// setup temp dir
	dir, err := os.MkdirTemp("", "")
	if err != nil {
		t.Fatal("unexpected error", err)
	}
	defer func() {
		err = os.RemoveAll(dir)
		if err != nil {
			t.Fatal("unexpected error", err)
		}
	}()

	// add two temp files to dir
	file, err := setupTempFile(dir, "line 1\nline 2\n\t\nline 3")
	if err != nil {
		t.Fatal("unexpected error", err)
	}
	defer file.Close()

	file, err = setupTempFile(dir, "\t\t\n\t\n\n")
	if err != nil {
		t.Fatal("unexpected error", err)
	}
	defer file.Close()

	got, err := count.Dir(dir, &count.Options{})
	if err != nil {
		t.Fatal("unexpected error", err)
	}

	wantCount := 3
	if wantCount != int(got.Count) {
		t.Errorf("got %v; want %v", got.Count, wantCount)
	}

	wantFiles := 2
	if wantFiles != int(got.Files) {
		t.Errorf("got %v; want %v", got.Files, wantFiles)
	}
}
