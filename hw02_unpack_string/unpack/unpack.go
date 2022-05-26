package unpack

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("неверный формат строки")

func Unpack(str string) (string, error) {
	var stringBuffer string
	var flag bool
	var resultString strings.Builder

	for _, symbol := range str {
		if unicode.IsDigit(symbol) {
			if stringBuffer == "" {
				return "", ErrInvalidString
			}

			if stringBuffer != "" && !flag {
				count, _ := strconv.Atoi(string(symbol))
				// Дублируем символ количеством - count и записываем в буфер
				resultString.WriteString(strings.Repeat(stringBuffer, count))
				stringBuffer = ""
				continue
			}
		}

		if stringBuffer != "" {
			if isSymbolEscape(symbol) && !flag {
				flag = true
				continue
			}

			resultString.WriteString(stringBuffer)
			stringBuffer = string(symbol)
			flag = false
			continue
		}
		stringBuffer = string(symbol)
	}
	resultString.WriteString(stringBuffer)

	// Методом String() - окончательно собираем строку из буфера
	return resultString.String(), nil
}

// Проверка. Символ это пробел.
func isSymbolEscape(r rune) bool {
	return r == 92
}
