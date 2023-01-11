package count

import (
	"bytes"
	"testing"
)

func TestCountLines(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		fileContent := "hello world\nhello Go\nbye world\nbye Go"
		r := bytes.NewReader([]byte(fileContent))

		want := 4
		got := countlines(r)
		
		if got != want {
			t.Errorf("got: %v; want: %v", got, want)
		}
	})

	t.Run("empty file", func(t *testing.T) {
		fileContent := ""
		r := bytes.NewReader([]byte(fileContent))

		want := 0
		got := countlines(r)
		
		if got != want {
			t.Errorf("got: %v; want: %v", got, want)
		}

		fileContent = "\n\n\n\t\t\t\n"
		r = bytes.NewReader([]byte(fileContent))

		want = 0
		got = countlines(r)
		
		if got != want {
			t.Errorf("got: %v; want: %v", got, want)
		}
	})

	t.Run("mixed file", func(t *testing.T) {
		fileContent := "hello world\n\n\nhello Go\n\t\nbye world\nbye Go\n\n\n"
		r := bytes.NewReader([]byte(fileContent))

		want := 4
		got := countlines(r)
		
		if got != want {
			t.Errorf("got: %v; want: %v", got, want)
		}
	})
}
