package main

import (
	"errors"
	"io"
	"os"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	fileFrom, err := os.OpenFile(fromPath, os.O_RDONLY, 0)
	if err != nil {
		return err
	}
	defer fileFrom.Close()
	fileFromStat, err := fileFrom.Stat()
	if err != nil {
		return err
	}
	if offset > 0 && fileFromStat.Size() < offset {
		return ErrOffsetExceedsFileSize
	}
	if limit == 0 {
		limit = fileFromStat.Size() - offset
	}
	copyContent := make([]byte, limit+1)
	n, err := fileFrom.ReadAt(copyContent, offset)
	if err != nil {
		if err == io.EOF {
			copyContent = copyContent[:n]
		} else {
			return err
		}
	}
	fileTo, err := os.Create(toPath)
	if err != nil {
		return err
	}
	defer fileTo.Close()
	_, err = fileTo.Write(copyContent)
	if err != nil {
		return err
	}
	return nil
}
