package gositemap

import (
	"compress/gzip"
	"log"
	"os"
)

func writeFile(dir, fullPath string, data []byte, compress bool) {
	if compress {
		fullPath = fullPath + ".gz"
	}

	fi, err := os.Stat(dir)
	if err != nil {
		_ = os.MkdirAll(dir, 0755)
	} else if !fi.IsDir() {
		log.Fatalf("[F] %s should be a directory", dir)
	}

	file, _ := os.OpenFile(fullPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	fi, err = file.Stat()
	if err != nil {
		log.Fatalf("[F] %s file not exists", fullPath)
	} else if !fi.Mode().IsRegular() {
		log.Fatalf("[F] %s should be a filename", fullPath)
	}

	if compress {
		gzipSitemap(file, data)
	} else {
		file.Write(data)
		defer file.Close()
	}
}

func gzipSitemap(file *os.File, data []byte) {
	gz := gzip.NewWriter(file)
	defer gz.Close()
	gz.Write(data)
}
