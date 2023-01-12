package count

import (
	"archive/zip"
	"io"
	"sync/atomic"
	"time"
)

type zipJob struct {
	rd io.ReadCloser
	name string
}

func Zip(path string, config *Config) (CountResult, error) {
	var count int64
	jobs := make(chan zipJob)

	if config.Workers == 0 {
		config.Workers = 10
	}

	for w := 1; w <= config.Workers; w++ {
		go zipWorker(&count, config, jobs)
	}

	zr, err := zip.OpenReader(path)
	if err != nil {
		return CountResult{}, err
	}
	defer zr.Close()

	startTime := time.Now()
	for _, f := range zr.File {
		if f.FileInfo().IsDir() {
			continue
		}

		file, err := f.Open()
		if err != nil {
			config.Logger.Panicln(err)
			continue
		}

		config.wg.Add(1)
		go func(name string) {
			jobs <- zipJob{rd: file, name: name}
		}(f.Name)
	}

	config.Logger.Println("walk zip complete")
	config.wg.Wait()

	return CountResult{Count: count, Files: config.files, Time: time.Since(startTime)}, nil
}

func zipWorker(count *int64, config *Config, jobs <-chan zipJob) {
	for job := range jobs {
		c := countlines(job.rd)
		atomic.AddInt64(count, int64(c))
		config.files++

		if config.Verbose {
			config.Logger.Println(job.name, c)
		}

		job.rd.Close()
		config.wg.Done()
	}
}
