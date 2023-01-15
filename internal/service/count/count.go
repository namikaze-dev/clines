package count

import (
	"bufio"
	"io"
	"log"
	"strings"
	"sync"
	"time"
)

type Options struct {
	Verbose bool
	Workers int
	Logger *log.Logger
}

type Config struct {
	wg   sync.WaitGroup
	files int
	count int64
}

type Result struct {
	Count int64
	Time time.Duration
	Files int
}

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

func defaultifyOptions(options *Options) {
	if options.Workers == 0 {
		options.Workers = 10
	}
	
	if options.Logger == nil {
		options.Logger = log.Default()
	}
}