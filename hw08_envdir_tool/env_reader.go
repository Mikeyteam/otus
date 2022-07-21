package main

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"log"
	"os"
	"path"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	_, err := haveFileinDir(dir)
	if err != nil {
		return nil, err
	}

	environmentMap := make(map[string]EnvValue)
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if !validateFileName(file.Name()) {
			continue
		}
		if !file.IsDir() {
			pathToFile := path.Join(dir, file.Name())
			fileInfo, err := file.Info()
			if err != nil {
				return nil, err
			}
			if fileInfo.Size() == 0 {
				environmentMap[file.Name()] = EnvValue{Value: "", NeedRemove: true}
				continue
			}
			fileContent, err := readingFile(pathToFile)
			if err != nil {
				return nil, err
			}
			environmentMap[file.Name()] = EnvValue{Value: fileContent, NeedRemove: false}
		}
	}

	return environmentMap, nil
}

// haveFileinDir returns whether the given file or directory exists.
func haveFileinDir(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}

	return false, err
}

// validateFileName read fileName and validate by symbol.
func validateFileName(fileName string) bool {
	for _, symbol := range fileName {
		if symbol == '=' {
			return false
		}
	}

	return true
}

// validateFileName read fileName and validate by symbol.
func readingFile(pathToFile string) (string, error) {
	file, err := os.Open(pathToFile)
	if err != nil {
		log.Println("Open file ", err)
		return "", err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	fileContent, err := reader.ReadBytes(byte('\n'))

	if err != nil && !errors.Is(err, io.EOF) {
		log.Println("BufferReader ", err)
		return "", err
	}

	fileContent = bytes.ReplaceAll(fileContent, []byte{0x00}, []byte{'\n'})

	return strings.TrimRight(string(fileContent), " \t\n"), nil
}
