package strset

import (
	"strconv"
	"strings"
	"unsafe"
)

func ToBytes(s string) []byte {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}

func BytesToString(b []byte) string {
	return unsafe.String(unsafe.SliceData(b), len(b))
}

func IsEmpty(s string) bool {
	return len(strings.TrimSpace(s)) == 0
}

func IsSliceEmpty(vals ...string) bool {
	for _, v := range vals {
		if IsEmpty(v) {
			return true
		}
	}
	return false
}

var replacer = strings.NewReplacer("\t", "", "\n", " ", "\r", "", "\x00", "")

func RemoveSpecialChars(s string) string {
	var sb strings.Builder
	sb.Grow(len(s))
	sb.WriteString(replacer.Replace(s))
	return sb.String()
}

func ToInt(str string) int {
	i, _ := strconv.Atoi(str)
	return i
}

func GetDigits(str string) string {
	var result strings.Builder
	for _, ch := range str {
		if ch >= '0' && ch <= '9' {
			result.WriteRune(ch)
		}
	}
	return result.String()
}

func IntToStr(i int) string {
	return strconv.Itoa(i)
}
