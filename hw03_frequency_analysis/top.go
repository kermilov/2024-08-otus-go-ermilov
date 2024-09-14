package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

func Top10(in string) []string {
	if len(in) == 0 {
		return []string{}
	}
	words := strings.Fields(in)
	countByWord := make(map[string]int, len(words))
	wordsSet := make([]string, 0, len(words))
	for _, word := range words {
		count, isExist := countByWord[word]
		countByWord[word] = count + 1
		if !isExist {
			wordsSet = append(wordsSet, word)
		}
	}
	sort.Slice(wordsSet, func(i, j int) bool {
		left := wordsSet[i]
		right := wordsSet[j]
		if countByWord[left] == countByWord[right] {
			return left < right
		}
		return countByWord[left] > countByWord[right]
	})
	if len(wordsSet) > 10 {
		return wordsSet[0:10]
	}
	return wordsSet
}
