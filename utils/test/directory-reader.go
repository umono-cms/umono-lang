package test

import (
	"os"
	"strings"
)

type DirectoryReader struct {
	path string
}

func NewDirectoryReader(path string) *DirectoryReader {
	return &DirectoryReader{
		path: path,
	}
}

func (dr *DirectoryReader) ReadWithoutExt() ([]string, error) {
	entries, err := os.ReadDir(dr.path)
	if err != nil {
		return []string{}, err
	}

	files := []string{}
	for _, entry := range entries {
		if !entry.IsDir() {
			pieces := strings.Split(entry.Name(), ".")
			files = append(files, pieces[0])
		}
	}

	return files, nil
}
