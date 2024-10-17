package helpers

import "unicode/utf8"

// Levenshtein computes the minimum number of single-character edits (insertions,
// deletions, or substitutions) required to change one string into another.
func Levenshtein(str1, str2 []rune) int {
	str1len := utf8.RuneCountInString(string(str1)) // Get the length of the first string.
	str2len := utf8.RuneCountInString(string(str2)) // Get the length of the second string.

	// Initialize a column to store distances for comparison.
	column := make([]int, str1len+1)

	// Set the initial values for the column (representing insertions).
	for y := 1; y <= str1len; y++ {
		column[y] = y
	}

	// Iterate through both strings to compute the Levenshtein distance.
	for x := 1; x <= str2len; x++ {
		column[0] = x
		lastKey := x - 1

		for y := 1; y <= str1len; y++ {
			oldKey := column[y]

			// Determine if the characters differ (increment by 1 if so).
			incr := 0
			if str1[y-1] != str2[x-1] {
				incr = 1
			}

			// Compute the minimum cost among insert, delete, or replace operations.
			column[y] = min(column[y]+1, column[y-1]+1, lastKey+incr)
			lastKey = oldKey
		}
	}
	return column[str1len] // Return the final Levenshtein distance.
}
