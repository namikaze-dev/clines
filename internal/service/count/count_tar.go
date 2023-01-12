package count

import (
	"archive/tar"
	"errors"
	"io"
	"os"
	"time"
)

func Tar(path string, config *Config) (CountResult, error) {
	file, err := os.Open(path)
	if err != nil {
		return CountResult{}, err
	}

	var count int64
	trd := tar.NewReader(file)
	startTime := time.Now()
	for {
		thdr, err := trd.Next()
		if errors.Is(err, io.EOF) {
			break
		}

		if err != nil {
			config.Logger.Println(err)
			continue
		}

		if thdr.FileInfo().IsDir() {
			continue
		}

		c := countlines(trd)
		config.files++
		if config.Verbose {
			config.Logger.Println(thdr.Name, c)
		}
		count += int64(c)
	}

	return CountResult{Count: count, Files: config.files, Time: time.Since(startTime)}, nil
}
