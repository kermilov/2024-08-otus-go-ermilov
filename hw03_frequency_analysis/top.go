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
		_, isExist := countByWord[word]
		if !isExist {
			wordsSet = append(wordsSet, word)
		}
		countByWord[word]++
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
		result := make([]string, 10)
		copy(result, wordsSet[0:10])
		return result
	}
	return wordsSet
}
