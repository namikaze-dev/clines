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

func Zip(path string, options *Options) (*Result, error) {
	var jobs = make(chan zipJob)
	var config = &Config{}
	defaultifyOptions(options)

	for w := 1; w <= options.Workers; w++ {
		go zipWorker(&config.count, options, config, jobs)
	}

	zr, err := zip.OpenReader(path)
	if err != nil {
		return nil, err
	}
	defer zr.Close()

	startTime := time.Now()
	for _, f := range zr.File {
		if f.FileInfo().IsDir() {
			continue
		}

		file, err := f.Open()
		if err != nil {
			options.Logger.Println(err)
			continue
		}

		config.wg.Add(1)
		go func(name string) {
			jobs <- zipJob{rd: file, name: name}
		}(f.Name)
	}

	if options.Verbose {
		options.Logger.Println("walk dir complete")
	}
	config.wg.Wait()

	return &Result{Count: config.count, Files: config.files, Time: time.Since(startTime)}, nil
}

func zipWorker(count *int64, options *Options, config *Config, jobs <-chan zipJob) {
	for job := range jobs {
		c := countlines(job.rd)
		atomic.AddInt64(count, int64(c))
		config.files++

		if options.Verbose {
			options.Logger.Println(job.name, c)
		}

		job.rd.Close()
		config.wg.Done()
	}
}
