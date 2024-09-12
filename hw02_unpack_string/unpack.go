package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode/utf8"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(in string) (string, error) {
	var strBuilder strings.Builder
	var prevRune rune
	var prevIsCounter bool
	var prevIsSafe bool
	var runeCount int = utf8.RuneCountInString(in)
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
