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
		//если первый символ цифра, то возвращаем "" + ошибку
		if index == 0 && asciiToInteger > 0 {
			return "", ErrInvalidString
		}

		//если после цифры идет снова цифра
		prevCharAsciiToInteger, _ := strconv.Atoi(prevChar)
		if (prevCharAsciiToInteger > 0 && int(char) == 48) || (prevCharAsciiToInteger > 0 && asciiToInteger > 0) {
			return "", ErrInvalidString
		}

		if asciiToInteger > 0 { //asciiToInteger будет больше нуля только в случае, если в строке не буква.
			builder.WriteString(strings.Repeat(prevChar, asciiToInteger-1)) //Поэтому повторяем последний символ z-1 раз
			prevChar = string(char)
		}

		if asciiToInteger == 0 {
			if int(char) == 48 { //если val == 0 (0 имеет значение 48 в таблице ASCII), то удаляем последнее добавленное значение
				s := builder.String()
				result := s[:len(s)-1]
				builder.Reset()
				builder.WriteString(result)
				prevChar = string(char)
			} else { //если нет, то добавляем в билдер текущий символ
				prevChar = string(char)
				builder.WriteString(prevChar)
			}
		}
	}
	return builder.String(), nil
}
