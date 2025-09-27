package pack

import (
	"errors"
	"strconv"
	"unicode"
)

var (
	ErrInvalidString = errors.New("некорректная строка")
)

const reverseSolidus = '\\'

// %s	the uninterpreted bytes of the string or slice
// %q	a double-quoted string safely escaped with Go syntax
// %x	base 16, lower-case, two characters per byte
// %X	base 16, upper-case, two characters per byte

func UnpackString(input string, escapeEnabled bool) (string, error) {
	if len(input) == 0 {
		return "", nil
	}

	runes := []rune(input)
	var result []rune

	for i := 0; i < len(runes); i++ {
		v := runes[i]

		if escapeEnabled && v == reverseSolidus {
			if i+1 >= len(runes) {
				return "", ErrInvalidString
			}

			nextChar := runes[i+1]
			if !unicode.IsDigit(nextChar) && nextChar != reverseSolidus {
				return "", ErrInvalidString
			}

			result = append(result, nextChar)
			i++
			continue
		}

		switch {
		case unicode.IsDigit(v) && len(result) == 0:
			return "", ErrInvalidString
		case unicode.IsDigit(v):
			if i+1 < len(runes) && unicode.IsDigit(runes[i+1]) {
				return "", ErrInvalidString
			}

			if v == '0' {
				if len(result) > 0 {
					result = result[:len(result)-1]
				}
			} else {
				number, err := strconv.Atoi(string(v))
				if err != nil {
					return "", ErrInvalidString
				}

				if len(result) > 0 {
					prevChar := result[len(result)-1]
					for j := 1; j < number; j++ {
						result = append(result, prevChar)
					}
				}
			}
		default:
			result = append(result, v)
		}
	}

	if escapeEnabled && len(runes) > 0 && runes[len(runes)-1] == reverseSolidus {
		return "", ErrInvalidString
	}

	return string(result), nil
}

func PackString(input string) (string, error) {
	if len(input) == 0 {
		return "", nil
	}

	runes := []rune(input)
	var result []rune

	for i := 0; i < len(runes); {
		currentChar := runes[i]
		count := 1

		for i+count < len(runes) && runes[i+count] == currentChar {
			count++
		}

		result = append(result, currentChar)

		if count > 1 {
			result = append(result, rune('0'+count))
		}

		i += count
	}

	return string(result), nil
}
