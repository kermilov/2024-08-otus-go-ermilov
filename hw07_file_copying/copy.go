package main

import (
	"errors"
	"io"
	"os"

	"github.com/cheggaaa/pb/v3"
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
	fileTo, err := os.Create(toPath)
	if err != nil {
		return err
	}
	defer fileTo.Close()
	bar := pb.Full.Start64(limit)
	defer bar.Finish()
	fileFromSectionReader := io.NewSectionReader(fileFrom, offset, limit)
	barReader := bar.NewProxyReader(fileFromSectionReader)
	_, err = io.CopyN(fileTo, barReader, limit)
	if err != nil {
		return err
	}
	return nil
}
