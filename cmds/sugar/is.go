package utils

import (
	"regexp"
	"unicode/utf8"
)

func IsNumber(s string) bool {
	return regexp.MustCompile(`[0-9]+`).MatchString(s)
}

func IsLetter(s string) bool {
	return regexp.MustCompile(`[a-zA-Z]+`).MatchString(s)
}

func IsLetterOrNumber(s string) bool {
	return regexp.MustCompile(`[a-zA-Z0-9]+`).MatchString(s)
}

func IsLetterOrNumberOrSymbol(s string) bool {
	return regexp.MustCompile(`[a-zA-Z0-9_,.-/+:!@#$%^&*[](){}<> ]+`).MatchString(s)
}

func IsEmoji(s string) bool {
	return regexp.MustCompile(`[\u4e00-\u9fa5]+`).MatchString(s) && utf8.RuneCountInString(s) == 1
}

func IsChinese(s string) bool {
	return regexp.MustCompile(`[\u4e00-\u9fa5]+`).MatchString(s)
}
