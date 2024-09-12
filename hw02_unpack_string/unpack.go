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
	var runeCount int = utf8.RuneCountInString(in)
	for i, currentRune := range in {
		currentStr := string(currentRune)

		repeatCount, err := strconv.Atoi(currentStr)
		currentIsCounter := err == nil	
		
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
			strBuilder.WriteString(strings.Repeat(string(prevRune), repeatCount))
		} else if !prevIsCounter {
			strBuilder.WriteRune(prevRune)
		}
		if i == runeCount-1 && !currentIsCounter {
			strBuilder.WriteString(currentStr)
		}
		prevRune = currentRune
		prevIsCounter = currentIsCounter
	}
	return strBuilder.String(), nil
}
