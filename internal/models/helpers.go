package models

import (
	"regexp"
	"strings"
)

func CreateSlugName(input string, maxLength int) string {
	formatted := strings.ToLower(input)
	re := regexp.MustCompile(`[^a-zA-Z0-9\-\ ]+`)
	formatted = re.ReplaceAllString(formatted, "")

	words := strings.Split(formatted, " ")

	// Rebuild string without exceeding maxLength,
	// and avoid cutting off words
	var result []string
	currentLength := 0

	for _, word := range words {
		wordLength := len(word)

		if currentLength+wordLength+1 > maxLength {
			break
		}

		result = append(result, word)
		// include the space between words in count
		currentLength += wordLength + 1
	}

	return strings.Join(result, "-")
}
