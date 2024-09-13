package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode/utf8"
)

var ErrInvalidString = errors.New("invalid string")

var allowToSafe = map[string]struct{}{
	"1":  {},
	"2":  {},
	"3":  {},
	"4":  {},
	"5":  {},
	"6":  {},
	"7":  {},
	"8":  {},
	"9":  {},
	"\\": {},
}

func Unpack(in string) (string, error) {
	runeCount := utf8.RuneCountInString(in)
	if runeCount == 0 {
		return "", nil
	}
	var prevIsCounter bool
	var prevIsSafe bool
	inRunes := []rune(in)
	prevStr := string(inRunes[0])
	if _, err := strconv.Atoi(prevStr); err == nil {
		return "", ErrInvalidString
	}
	var strBuilder strings.Builder
	for i := 1; i < runeCount; i++ {
		currentStr := string(inRunes[i])
		currentIsSafe := currentStr == `\` && !prevIsSafe

		repeatCount, err := strconv.Atoi(currentStr)
		currentIsCounter := err == nil && !prevIsSafe

		if prevIsCounter && currentIsCounter {
			return "", ErrInvalidString
		}

		if _, currentIsAllowToSafe := allowToSafe[currentStr]; prevIsSafe && !currentIsAllowToSafe {
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
	}
	return strBuilder.String(), nil
}
