package utils

import "regexp"

func IsNumber(s string) bool {
	return regexp.MustCompile(`^[0-9]+$`).MatchString(s)
}

func IsLetter(s string) bool {
	return regexp.MustCompile(`^[a-zA-Z]+$`).MatchString(s)
}

func IsLetterOrNumber(s string) bool {
	return regexp.MustCompile(`^[a-zA-Z0-9]+$`).MatchString(s)
}

func IsLetterOrNumberOrSymbol(s string) bool {
	return regexp.MustCompile(`^[a-zA-Z0-9_,.\-\/+:!@#$%^&*\[\](){}<> ]+$`).MatchString(s)
}

func IsChinese(s string) bool {
	return regexp.MustCompile(`^[\p{Han}]+$`).MatchString(s)
}

func IsEmail(s string) bool {
	return regexp.MustCompile(`^[\w._%+\-]+@[\w.\-]+\.[a-zA-Z]{2,}$`).MatchString(s)
}

func IsPhone(s string) bool {
	return regexp.MustCompile(`^1(3\d|4[5-9]|5[0-35-9]|6[2567]|7[0-8]|8\d|9[0-35-9])\d{8}$`).MatchString(s)
}

func IsURL(s string) bool {
	return regexp.MustCompile("^(http:|https:|www.|file://).*").MatchString(s)
}

func IsIPv4(s string) bool {
	return regexp.MustCompile(`^((2[0-4]\d|25[0-5]|[01]?\d\d?)\.){3}(2[0-4]\d|25[0-5]|[01]?\d\d?)$`).MatchString(s)
}

func IsIPv6(s string) bool {
	return regexp.MustCompile(`^([0-9A-Fa-f]{1,4}:){7}[0-9A-Fa-f]{1,4}$`).MatchString(s)
}

func IsIP(s string) bool {
	return IsIPv4(s) || IsIPv6(s)
}
