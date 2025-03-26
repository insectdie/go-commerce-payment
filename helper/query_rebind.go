package helper

import (
	"fmt"
	"strings"
)

func RebindQuery(query string) string {
	var count int
	// Replace each occurrence of "?" with "$<number>"
	var builder strings.Builder
	for _, char := range query {
		if char == '?' {
			count++
			builder.WriteString(fmt.Sprintf("$%d", count))
		} else {
			builder.WriteRune(char)
		}
	}

	return builder.String()
}
