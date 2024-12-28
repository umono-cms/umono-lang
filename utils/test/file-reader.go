package test

import (
	"io"
	"os"
	"strings"
)

type FileReader struct {
	path      string
	extension string
}

func NewFileReader(path, extension string) *FileReader {
	return &FileReader{
		path:      path,
		extension: extension,
	}
}

func (fr *FileReader) Read(fileName string, withoutNewLine bool) (string, error) {

	file, err := os.Open(fr.path + "/" + fileName + "." + fr.extension)
	if err != nil {
		return "", err
	}
	defer file.Close()

	cont, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}

	if withoutNewLine {
		return strings.TrimRight(string(cont), "\n"), nil
	}

	return string(cont), nil
}
