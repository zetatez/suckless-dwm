package utils

import "bytes"

func detectEncoding(content []byte) string {
	if len(content) == 0 {
		return "empty"
	}
	if bytes.HasPrefix(content, []byte{0xEF, 0xBB, 0xBF}) {
		return "utf-8-bom"
	}
	if bytes.HasPrefix(content, []byte{0xFF, 0xFE}) {
		return "utf-16-le"
	}
	if bytes.HasPrefix(content, []byte{0xFE, 0xFF}) {
		return "utf-16-be"
	}
	return "utf-8"
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
