package utils

import "regexp"

var (
	reNumber              = regexp.MustCompile(`^[0-9]+$`)
	reLetter              = regexp.MustCompile(`^[a-zA-Z]+$`)
	reLetterOrNumber      = regexp.MustCompile(`^[a-zA-Z0-9]+$`)
	reLetterOrNumberOrSym = regexp.MustCompile(`^[a-zA-Z0-9_,.\-\/+:!@#$%^&*\[\](){}<> ]+$`)
	reChinese             = regexp.MustCompile(`^[\p{Han}]+$`)
	reEmail               = regexp.MustCompile(`^[\w._%+\-]+@[\w.\-]+\.[a-zA-Z]{2,}$`)
	rePhone               = regexp.MustCompile(`^1(3\d|4[5-9]|5[0-35-9]|6[2567]|7[0-8]|8\d|9[0-35-9])\d{8}$`)
	reURL                 = regexp.MustCompile("^(http:|https:|www.|file://).*")
	reIPv4                = regexp.MustCompile(`^((2[0-4]\d|25[0-5]|[01]?\d\d?)\.){3}(2[0-4]\d|25[0-5]|[01]?\d\d?)$`)
	reIPv6                = regexp.MustCompile(`^([0-9A-Fa-f]{1,4}:){7}[0-9A-Fa-f]{1,4}$`)
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
	return reLetterOrNumberOrSym.MatchString(s)
}

func IsChinese(s string) bool {
	return reChinese.MatchString(s)
}

func IsEmail(s string) bool {
	return reEmail.MatchString(s)
}

func IsPhone(s string) bool {
	return rePhone.MatchString(s)
}

func IsURL(s string) bool {
	return reURL.MatchString(s)
}

func IsIPv4(s string) bool {
	return reIPv4.MatchString(s)
}

func IsIPv6(s string) bool {
	return reIPv6.MatchString(s)
}

func IsIP(s string) bool {
	return IsIPv4(s) || IsIPv6(s)
}
