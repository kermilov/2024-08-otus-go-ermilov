package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode/utf8"
)

var ErrInvalidString = errors.New("invalid string")

var allowToSafe = map[rune]struct{}{
	'1':  {},
	'2':  {},
	'3':  {},
	'4':  {},
	'5':  {},
	'6':  {},
	'7':  {},
	'8':  {},
	'9':  {},
	'\\': {},
}

func Unpack(in string) (string, error) {
	var strBuilder strings.Builder
	var prevStr string
	var prevIsCounter bool
	var prevIsSafe bool
	var i int
	runeCount := utf8.RuneCountInString(in)
	for _, currentRune := range in {
		currentStr := string(currentRune)
		currentIsSafe := currentStr == `\` && !prevIsSafe

		repeatCount, err := strconv.Atoi(currentStr)
		currentIsCounter := err == nil && !prevIsSafe

		if i == 0 {
			if currentIsCounter {
				return "", ErrInvalidString
			}
			prevStr = currentStr
			i++
			continue
		}

		if prevIsCounter && currentIsCounter {
			return "", ErrInvalidString
		}

		if _, currentIsAllowToSafe := allowToSafe[currentRune]; prevIsSafe && !currentIsAllowToSafe {
			return "", ErrInvalidString
		}

		if currentIsCounter {
			strBuilder.WriteString(strings.Repeat(prevStr, repeatCount))
		} else if !prevIsCounter && !prevIsSafe {
			strBuilder.WriteString(prevStr)
		}
		if i == runeCount-1 {
			if currentIsSafe {
				return "", ErrInvalidString
			}
			if !currentIsCounter {
				strBuilder.WriteString(currentStr)
			}
		}
		prevStr = currentStr
		prevIsCounter = currentIsCounter
		prevIsSafe = currentIsSafe
		i++
	}
	return strBuilder.String(), nil
}
