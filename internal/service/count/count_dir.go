package count

import (
	"io/fs"
	"os"
	"path/filepath"
	"sync/atomic"
	"time"
)

func Dir(path string, options *Options) (*Result, error) {
	var jobs = make(chan string)
	var config = &Config{}
	defaultifyOptions(options)

	for w := 1; w <= options.Workers; w++ {
		go countWorker(&config.count, options, config, jobs)
	}

	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		options.Logger.Println(err)
		return nil, err
	}

	startTime := time.Now()
	filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			options.Logger.Println(err)
			return nil
		}

		if d.IsDir() {
			return nil
		}

		config.wg.Add(1)
		config.files++
		go func() {
			jobs <- path
		}()

		return nil
	})

	if options.Verbose {
		options.Logger.Println("walk dir complete")
	}
	config.wg.Wait()

	return &Result{Count: config.count, Files: config.files, Time: time.Since(startTime)}, nil
}

func countWorker(count *int64, options *Options, config *Config, jobs chan string) {
	for path := range jobs {
		f, err := os.Open(path)
		if err != nil {
			options.Logger.Println(err)
			continue
		}

		c := countlines(f)
		atomic.AddInt64(count, int64(c))

		if options.Verbose {
			options.Logger.Println(filepath.Base(path), c)
		}

		f.Close()
		config.wg.Done()
	}
}
