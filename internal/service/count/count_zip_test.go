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

	zw := zip.NewWriter(file)
	// add two test files 
	f, err := zw.Create("file.txt")
	if err != nil {
		t.Fatal("unexpected error", err)
	}
	f.Write([]byte("line 1\nline 2\n\t\t\nline 3\n\t\nline 4\nline 5"))

	f, err = zw.Create("file2.txt")
	if err != nil {
		t.Fatal("unexpected error", err)
	}
	f.Write([]byte("\n\t\t\n\n\t\t\n\n\t\t\nline 1\nline 2\n\t\t\n"))
	zw.Close()

	got, err := count.Zip(file.Name(), &count.Options{})
	if err != nil {
		t.Fatal("unexpected error", err)
	}

	wantCount := 7
	if wantCount != int(got.Count) {
		t.Errorf("got %v; want %v", got.Count, wantCount)
	}

	wantFiles := 2
	if wantFiles != int(got.Files) {
		t.Errorf("got %v; want %v", got.Files, wantFiles)
	}
}
