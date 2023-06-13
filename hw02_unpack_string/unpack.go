package main

import (
	"errors"
	"strconv"
	"strings"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(str string) (string, error) {
	var builder strings.Builder
	var prevChar string

	for index, char := range str {
		asciiToInteger, _ := strconv.Atoi(string(char))
		if index == 0 && asciiToInteger > 0 {
			return "", ErrInvalidString
		}

		prevCharASCIIToInteger, _ := strconv.Atoi(prevChar)
		if (prevCharASCIIToInteger > 0 && int(char) == 48) || (prevCharASCIIToInteger > 0 && asciiToInteger > 0) {
			return "", ErrInvalidString
		}

		if asciiToInteger > 0 {
			builder.WriteString(strings.Repeat(prevChar, asciiToInteger-1))
			prevChar = string(char)
		}

		if asciiToInteger == 0 {
			if int(char) == 48 {
				s := builder.String()
				result := s[:len(s)-1]
				builder.Reset()
				builder.WriteString(result)
				prevChar = string(char)
			} else {
				prevChar = string(char)
				builder.WriteString(prevChar)
			}
		}
	}
	return builder.String(), nil
}
