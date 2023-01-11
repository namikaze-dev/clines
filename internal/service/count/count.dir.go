package count

import (
	"io/fs"
	"os"
	"path/filepath"
	"sync/atomic"
	"time"
)

func Dir(path string, config *Config) (CountResult, error) {
	var count int64
	config.jobs = make(chan string)

	if config.Workers == 0 {
		config.Workers = 10
	}

	for w := 1; w <= config.Workers; w++ {
		go countWorker(&count, config)
	}

	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		config.Logger.Println(err)
		return CountResult{}, err
	}

	startTime := time.Now()
	filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			config.Logger.Println(err)
			return nil
		}

		if d.IsDir() {
			return nil
		}

		config.wg.Add(1)
		config.files++
		go func() {
			config.jobs <- path
		}()

		return nil
	})

	config.Logger.Println("walk dir complete")
	config.wg.Wait()

	return CountResult{Count: count, Files: config.files, Time: time.Since(startTime)}, nil
}

func countWorker(count *int64, config *Config) {
	for path := range config.jobs {
		f, err := os.Open(path)
		if err != nil {
			config.Logger.Println(err)
			continue
		}

		c := countlines(f)
		atomic.AddInt64(count, int64(c))

		if config.Verbose {
			config.Logger.Println(filepath.Base(path), c)
		}

		f.Close()
		config.wg.Done()
	}
}
