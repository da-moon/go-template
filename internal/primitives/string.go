package primitives

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net"
	"net/url"
	"path"
	"path/filepath"
	"reflect"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"
	"unsafe"

	stacktrace "github.com/palantir/stacktrace"
)

var camelingRegex = regexp.MustCompile("[0-9A-Za-z]+")

// HasPrefix ...
func HasPrefix(s string, prefix string) bool {
	if runtime.GOOS == "windows" {
		return strings.HasPrefix(strings.ToLower(s), strings.ToLower(prefix))
	}
	return strings.HasPrefix(s, prefix)
}

// HasSuffix ...
func HasSuffix(s string, suffix string) bool {
	if runtime.GOOS == "windows" {
		return strings.HasSuffix(strings.ToLower(s), strings.ToLower(suffix))
	}
	return strings.HasSuffix(s, suffix)
}

// IsStringEqual ...
func IsStringEqual(s1 string, s2 string) bool {
	if runtime.GOOS == "windows" {
		return strings.EqualFold(s1, s2)
	}
	return s1 == s2
}

func NormalizePath(p string) string {
	lowecased := strings.ToLower(p)
	toslash := filepath.ToSlash(lowecased)
	cleaned := path.Clean(toslash)
	trimmed := strings.Trim(cleaned, "/.")
	return trimmed
}

// PathJoin ...
func PathJoin(elem ...string) string {
	trailingSlash := ""
	if len(elem) > 0 {
		if HasSuffix(elem[len(elem)-1], "/") {
			trailingSlash = "/"
		}
	}
	return path.Join(elem...) + trailingSlash
}

// StringFromByteSlice converts a slice of bytes into a string without performing a copy.
// This is an unsafe operation and may lead to problems if the bytes
// passed as argument are changed while the string is used.
func StringFromByteSlice(bytes []byte) string {
	hdr := *(*reflect.SliceHeader)(unsafe.Pointer(&bytes))
	return *(*string)(unsafe.Pointer(&reflect.StringHeader{
		Data: hdr.Data,
		Len:  hdr.Len,
	}))
}

// ReverseString returns the string reversed
func ReverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// LowerFirst returns a string with the first character lowercased
func LowerFirst(s string) string {
	if s == "" {
		return ""
	}
	r, n := utf8.DecodeRuneInString(s)
	return string(unicode.ToLower(r)) + s[n:]
}

// UpperFirst returns a string with the first character uppercased
func UpperFirst(s string) string {
	if s == "" {
		return ""
	}
	r, n := utf8.DecodeRuneInString(s)
	return string(unicode.ToUpper(r)) + s[n:]
}

// ToSnakeCase converts a string to snake_case
func ToSnakeCase(s string) string {
	s = strings.Replace(strings.Replace(s, "-", "_", -1), " ", "_", -1)
	runes := []rune(s)
	length := len(runes)
	var out []rune
	for i := 0; i < length; i++ {
		if i > 0 && unicode.IsUpper(runes[i]) && ((i+1 < length && unicode.IsLower(runes[i+1])) || unicode.IsLower(runes[i-1])) {
			out = append(out, '_')
		}
		out = append(out, unicode.ToLower(runes[i]))
	}
	return string(out)
}

// ToCamelCase converts a string to camelCase
func ToCamelCase(s string) string {
	byteSrc := []byte(s)
	chunks := camelingRegex.FindAll(byteSrc, -1)
	for idx, val := range chunks {
		val = bytes.ToLower(val)
		if idx > 0 {
			chunks[idx] = bytes.Title(val)
		}
	}
	return string(bytes.Join(chunks, nil))
}

// ToPascalCase converts a string to PascalCase
func ToPascalCase(s string) string {
	byteSrc := []byte(s)
	chunks := camelingRegex.FindAll(byteSrc, -1)
	for idx, val := range chunks {
		chunks[idx] = bytes.Title(val)
	}
	return string(bytes.Join(chunks, nil))
}

// ToKebabCase converts a string to kebab-case
func ToKebabCase(s string) string {
	s = strings.Replace(strings.Replace(s, "_", "-", -1), " ", "-", -1)
	runes := []rune(s)
	length := len(runes)
	var out []rune
	for i := 0; i < length; i++ {
		if i > 0 && unicode.IsUpper(runes[i]) && ((i+1 < length && unicode.IsLower(runes[i+1])) || unicode.IsLower(runes[i-1])) {
			out = append(out, '-')
		}
		out = append(out, unicode.ToLower(runes[i]))
	}
	return string(out)
}

// ToInt returns an int from a string
func ToInt(s string) int {
	if strings.Contains(s, ".") {
		s = strings.Split(s, ".")[0]
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return v
}

// ToInt64 returns an int64 from a string
func ToInt64(s string) int64 {
	if strings.Contains(s, ".") {
		s = strings.Split(s, ".")[0]
	}
	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0
	}
	return v
}

// ToUint returns a uint from a string
func ToUint(s string) uint {
	if strings.Contains(s, ".") {
		s = strings.Split(s, ".")[0]
	}
	v, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		return uint(0)
	}
	return uint(v)
}

// ToBase64 encodes a string in base64
func ToBase64(s string) string {
	return base64.StdEncoding.EncodeToString([]byte(s))
}

// FromBase64 decodes a string from base64
func FromBase64(s string) string {
	str, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return ""
	}
	return string(str)
}

// IsBool checks if a string is a boolean
func IsBool(s string) bool {
	_, err := strconv.ParseBool(s)
	return err == nil
}

