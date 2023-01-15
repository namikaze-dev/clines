package count_test

import (
	"archive/tar"
	"io/fs"
	"os"
	"testing"

	"github.com/namikaze-dev/clines/internal/service/count"
)

func TestTar(t *testing.T) {
	// should error on invalid in zip path
	_, err := count.Tar("invalid/dir/path.tar", &count.Options{})
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

	tw := tar.NewWriter(file)
	// add two test files
	var files = []struct {
		Name, Body string
	}{
		{"readme.txt", "This archive contains some text files."},
		{"gopher.txt", "Gopher names:\namikaze\neren\njack"},
		{"todo.txt", "\n\n\t\t\nline 1"},
	}
	for _, file := range files {
		hdr := &tar.Header{
			Name: file.Name,
			Mode: 0600,
			Size: int64(len(file.Body)),
		}
		if err := tw.WriteHeader(hdr); err != nil {
			t.Fatal("unexpected error", err)
		}
		if _, err := tw.Write([]byte(file.Body)); err != nil {
			t.Fatal("unexpected error", err)
		}
	}
	if err := tw.Close(); err != nil {
		t.Fatal("unexpected error", err)
	}

	got, err := count.Tar(file.Name(), &count.Options{Verbose: true})
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
