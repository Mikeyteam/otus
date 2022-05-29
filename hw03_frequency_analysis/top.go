package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
)

// Top10 finds top count words in text.
func Top10(s string) []string {
	maxCountWords := 10
	result := sortingMap(findRepeatWords(cleanText(s)))
	if totalWords := len(result); totalWords < maxCountWords {
		maxCountWords = totalWords
	}

	return result[:maxCountWords]
}

// cleanText delete from text special char, change this on space.
func cleanText(text string) []string {
	prepareStr := strings.Fields(text)
	str := strings.Join(prepareStr[:], " ")
	replacer := regexp.MustCompile(`[/s]+`)
	replaceStr := replacer.ReplaceAllString(str, " ")
	result := strings.Split(strings.TrimSpace(replaceStr), " ")

	return result

}

// findRepeatWords find count repeat word in text.
func findRepeatWords(text []string) map[string]int {
	words := make(map[string]int)
	for _, word := range text {
		words[word]++
	}

	return words
}

// sortingMap sorting map and change to slice.
func sortingMap(list map[string]int) []string {
	slice := make([]string, 0, len(list))
	for word := range list {
		slice = append(slice, word)
	}

	sort.Slice(slice, func(i, j int) bool {
		return list[slice[i]] > list[slice[j]]
	})

	return slice
}
