package count

import (
	"bufio"
	"io"
	"strings"
)

func countlines(r io.Reader) int {
	var count int
	bs := bufio.NewScanner(r)
	for bs.Scan() {
		line := bs.Text()
		if strings.TrimSpace(line) != "" {
			count++
		}
	}
	return count
}
