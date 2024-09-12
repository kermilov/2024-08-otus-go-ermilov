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
	var prevRune rune
	var prevIsCounter bool
	var prevIsSafe bool
	runeCount := utf8.RuneCountInString(in)
	for i, currentRune := range in {
		currentStr := string(currentRune)
		prevStr := string(prevRune)
		currentIsSafe := currentStr == `\` && (prevStr != `\` || !prevIsSafe)

		repeatCount, err := strconv.Atoi(currentStr)
		currentIsCounter := err == nil && !prevIsSafe

		if i == 0 {
			if currentIsCounter {
				return "", ErrInvalidString
			}
			prevRune = currentRune
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
			strBuilder.WriteRune(prevRune)
		}
		if i == runeCount-1 && !currentIsCounter {
			strBuilder.WriteString(currentStr)
		}
		prevRune = currentRune
		prevIsCounter = currentIsCounter
		prevIsSafe = currentIsSafe
	}
	return strBuilder.String(), nil
}
