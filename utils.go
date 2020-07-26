package gositemap

import (
	"compress/gzip"
	"fmt"
	"os"
)

func writeFile(dir, fullPath string, data []byte, compress bool) error {
	if compress {
		fullPath = fullPath + ".gz"
	}

	fi, err := os.Stat(dir)
	if err != nil {
		_ = os.MkdirAll(dir, 0750)
	} else if !fi.IsDir() {
		return fmt.Errorf("[F] %s should be a directory", dir)
	}

	file, _ := os.OpenFile(fullPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	fi, err = file.Stat()

	if err != nil {
		return fmt.Errorf("[F] %s file not exists", fullPath)
	} else if !fi.Mode().IsRegular() {
		return fmt.Errorf("[F] %s should be a filename", fullPath)
	}

	if compress {
		return gzipSitemap(file, data)
	}

	_, err = file.Write(data)
	if err != nil {
		return err
	}

	return file.Close()
}

func gzipSitemap(file *os.File, data []byte) error {
	gz := gzip.NewWriter(file)

	_, err := gz.Write(data)
	if err != nil {
		return err
	}

	return gz.Close()
}
