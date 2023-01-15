package count

import (
	"archive/tar"
	"errors"
	"io"
	"os"
	"time"
)

func Tar(path string, options *Options) (*Result, error) {
	var config = &Config{}
	defaultifyOptions(options)

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	trd := tar.NewReader(file)
	startTime := time.Now()
	for {
		thdr, err := trd.Next()
		if errors.Is(err, io.EOF) {
			break
		}

		if err != nil {
			options.Logger.Println(err)
			continue
		}

		if thdr.FileInfo().IsDir() {
			continue
		}

		c := countlines(trd)
		config.files++
		if options.Verbose {
			options.Logger.Println(thdr.Name, c)
		}
		config.count += int64(c)
	}

	return &Result{Count: config.count, Files: config.files, Time: time.Since(startTime)}, nil
}
