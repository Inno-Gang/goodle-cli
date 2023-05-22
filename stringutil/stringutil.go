package stringutil

import (
	"fmt"
	"github.com/Inno-Gang/goodle-cli/icon"
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

func Trim(s string, max int) string {
	const ellipsis = icon.Ellipsis

	if len(s)-len(ellipsis) >= max {
		return s[:max-len(ellipsis)] + icon.Ellipsis
	}

	return s
}
