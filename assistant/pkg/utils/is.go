package utils

import (
	"net"
	"net/url"
	"regexp"
	"strings"
)

var (
	reNumber               = regexp.MustCompile(`^[0-9]+$`)
	reLetter               = regexp.MustCompile(`^[a-zA-Z]+$`)
	reLetterOrNumber       = regexp.MustCompile(`^[a-zA-Z0-9]+$`)
	reLetterNumberOrSymbol = regexp.MustCompile(`^[a-zA-Z0-9_,.\-\/+:!@#$%^&*\[\](){}<> ]+$`)
	reChinese              = regexp.MustCompile(`^[\p{Han}]+$`)
	reEmail                = regexp.MustCompile(`^[\w._%+\-]+@[\w.\-]+\.[a-zA-Z]{2,}$`)
	rePhoneCN              = regexp.MustCompile(`^1(3\d|4[5-9]|5[0-35-9]|6[2567]|7[0-8]|8\d|9[0-35-9])\d{8}$`)
)

func IsNumber(s string) bool {
	return reNumber.MatchString(s)
}

func IsLetter(s string) bool {
	return reLetter.MatchString(s)
}

func IsLetterOrNumber(s string) bool {
	return reLetterOrNumber.MatchString(s)
}

func IsLetterOrNumberOrSymbol(s string) bool {
	return reLetterNumberOrSymbol.MatchString(s)
}

func IsChinese(s string) bool {
	return reChinese.MatchString(s)
}

func IsEmail(s string) bool {
	return reEmail.MatchString(s)
}

func IsPhone(s string) bool {
	return rePhoneCN.MatchString(s)
}

func IsURL(s string) bool {
	if strings.ContainsAny(s, " \t\r\n") {
		return false
	}
	u, err := url.Parse(s)
	if err != nil {
		return false
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return false
	}
	return u.Host != ""
}

func IsIPv4(s string) bool {
	ip := net.ParseIP(s)
	return ip != nil && ip.To4() != nil && strings.Count(s, ".") == 3
}

func IsIPv6(s string) bool {
	ip := net.ParseIP(s)
	return ip != nil && strings.Contains(s, ":")
}

func IsIP(s string) bool {
	return IsIPv4(s) || IsIPv6(s)
}
