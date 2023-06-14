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
		charIsNumber, _ := strconv.Atoi(string(char))
		if index == 0 && charIsNumber > 0 {
			return "", ErrInvalidString
		}

		prevCharIsNumber, _ := strconv.Atoi(prevChar)
		if (prevCharIsNumber > 0 && char == '0') || (prevCharIsNumber > 0 && charIsNumber > 0) {
			return "", ErrInvalidString
		}

		if charIsNumber > 0 {
			builder.WriteString(strings.Repeat(prevChar, charIsNumber-1))
			prevChar = string(char)
		}

		if charIsNumber == 0 {
			if char == '0' {
				s := builder.String()
				runes := []rune(s)
				result := runes[:len(runes)-1]
				builder.Reset()
				builder.WriteString(string(result))
				prevChar = string(char)
			} else {
				prevChar = string(char)
				builder.WriteString(prevChar)
			}
		}
	}
	return builder.String(), nil
}
