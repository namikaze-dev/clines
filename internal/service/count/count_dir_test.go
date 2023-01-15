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

	// Add some files to the dir.
	var files = []struct {
		Name, Body string
	}{
		{"readme.txt", "This archive contains some text files."},
		{"gopher.txt", "Gopher names:\namikaze\neren\njack"},
		{"todo.txt", "\n\n\t\t\nline 1"},
	}
	for _, file := range files {
		f, err := setupTempFile(dir, file.Body)
		if err != nil {
			t.Fatal("unexpected error", err)
		}
		f.Close()
	}

	got, err := count.Dir(dir, &count.Options{})
	if err != nil {
		t.Fatal("unexpected error", err)
	}

	wantCount := 6
	if wantCount != int(got.Count) {
		t.Errorf("got %v; want %v", got.Count, wantCount)
	}

	wantFiles := 3
	if wantFiles != int(got.Files) {
		t.Errorf("got %v; want %v", got.Files, wantFiles)
	}
}