// IsEmail returns true if the provided string is a valid email address
func IsEmail(s string) bool {
	emailRegexp := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	return emailRegexp.MatchString(s)
}

// IsURL returns true if the provided string is a valid url
func IsURL(s string) bool {
	_, err := url.ParseRequestURI(s)
	return err == nil
}

// IsJSON returns true if the provided string is a valid JSON document
func IsJSON(s string) bool {
	var js json.RawMessage
	return json.Unmarshal([]byte(s), &js) == nil
}

// IsIP returns true if the provided string is a valid IPv4
func IsIP(s string) bool {
	ip := net.ParseIP(s)
	return ip != nil && strings.Count(s, ".") == 3
}

// IsHexColor returns true if the provided string is a valid HEX color
func IsHexColor(s string) bool {
	if s == "" {
		return false
	}
	if s[0] == '#' {
		s = s[1:]
	}
	if len(s) != 3 && len(s) != 6 {
		return false
	}
	for _, c := range s {
		if ('F' < c || c < 'A') && ('f' < c || c < 'a') && ('9' < c || c < '0') {
			return false
		}
	}
	return true
}

// IsRGB returns true if the provided string is a valid RGB color
func IsRGB(s string) bool {
	if s == "" || len(s) < 10 {
		return false
	}
	if s[0:4] != "rgb(" || s[len(s)-1] != ')' {
		return false
	}
	s = s[4 : len(s)-1]
	s = strings.TrimSpace(s)
	if strings.Count(s, ",") != 2 {
		return false
	}
	for _, p := range strings.Split(s, ",") {
		if len(p) > 1 && p[0] == '0' {
			return false
		}
		p = strings.TrimSpace(p)
		if i, e := strconv.Atoi(p); (255 < i || i < 0) || e != nil {
			return false
		}
	}
	return true
}

// IsCreditCard returns true if the provided string is a valid credit card
func IsCreditCard(s string) bool {
	rxCreditCard := regexp.MustCompile(`^(?:4[0-9]{12}(?:[0-9]{3})?|5[1-5][0-9]{14}|6(?:011|5[0-9][0-9])[0-9]{12}|3[47][0-9]{13}|3(?:0[0-5]|[68][0-9])[0-9]{11}|(?:2131|1800|35\\d{3})\\d{11})$`)
	r, _ := regexp.Compile(`[^0-9]+`)
	sanitized := r.ReplaceAll([]byte(s), []byte(""))
	if !rxCreditCard.MatchString(string(sanitized)) {
		return false
	}
	var sum int
	var digit string
	var tmpNum int
	var shouldDouble bool
	for i := len(sanitized) - 1; i >= 0; i-- {
		digit = string(sanitized[i:(i + 1)])
		tmpNum = ToInt(digit)
		if shouldDouble {
			tmpNum *= 2
			if tmpNum >= 10 {
				sum += ((tmpNum % 10) + 1)
			} else {
				sum += tmpNum
			}
		} else {
			sum += tmpNum
		}
		shouldDouble = !shouldDouble
	}
	return sum%10 == 0
}

// IsOnlyDigits returns true if the provided string is composed only by numbers
func IsOnlyDigits(s string) bool {
	r := regexp.MustCompile("^[0-9]$")
	return r.MatchString(s)
}

// IsOnlyLetters returns true if the provided string is composed only by letters
func IsOnlyLetters(s string) bool {
	r := regexp.MustCompile("^[A-Za-z]+$")
	return r.MatchString(s)
}

// IsOnlyAlphaNumeric returns true if the provided string is composed only by letters and numbers
func IsOnlyAlphaNumeric(s string) bool {
	r := regexp.MustCompile("^[A-Za-z0-9]+$")
	return r.MatchString(s)
}

// Unique a unique token based on current time and crypted in sha256
func Unique() string {
	t := strconv.Itoa(int(time.Now().UnixNano()))
	return fmt.Sprintf("%x", sha256.Sum256([]byte("b33f"+t)))
}

// GetSizeString ...
func IntToFileSizeString(size, unit int) string {
	result := strconv.Itoa(size)
	if unit < Ki {
		if unit*size >= Ki {
			return ""
		}
		result = strconv.Itoa(size*unit) + " Bytes"
	}
	if unit >= Ki && unit < Mi {
		result = result + " Kilobytes"
	}
	if unit >= Mi {
		result = result + " Megabytes"
	}
	return result
}

// FileSizeStringToInt ...
func FileSizeStringToInt(s string) (int64, error) {
	ss := []byte(strings.ToUpper(s))
	if !(strings.Contains(string(ss), "K") || strings.Contains(string(ss), "KB") ||
		strings.Contains(string(ss), "M") || strings.Contains(string(ss), "MB") ||
		strings.Contains(string(ss), "G") || strings.Contains(string(ss), "GB") ||
		strings.Contains(string(ss), "T") || strings.Contains(string(ss), "TB")) {
		return -1, stacktrace.NewError("wrong format for input string")
	}
	var unit int64 = 1
	p, _ := strconv.Atoi(string(ss[:len(ss)-1]))
	unitstr := string(ss[len(ss)-1])
	if ss[len(ss)-1] == 'B' {
		p, _ = strconv.Atoi(string(ss[:len(ss)-2]))
		unitstr = string(ss[len(ss)-2:])
	}
	switch unitstr {
	default:
		// fallthrough
	case "T", "TB":
		unit *= 1024
		fallthrough
	case "G", "GB":
		unit *= 1024
		fallthrough
	case "M", "MB":
		unit *= 1024
		fallthrough
	case "K", "KB":
		unit *= 1024
	}
	return int64(p) * unit, nil
}
