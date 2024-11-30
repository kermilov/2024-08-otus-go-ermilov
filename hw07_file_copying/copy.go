package main

import (
	"errors"
	"io"
	"os"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrSamePath              = errors.New("fromPath & toPath are the same")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	if isUnknownLengthFile(fromPath) {
		return ErrUnsupportedFile
	}
	if areSameFile(fromPath, toPath) {
		return ErrSamePath
	}
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
	if limit > fileFromStat.Size() {
		limit = fileFromStat.Size()
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
		if errors.Is(err, io.EOF) {
			return nil
		}
		return err
	}
	return nil
}

func isUnknownLengthFile(filename string) bool {
	fi, err := os.Stat(filename)
	if err != nil {
		return false
	}
	return fi.Size() == 0 && fi.Mode()&os.ModeDevice != 0
}

func areSameFile(path1, path2 string) bool {
	fi1, err := os.Stat(path1)
	if err != nil {
		return false
	}
	fi2, err := os.Stat(path2)
	if err != nil {
		return false
	}
	return os.SameFile(fi1, fi2)
}
