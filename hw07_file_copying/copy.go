package main

import (
	"errors"
	"fmt"
	"github.com/cheggaaa/pb/v3"
	"io"
	"os"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {

	fileInfo, err := os.Stat(fromPath)
	if err != nil {
		return fmt.Errorf("could not get stat from file: %w", err)
	}
	fileSize := fileInfo.Size()

	if offset > fileSize {
		return ErrOffsetExceedsFileSize
	}

	fileFrom, err := os.Open(fromPath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("error file does not exist: %w", err)
		}
		return fmt.Errorf("error open from file: %w", err)
	}
	defer fileFrom.Close()

	_, err = fileFrom.Seek(offset, io.SeekStart)
	if err != nil {
		return fmt.Errorf("error seeking to offset: %w", err)
	}

	var bufSize int64
	if limit+offset > fileSize {
		bufSize = fileSize - offset
	} else if limit+offset <= fileSize && limit > 0 {
		bufSize = limit
	} else if limit == 0 {
		bufSize = fileSize - offset
	}

	fileTo, err := os.Create(toPath)
	if err != nil {
		return fmt.Errorf("error creating destination file: %w", err)
	}
	defer fileTo.Close()
	bar := pb.Full.Start64(fileSize)

	_, err = io.Copy(fileTo, bar.NewProxyReader(io.LimitReader(fileFrom, bufSize)))
	if err != nil {
		return fmt.Errorf("error copying data: %w", err)
	}
	bar.Finish()
	return nil
}
