package gotodo

import (
	"math"
	"strings"
)

var letters = [...]string{
	"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M",
	"N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}

// letterIndex returns the 0 based index for an uppercased letter or -1
func letterIndex(char string) int {
	for idx, val := range letters {
		if val == char {
			return idx
		}
	}

	return -1
}

// isCompleteToken determines whether or not a token is a todo.txt complete flag
func isCompleteToken(token string) bool {
	return len(token) == 1 && string(token[0]) == "x"
}

// isPriorityToken determines whether or not a token is a valid todo.txt priority
func isPriorityToken(arg string) bool {
	// We need at least 3 characters in a priority- open paren, close paren, priority letter(s)
	if len(arg) < 3 {
		return false
	}

	start := 0
	end := len(arg) - 1
	if string(arg[start]) != "(" || string(arg[end]) != ")" {
		return false
	}

	rest := arg[start+1 : end]

	return IsPriorityString(rest)
}

// IsPriorityString determines whether or not a string is a valid todo.txt priority
func IsPriorityString(arg string) bool {
	// Do base-26 math to determine a priority score, where A is 1, AA is 27, AAA is 677, etc
	// If any character is not a capital letter, invalidate the score and return 0
	for _, char := range arg {
		value := letterIndex(strings.ToUpper(string(char)))
		if value == -1 {
			return false
		}
	}

	return true
}

// parsePriority converts a priority string into an integer score
func parsePriority(arg string) int {
	total := 0

	// Do base-26 math to determine a priority score, where A is 1, AA is 27, AAA is 677, etc
	// If any character is not a capital letter, invalidate the score and return 0
	for idx, char := range arg {
		pos := len(arg) - 1 - idx
		char := string(char)

		value := letterIndex(strings.ToUpper(string(char)))
		if value == -1 {
			return 0
		}

		oneBased := value + 1
		subtotal := math.Pow(float64(26), float64(pos)) * float64(oneBased)
		total += int(subtotal)
	}

	return total
}

// unparsePriority converts a base-26 number to a todo.txt priority score
func unparsePriority(priority int) string {
	quotient := priority
	remainder := 0
	scores := make([]string, 0)
	for {
		remainder = quotient % 26
		quotient = quotient / 26
		scores = append(scores, letters[remainder-1])
		if quotient == 0 {
			break
		}
	}

	for i := len(scores)/2 - 1; i >= 0; i-- {
		j := len(scores) - 1 - i
		scores[i], scores[j] = scores[j], scores[i]
	}

	return strings.Join(scores, "")
}

// parseDate determines whether or not a string is formatted according to TimeFormat.
// It returns a valid NullTime if the string is formatted properly, or invalid NullTime
func parseDate(arg string) NullTime {
	return NewNullTime(string(arg))
}

// parseTags extracts todo.txt projects and contexts from a slice of strings
func parseTags(parts []string) (Tags, Tags) {
	projects := make(Tags)
	contexts := make(Tags)
	var elem void

	for _, part := range parts {
		if len(part) == 0 {
			continue
		}
		switch string(part[0]) {
		case "+":
			projects[part[1:]] = elem
		case "@":
			contexts[part[1:]] = elem
		}
	}

	return projects, contexts
}
