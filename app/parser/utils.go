package parser

import "regexp"

func StripColorCodes(input string) string {
	regex := regexp.MustCompile(`\^\d+`)
	return regex.ReplaceAllString(input, "")
}
