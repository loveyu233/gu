package tools

import (
	"io"
	"net/http"
	"os"
)

func GetFileContentType(file *os.File) (string, error) {
	buffer := make([]byte, 512)

	_, err := io.ReadFull(file, buffer)
	if err != nil {
		return "", err
	}

	if _, err = file.Seek(0, io.SeekStart); err != nil {
		return "", err
	}

	contentType := http.DetectContentType(buffer)
	return contentType, nil
}
