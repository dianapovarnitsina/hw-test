package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
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
	env := make(Environment)
	dirEntry, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("system cannot read the directory, error: %w", err)
	}

	for _, e := range dirEntry {
		if e.IsDir() {
			continue
		}

		filePath := filepath.Join(dir, e.Name())

		buf, err := fileToByteSlice(filePath)

		if errors.Is(err, ErrNeedRemove) {
			env[e.Name()] = EnvValue{"", true}
			continue
		}
		if err != nil {
			log.Printf("cannot convert file's content to bytes, error: %v\n", err)
			continue
		}

		r := bufio.NewReader(bytes.NewReader(buf))
		envVal, err := r.ReadString('\n')

		if err != nil && !errors.Is(err, io.EOF) {
			log.Printf("cannot read string from buf, error: %v\n", err)
			continue
		}

		envVal = strings.TrimRight(envVal, "\n")
		envVal = strings.TrimRight(envVal, " ")
		for i, v := range envVal {
			if v == 0x00 {
				envVal = envVal[:i] + "\n" + envVal[i+1:]
			}
		}

		env[e.Name()] = EnvValue{envVal, false}
	}
	return env, nil
}

var ErrNeedRemove = errors.New("need to remove the env value")

func fileToByteSlice(fileName string) ([]byte, error) {
	file, err := os.OpenFile(fileName, os.O_RDONLY, 0o400)
	if err != nil {
		return nil, fmt.Errorf("cannot open file, name: %s error: %w", fileName, err)
	}
	defer file.Close()

	fInf, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("unsupported file, name: %s, error: %w", fileName, err)
	}
	if fInf.Size() == 0 {
		return nil, ErrNeedRemove
	}
	buf := make([]byte, fInf.Size())
	n, err := file.Read(buf)
	if err != nil && !errors.Is(err, io.EOF) {
		return nil, fmt.Errorf("cannot read the file, name: %s error: %w", fileName, err)
	}
	return buf[:n], nil
}
