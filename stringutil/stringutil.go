package stringutil

import (
	"fmt"
	"strings"
)

func Capitalize(s string) string {
	if len(s) == 0 {
		return s
	}

	return strings.ToUpper(string(s[0])) + s[1:]
}

func Quantify(number int, singular, plural string) string {
	var noun string
	if number == 1 {
		noun = singular
	} else {
		noun = plural
	}

	return fmt.Sprintf("%d %s", number, noun)
}
