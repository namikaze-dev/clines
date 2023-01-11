package count

import (
	"bytes"
	"testing"

	"github.com/namikaze-dev/clines/internal/assert"
)

func TestCountLines(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		fileContent := "hello world\nhello Go\nbye world\nbye Go"
		r := bytes.NewReader([]byte(fileContent))

		want, got := 4, countlines(r)		
		assert.Equal(t, got, want)
	})

	t.Run("empty file", func(t *testing.T) {
		fileContent := ""
		r := bytes.NewReader([]byte(fileContent))

		want, got := 0, countlines(r)
		assert.Equal(t, got, want)

		fileContent = "\n\n\n\t\t\t\n"
		r = bytes.NewReader([]byte(fileContent))

		want, got = 0, countlines(r)
		assert.Equal(t, got, want)
	})

	t.Run("mixed file", func(t *testing.T) {
		fileContent := "hello world\n\n\nhello Go\n\t\nbye world\nbye Go\n\n\n"
		r := bytes.NewReader([]byte(fileContent))

		want, got := 4, countlines(r)
		assert.Equal(t, got, want)
	})
}
