package unpack

import (
	"errors"
	"testing"
	"unicode"

	"github.com/stretchr/testify/require"
)

func TestUnpack(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{input: "a4bc2d5e", expected: "aaaabccddddde"},
		{input: "abccd", expected: "abccd"},
		{input: "", expected: ""},
		{input: "aaa0b", expected: "aab"},
		{input: `qwe\4\5`, expected: `qwe45`},
		{input: `qwe\45`, expected: `qwe44444`},
		{input: `qwe\\5`, expected: `qwe\\\\\`},
		{input: `qwe\\\3`, expected: `qwe\3`},
		{input: `Ä4Ñ2`, expected: `ÄÄÄÄÑÑ`},
		{input: `퉄3팿3`, expected: `퉄퉄퉄팿팿팿`},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			result, err := Unpack(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.expected, result)
		})
	}
}

func TestUnpackInvalidString(t *testing.T) {
	invalidStrings := []string{"3abc", "45", "aaa10b"}
	for _, tc := range invalidStrings {
		tc := tc
		t.Run(tc, func(t *testing.T) {
			_, err := Unpack(tc)
			require.Truef(t, errors.Is(err, ErrInvalidString), "actual error %q", err)
		})
	}
}

func TestStringOnlyDigit(t *testing.T) {
	input := "78784564"
	if result, err := Unpack(input); err == nil {
		t.Errorf("Строка содержит толькко цифры \"%s\", expect error", result)
	}
}

func TestWrongSymbolSequence(t *testing.T) {
	input := "a2b37"
	if result, err := Unpack(input); err == nil {
		t.Errorf("Неверная последовательность символов \"%s\", expect error", result)
	}
}

func TestIsSymbolDigit(t *testing.T) {
	digitsRunes := string([]rune("0123456789"))
	for _, digitRune := range digitsRunes {
		if !unicode.IsDigit(digitRune) {
			t.Errorf("Символ \"%s\" не цифра", string(digitRune))
		}
	}

	invalidDigitsRunes := string([]rune("abcXYZ!@#"))
	for _, digitRune := range invalidDigitsRunes {
		if unicode.IsDigit(digitRune) {
			t.Errorf("Символ \"%s\" цифра", string(digitRune))
		}
	}
}

func TestIsSymbolEscape(t *testing.T) {
	validEscapeRunes := string([]rune(`\`))
	for _, escapeRune := range validEscapeRunes {
		if !isSymbolEscape(escapeRune) {
			t.Errorf("Символ \"%s\" не пробел", string(escapeRune))
		}
	}

	invalidEscapeRunes := string([]rune(`|/#@!`))
	for _, escapeRune := range invalidEscapeRunes {
		if isSymbolEscape(escapeRune) {
			t.Errorf("Символ \"%s\" пробел", string(escapeRune))
		}
	}
}
