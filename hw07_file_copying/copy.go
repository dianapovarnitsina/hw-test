package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func isSupportedFile(filename string) bool {
	supportedExtensions := []string{".txt", ".rtf", ".doc"}
	ext := filepath.Ext(filename)
	for _, supportedExt := range supportedExtensions {
		if ext == supportedExt {
			return true
		}
	}
	return false
}

func Copy(fromPath, toPath string, offset, limit int64) error {
	fileInfo, err := os.Stat(fromPath)
	if err != nil {
		return fmt.Errorf("could not get stat from file: %w", err)
	}
	fileSize := fileInfo.Size()

	if offset > fileSize {
		return ErrOffsetExceedsFileSize
	}

	if !isSupportedFile(fromPath) {
		return ErrUnsupportedFile
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
	}
	if limit+offset <= fileSize && limit > 0 {
		bufSize = limit
	}
	if limit == 0 {
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
