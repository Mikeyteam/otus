package main

import (
	"errors"
	"io"
	"os"

	"github.com/cheggaaa/pb"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

// Copy content file to new file.
func Copy(fromPath, toPath string, offset, limit int64) error {
	file, err := os.Open(fromPath)
	if err != nil {
		return err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}

	lenghtFile := fileInfo.Size()
	if offset > 0 {
		if offset > lenghtFile {
			return ErrOffsetExceedsFileSize
		}
		file.Seek(offset, io.SeekStart)
		lenghtFile -= offset
	}
	if lenghtFile == 0 {
		return ErrUnsupportedFile
	}
	if limit == 0 || limit > lenghtFile {
		limit = lenghtFile
	}
	newFile, err := os.Create(toPath)
	if err != nil {
		return err
	}
	defer newFile.Close()

	progressBar := pb.New(int(limit)).SetUnits(pb.U_BYTES)
	progressBar.Start()
	_, err = io.CopyN(newFile, progressBar.NewProxyReader(file), limit)
	progressBar.Finish()

	return err
}
