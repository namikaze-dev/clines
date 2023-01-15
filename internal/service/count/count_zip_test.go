package count_test

import (
	"archive/zip"
	"io/fs"
	"os"
	"testing"

	"github.com/namikaze-dev/clines/internal/service/count"
)

func TestZip(t *testing.T) {
	// should error on invalid in zip path
	_, err := count.Zip("invalid/dir/path.zip", &count.Options{})
	if err == nil {
		t.Errorf("want error %v, got nil", fs.ErrNotExist)
	}

	// setup temp zip files
	file, err := os.CreateTemp("", "*.zip")
	if err != nil {
		t.Fatal("unexpected error", err)
	}
	defer os.Remove(file.Name())
	defer file.Close()

	w := zip.NewWriter(file)

	// Add some files to the archive.
	var files = []struct {
		Name, Body string
	}{
		{"readme.txt", "This archive contains some text files."},
		{"gopher.txt", "Gopher names:\namikaze\neren\njack"},
		{"todo.txt", "\n\n\t\t\nline 1"},
	}
	for _, file := range files {
		f, err := w.Create(file.Name)
		if err != nil {
			t.Fatal("unexpected error", err)
		}
		_, err = f.Write([]byte(file.Body))
		if err != nil {
			t.Fatal("unexpected error", err)
		}
	}

	// Make sure to check the error on Close.
	err = w.Close()
	if err != nil {
		t.Fatal("unexpected error", err)
	}

	got, err := count.Zip(file.Name(), &count.Options{})
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
