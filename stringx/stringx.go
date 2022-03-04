package stringx

import (
	"bytes"
	crrand "crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"html"
	"math"
	"math/rand"
	"net/url"
	"path"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"
)

var (
	rxNumber = regexp.MustCompile(Numeric)
	rxEmail  = regexp.MustCompile(EMAIL)
	rxUrl    = regexp.MustCompile(URLPath)
	rxFloat  = regexp.MustCompile(Float)
	rxHost   = regexp.MustCompile(Host)
	rxPhone  = regexp.MustCompile(PHONE)
)

//截取字符串 start 起点下标 end 终点下标(不包括)
func Substr(str string, start int, end int) string {
	rs := []rune(str)
	length := len(rs)

	if start < 0 || end > length || (end >= 0 && start > end) {
		panic("out of range")
	}
	if len(str) == 0 && (start > 0 || end > 0) {
		panic("out of range")
	}

	if end < 0 {
		end += length
	}
	return string(rs[start:end])
}

func Substrs(str string, minLen int) []string {
	// Convert to a rune to safely handle unicode
	r_str := []rune(str)
	totLen := len(r_str)
	var subsCount, length int

	if minLen <= 1 {
		minLen = 1
		length = totLen
		subsCount = factorial(length)
	} else {
		length = totLen - minLen + 1
		subsCount = factorial(length)
	}

	subs := make([]string, subsCount)

	next := 0
	for i := 0; i < length; i++ {
		for j := i; j < totLen; j++ {
			lenbtw := j - i + 1
			if lenbtw >= minLen {
				subs[next] = Substr(str, i, lenbtw)
				next++
			}
		}
	}
	return subs
}

func RuneSlices(r []rune, minLen int) [][]rune {
	totLen := len(r)
	var subsCount, length int

	if minLen <= 1 {
		minLen = 1
		length = totLen
		subsCount = factorial(length)
	} else {
		length = totLen - minLen + 1
		subsCount = factorial(length)
	}

	subs := make([][]rune, subsCount)

	next := 0
	for i := 0; i < length; i++ {
		for j := i; j < totLen; j++ {
			if j-i+1 >= minLen {
				subs[next] = r[i : j+1]
				next++
			}
		}
	}
	return subs
}

type sortRunes []rune

func (s sortRunes) Less(i, j int) bool {
	return s[i] < s[j]
}

// Swap the runes at indexes i and j
func (s sortRunes) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s sortRunes) Len() int {
	return len(s)
}

func SortString(str string) string {
	r := []rune(str)
	sort.Sort(sortRunes(r))
	return string(r)
}

// Cache factorials
var factorials = make(map[int]int)

//阶乘
func factorial(n int) int {
	if fac, ok := factorials[n]; ok {
		return fac
	}

	fac := n
	for m := fac - 1; m > 0; m-- {
		fac += m
	}
	factorials[n] = fac
	return fac
}

// Contains check if the string contains the substring.
func Contains(str, substring string) bool {
	return strings.Contains(str, substring)
}

// Matches check if string matches the pattern (pattern is regular expression)
// In case of error return false
func Matches(str, pattern string) bool {
	match, _ := regexp.MatchString(pattern, str)
	return match
}

// LeftTrim trim characters from the left-side of the input.
// If second argument is empty, it's will be remove leading spaces.
func LeftTrim(str, chars string) string {
	if chars == "" {
		return strings.TrimLeftFunc(str, unicode.IsSpace)
	}
	r, _ := regexp.Compile("^[" + chars + "]+")
	return r.ReplaceAllString(str, "")
}

// RightTrim trim characters from the right-side of the input.
// If second argument is empty, it's will be remove spaces.
func RightTrim(str, chars string) string {
	if chars == "" {
		return strings.TrimRightFunc(str, unicode.IsSpace)
	}
	r, _ := regexp.Compile("[" + chars + "]+$")
	return r.ReplaceAllString(str, "")
}

// Trim trim characters from both sides of the input.
// If second argument is empty, it's will be remove spaces.
func Trim(str, chars string) string {
	return LeftTrim(RightTrim(str, chars), chars)
}

// WhiteList remove characters that do not appear in the whitelist.
func WhiteList(str, chars string) string {
	pattern := "[^" + chars + "]+"
	r, _ := regexp.Compile(pattern)
	return r.ReplaceAllString(str, "")
}

// BlackList remove characters that appear in the blacklist.
func BlackList(str, chars string) string {
	pattern := "[" + chars + "]+"
	r, _ := regexp.Compile(pattern)
	return r.ReplaceAllString(str, "")
}

// StripLow remove characters with a numerical value < 32 and 127, mostly control characters.
// If keep_new_lines is true, newline characters are preserved (\n and \r, hex 0xA and 0xD).
func StripLow(str string, keepNewLines bool) string {
	chars := ""
	if keepNewLines {
		chars = "\x00-\x09\x0B\x0C\x0E-\x1F\x7F"
	} else {
		chars = "\x00-\x1F\x7F"
	}
	return BlackList(str, chars)
}

// ReplacePattern replace regular expression pattern in string
func ReplacePattern(str, pattern, replace string) string {
	r, _ := regexp.Compile(pattern)
	return r.ReplaceAllString(str, replace)
}

// Escape replace <, >, & and " with HTML entities.
var Escape = html.EscapeString

func addSegment(inrune, segment []rune) []rune {
	if len(segment) == 0 {
		return inrune
	}
	if len(inrune) != 0 {
		inrune = append(inrune, '_')
	}
	inrune = append(inrune, segment...)
	return inrune
}

// UnderscoreToCamelCase converts from underscore separated form to camel case form.
// Ex.: my_func => MyFunc
func UnderscoreToCamelCase(s string) string {
	return strings.Replace(strings.Title(strings.Replace(strings.ToLower(s), "_", " ", -1)), " ", "", -1)
}

// CamelCaseToUnderscore converts from camel case form to underscore separated form.
// Ex.: MyFunc => my_func
func CamelCaseToUnderscore(str string) string {
	var output []rune
	var segment []rune
	for _, r := range str {
		if !unicode.IsLower(r) && string(r) != "_" {
			output = addSegment(output, segment)
			segment = nil
		}
		segment = append(segment, unicode.ToLower(r))
	}
	output = addSegment(output, segment)
	return string(output)
}

// Reverse return reversed string
func Reverse(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < j; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}

// GetLines split string by "\n" and return array of lines
func GetLines(s string) []string {
	return strings.Split(s, "\n")
}

// GetLine return specified line of multiline string
func GetLine(s string, index int) (string, error) {
	lines := GetLines(s)
	if index < 0 || index >= len(lines) {
		return "", errors.New("line index out of bounds")
	}
	return lines[index], nil
}

// RemoveTags remove all tags from HTML string
func RemoveTags(s string) string {
	return ReplacePattern(s, "<[^>]*>", "")
}

// SafeFileName return safe string that can be used in file names
func SafeFileName(str string) string {
	name := strings.ToLower(str)
	name = path.Clean(path.Base(name))
	name = strings.Trim(name, " ")
	separators, err := regexp.Compile(`[ &_=+:]`)
	if err == nil {
		name = separators.ReplaceAllString(name, "-")
	}
	legal, err := regexp.Compile(`[^[:alnum:]-.]`)
	if err == nil {
		name = legal.ReplaceAllString(name, "")
	}
	for strings.Contains(name, "--") {
		name = strings.Replace(name, "--", "-", -1)
	}
	return name
}

// NormalizeEmail canonicalize an email address.
// The local part of the email address is lowercased for all domains; the hostname is always lowercased and
// the local part of the email address is always lowercased for hosts that are known to be case-insensitive (currently only GMail).
// Normalization follows special rules for known providers: currently, GMail addresses have dots removed in the local part and
// are stripped of tags (e.g. some.one+tag@gmail.com becomes someone@gmail.com) and all @googlemail.com addresses are
// normalized to @gmail.com.
func NormalizeEmail(str string) (string, error) {
	if !IsEmail(str) {
		return "", fmt.Errorf("%s is not an email", str)
	}
	parts := strings.Split(str, "@")
	parts[0] = strings.ToLower(parts[0])
	parts[1] = strings.ToLower(parts[1])
	if parts[1] == "gmail.com" || parts[1] == "googlemail.com" {
		parts[1] = "gmail.com"
		parts[0] = strings.Split(ReplacePattern(parts[0], `\.`, ""), "+")[0]
	}
	return strings.Join(parts, "@"), nil
}

// Truncate a string to the closest length without breaking words.
func Truncate(str string, length int, ending string) string {
	var aftstr, befstr string
	if len(str) > length {
		words := strings.Fields(str)
		before, present := 0, 0
		for i := range words {
			befstr = aftstr
			before = present
			aftstr = aftstr + words[i] + " "
			present = len(aftstr)
			if present > length && i != 0 {
				if (length - before) < (present - length) {
					return Trim(befstr, " /\\.,\"'#!?&@+-") + ending
				}
				return Trim(aftstr, " /\\.,\"'#!?&@+-") + ending
			}
		}
	}

	return str
}

// PadLeft pad left side of string if size of string is less then indicated pad length
func PadLeft(str string, padStr string, padLen int) string {
	return buildPadStr(str, padStr, padLen, true, false)
}

// PadRight pad right side of string if size of string is less then indicated pad length
func PadRight(str string, padStr string, padLen int) string {
	return buildPadStr(str, padStr, padLen, false, true)
}

// PadBoth pad sides of string if size of string is less then indicated pad length
func PadBoth(str string, padStr string, padLen int) string {
	return buildPadStr(str, padStr, padLen, true, true)
}

// PadString either left, right or both sides, not the padding string can be unicode and more then one
// character
func buildPadStr(str string, padStr string, padLen int, padLeft bool, padRight bool) string {

	// When padded length is less then the current string size
	if padLen < utf8.RuneCountInString(str) {
		return str
	}

	padLen -= utf8.RuneCountInString(str)

	targetLen := padLen

	targetLenLeft := targetLen
	targetLenRight := targetLen
	if padLeft && padRight {
		targetLenLeft = padLen / 2
		targetLenRight = padLen - targetLenLeft
	}

	strToRepeatLen := utf8.RuneCountInString(padStr)

	repeatTimes := int(math.Ceil(float64(targetLen) / float64(strToRepeatLen)))
	repeatedString := strings.Repeat(padStr, repeatTimes)

	leftSide := ""
	if padLeft {
		leftSide = repeatedString[0:targetLenLeft]
	}

	rightSide := ""
	if padRight {
		rightSide = repeatedString[0:targetLenRight]
	}

	return leftSide + str + rightSide
}

// IsEmail check if the string is an email.
func IsPhone(str string) bool {
	// TODO uppercase letters are not supported
	return rxPhone.MatchString(str)
}

// IsEmail check if the string is an email.
func IsEmail(str string) bool {
	// TODO uppercase letters are not supported
	return rxEmail.MatchString(str)
}
func IsNumber(str string) bool {
	return rxNumber.MatchString(str)
}
func IsFloat(str string) bool {
	return rxFloat.MatchString(str)
}

func IsNumberMin(str string, min int) bool {
	if rxNumber.MatchString(str) {
		v, _ := strconv.Atoi(str)
		return v >= min
	}
	return false
}

func IsFloatMin(str string, min float64) bool {
	if rxFloat.MatchString(str) {
		v, _ := strconv.ParseFloat(str, 64)
		return v >= min
	}
	return false
}

func IsFloatGreat(str string, min float64) bool {
	if rxFloat.MatchString(str) {
		v, _ := strconv.ParseFloat(str, 64)
		return v > min
	}
	return false
}

func UpcaseFirst(str string) string {
	for i, v := range str {
		return string(unicode.ToUpper(v)) + str[i+1:]
	}
	return ""
}

func LowerFirst(str string) string {

	for i, v := range str {
		return string(unicode.ToLower(v)) + str[i+1:]
	}
	return ""
}

func IsFirstUp(str string) bool {
	for _, v := range str {
		if unicode.IsUpper(v) {
			return true
		}
		break
	}
	return false
}

//驼峰形式且首字母小写
func LowerTitle(str string) string {
	return LowerFirst(UnderscoreToCamelCase(str))
}

//以下代码来自xstrings  github.com/huandu/xstrings

const bufferMaxInitGrowSize = 2048

// Lazy initialize a buffer.
func allocBuffer(orig, cur string) *bytes.Buffer {
	output := &bytes.Buffer{}
	maxSize := len(orig) * 4

	// Avoid to reserve too much memory at once.
	if maxSize > bufferMaxInitGrowSize {
		maxSize = bufferMaxInitGrowSize
	}

	output.Grow(maxSize)
	output.WriteString(orig[:len(orig)-len(cur)])
	return output
}

// ToCamelCase can convert all lower case characters behind underscores
// to upper case character.
// Underscore character will be removed in result except following cases.
//     * More than 1 underscore.
//           "a__b" => "A_B"
//     * At the beginning of string.
//           "_a" => "_A"
//     * At the end of string.
//           "ab_" => "Ab_"
func ToCamelCase(str string) string {
	if len(str) == 0 {
		return ""
	}

	buf := &bytes.Buffer{}
	var r0, r1 rune
	var size int

	// leading '_' will appear in output.
	for len(str) > 0 {
		r0, size = utf8.DecodeRuneInString(str)
		str = str[size:]

		if r0 != '_' {
			break
		}

		buf.WriteRune(r0)
	}

	if len(str) == 0 {
		return buf.String()
	}

	r0 = unicode.ToUpper(r0)

	for len(str) > 0 {
		r1 = r0
		r0, size = utf8.DecodeRuneInString(str)
		str = str[size:]

		if r1 == '_' && r0 == '_' {
			buf.WriteRune(r1)
			continue
		}

		if r1 == '_' {
			r0 = unicode.ToUpper(r0)
		} else {
			r0 = unicode.ToLower(r0)
		}

		if r1 != '_' {
			buf.WriteRune(r1)
		}
	}

	buf.WriteRune(r0)
	return buf.String()
}

// ToSnakeCase can convert all upper case characters in a string to
// underscore format.
//
// Some samples.
//     "FirstName"  => "first_name"
//     "HTTPServer" => "http_server"
//     "NoHTTPS"    => "no_https"
//     "GO_PATH"    => "go_path"
//     "GO PATH"    => "go_path"      // space is converted to underscore.
//     "GO-PATH"    => "go_path"      // hyphen is converted to underscore.
func ToSnakeCase(str string) string {
	if len(str) == 0 {
		return ""
	}

	buf := &bytes.Buffer{}
	var prev, r0, r1 rune
	var size int

	r0 = '_'

	for len(str) > 0 {
		prev = r0
		r0, size = utf8.DecodeRuneInString(str)
		str = str[size:]

		switch {
		case r0 == utf8.RuneError:
			buf.WriteByte(byte(str[0]))

		case unicode.IsUpper(r0):
			if prev != '_' {
				buf.WriteRune('_')
			}

			buf.WriteRune(unicode.ToLower(r0))

			if len(str) == 0 {
				break
			}

			r0, size = utf8.DecodeRuneInString(str)
			str = str[size:]

			if !unicode.IsUpper(r0) {
				buf.WriteRune(r0)
				break
			}

			// find next non-upper-case character and insert `_` properly.
			// it's designed to convert `HTTPServer` to `http_server`.
			// if there are more than 2 adjacent upper case characters in a word,
			// treat them as an abbreviation plus a normal word.
			for len(str) > 0 {
				r1 = r0
				r0, size = utf8.DecodeRuneInString(str)
				str = str[size:]

				if r0 == utf8.RuneError {
					buf.WriteRune(unicode.ToLower(r1))
					buf.WriteByte(byte(str[0]))
					break
				}

				if !unicode.IsUpper(r0) {
					if r0 == '_' || r0 == ' ' || r0 == '-' {
						r0 = '_'

						buf.WriteRune(unicode.ToLower(r1))
					} else {
						buf.WriteRune('_')
						buf.WriteRune(unicode.ToLower(r1))
						buf.WriteRune(r0)
					}

					break
				}

				buf.WriteRune(unicode.ToLower(r1))
			}

			if len(str) == 0 || r0 == '_' {
				buf.WriteRune(unicode.ToLower(r0))
				break
			}

		default:
			if r0 == ' ' || r0 == '-' {
				r0 = '_'
			}

			buf.WriteRune(r0)
		}
	}

	return buf.String()
}

// SwapCase will swap characters case from upper to lower or lower to upper.
func SwapCase(str string) string {
	var r rune
	var size int

	buf := &bytes.Buffer{}

	for len(str) > 0 {
		r, size = utf8.DecodeRuneInString(str)

		switch {
		case unicode.IsUpper(r):
			buf.WriteRune(unicode.ToLower(r))

		case unicode.IsLower(r):
			buf.WriteRune(unicode.ToUpper(r))

		default:
			buf.WriteRune(r)
		}

		str = str[size:]
	}

	return buf.String()
}

// FirstRuneToUpper converts first rune to upper case if necessary.
func FirstRuneToUpper(str string) string {
	if str == "" {
		return str
	}

	r, size := utf8.DecodeRuneInString(str)

	if !unicode.IsLower(r) {
		return str
	}

	buf := &bytes.Buffer{}
	buf.WriteRune(unicode.ToUpper(r))
	buf.WriteString(str[size:])
	return buf.String()
}

// FirstRuneToLower converts first rune to lower case if necessary.
func FirstRuneToLower(str string) string {
	if str == "" {
		return str
	}

	r, size := utf8.DecodeRuneInString(str)

	if !unicode.IsUpper(r) {
		return str
	}

	buf := &bytes.Buffer{}
	buf.WriteRune(unicode.ToLower(r))
	buf.WriteString(str[size:])
	return buf.String()
}

// Shuffle randomizes runes in a string and returns the result.
// It uses default random source in `math/rand`.
func Shuffle(str string) string {
	if str == "" {
		return str
	}

	runes := []rune(str)
	index := 0

	for i := len(runes) - 1; i > 0; i-- {
		index = rand.Intn(i + 1)

		if i != index {
			runes[i], runes[index] = runes[index], runes[i]
		}
	}

	return string(runes)
}

// ShuffleSource randomizes runes in a string with given random source.
func ShuffleSource(str string, src rand.Source) string {
	if str == "" {
		return str
	}

	runes := []rune(str)
	index := 0
	r := rand.New(src)

	for i := len(runes) - 1; i > 0; i-- {
		index = r.Intn(i + 1)

		if i != index {
			runes[i], runes[index] = runes[index], runes[i]
		}
	}

	return string(runes)
}

// Successor returns the successor to string.
//
// If there is one alphanumeric rune is found in string, increase the rune by 1.
// If increment generates a "carry", the rune to the left of it is incremented.
// This process repeats until there is no carry, adding an additional rune if necessary.
//
// If there is no alphanumeric rune, the rightmost rune will be increased by 1
// regardless whether the result is a valid rune or not.
//
// Only following characters are alphanumeric.
//     * a - z
//     * A - Z
//     * 0 - 9
//
// Samples (borrowed from ruby's String#succ document):
//     "abcd"      => "abce"
//     "THX1138"   => "THX1139"
//     "<<koala>>" => "<<koalb>>"
//     "1999zzz"   => "2000aaa"
//     "ZZZ9999"   => "AAAA0000"
//     "***"       => "**+"
func Successor(str string) string {
	if str == "" {
		return str
	}

	var r rune
	var i int
	carry := ' '
	runes := []rune(str)
	l := len(runes)
	lastAlphanumeric := l

	for i = l - 1; i >= 0; i-- {
		r = runes[i]

		if ('a' <= r && r <= 'y') ||
			('A' <= r && r <= 'Y') ||
			('0' <= r && r <= '8') {
			runes[i]++
			carry = ' '
			lastAlphanumeric = i
			break
		}

		switch r {
		case 'z':
			runes[i] = 'a'
			carry = 'a'
			lastAlphanumeric = i

		case 'Z':
			runes[i] = 'A'
			carry = 'A'
			lastAlphanumeric = i

		case '9':
			runes[i] = '0'
			carry = '0'
			lastAlphanumeric = i
		}
	}

	// Needs to add one character for carry.
	if i < 0 && carry != ' ' {
		buf := &bytes.Buffer{}
		buf.Grow(l + 4) // Reserve enough space for write.

		if lastAlphanumeric != 0 {
			buf.WriteString(str[:lastAlphanumeric])
		}

		buf.WriteRune(carry)

		for _, r = range runes[lastAlphanumeric:] {
			buf.WriteRune(r)
		}

		return buf.String()
	}

	// No alphanumeric character. Simply increase last rune's value.
	if lastAlphanumeric == l {
		runes[l-1]++
	}

	return string(runes)
}

// Len returns str's utf8 rune length.
func Len(str string) int {
	return utf8.RuneCountInString(str)
}

// WordCount returns number of words in a string.
//
// Word is defined as a locale dependent string containing alphabetic characters,
// which may also contain but not start with `'` and `-` characters.
func WordCount(str string) int {
	var r rune
	var size, n int

	inWord := false

	for len(str) > 0 {
		r, size = utf8.DecodeRuneInString(str)

		switch {
		case isAlphabet(r):
			if !inWord {
				inWord = true
				n++
			}

		case inWord && (r == '\'' || r == '-'):
			// Still in word.

		default:
			inWord = false
		}

		str = str[size:]
	}

	return n
}

const minCJKCharacter = '\u3400'

// Checks r is a letter but not CJK character.
func isAlphabet(r rune) bool {
	if !unicode.IsLetter(r) {
		return false
	}

	switch {
	// Quick check for non-CJK character.
	case r < minCJKCharacter:
		return true

		// Common CJK characters.
	case r >= '\u4E00' && r <= '\u9FCC':
		return false

		// Rare CJK characters.
	case r >= '\u3400' && r <= '\u4D85':
		return false

		// Rare and historic CJK characters.
	case r >= '\U00020000' && r <= '\U0002B81D':
		return false
	}

	return true
}

// Width returns string width in monotype font.
// Multi-byte characters are usually twice the width of single byte characters.
//
// Algorithm comes from `mb_strwidth` in PHP.
// http://php.net/manual/en/function.mb-strwidth.php
func Width(str string) int {
	var r rune
	var size, n int

	for len(str) > 0 {
		r, size = utf8.DecodeRuneInString(str)
		n += RuneWidth(r)
		str = str[size:]
	}

	return n
}

// RuneWidth returns character width in monotype font.
// Multi-byte characters are usually twice the width of single byte characters.
//
// Algorithm comes from `mb_strwidth` in PHP.
// http://php.net/manual/en/function.mb-strwidth.php
func RuneWidth(r rune) int {
	switch {
	case r == utf8.RuneError || r < '\x20':
		return 0

	case '\x20' <= r && r < '\u2000':
		return 1

	case '\u2000' <= r && r < '\uFF61':
		return 2

	case '\uFF61' <= r && r < '\uFFA0':
		return 1

	case '\uFFA0' <= r:
		return 2
	}

	return 0
}

// ExpandTabs can expand tabs ('\t') rune in str to one or more spaces dpending on
// current column and tabSize.
// The column number is reset to zero after each newline ('\n') occurring in the str.
//
// ExpandTabs uses RuneWidth to decide rune's width.
// For example, CJK characters will be treated as two characters.
//
// If tabSize <= 0, ExpandTabs panics with error.
//
// Samples:
//     ExpandTabs("a\tbc\tdef\tghij\tk", 4) => "a   bc  def ghij    k"
//     ExpandTabs("abcdefg\thij\nk\tl", 4)  => "abcdefg hij\nk   l"
//     ExpandTabs("z中\t文\tw", 4)           => "z中 文  w"
func ExpandTabs(str string, tabSize int) string {
	if tabSize <= 0 {
		panic("tab size must be positive")
	}

	var r rune
	var i, size, column, expand int
	var output *bytes.Buffer

	orig := str

	for len(str) > 0 {
		r, size = utf8.DecodeRuneInString(str)

		if r == '\t' {
			expand = tabSize - column%tabSize

			if output == nil {
				output = allocBuffer(orig, str)
			}

			for i = 0; i < expand; i++ {
				output.WriteByte(byte(' '))
			}

			column += expand
		} else {
			if r == '\n' {
				column = 0
			} else {
				column += RuneWidth(r)
			}

			if output != nil {
				output.WriteRune(r)
			}
		}

		str = str[size:]
	}

	if output == nil {
		return orig
	}

	return output.String()
}

// LeftJustify returns a string with pad string at right side if str's rune length is smaller than length.
// If str's rune length is larger than length, str itself will be returned.
//
// If pad is an empty string, str will be returned.
//
// Samples:
//     LeftJustify("hello", 4, " ")    => "hello"
//     LeftJustify("hello", 10, " ")   => "hello     "
//     LeftJustify("hello", 10, "123") => "hello12312"
func LeftJustify(str string, length int, pad string) string {
	l := Len(str)

	if l >= length || pad == "" {
		return str
	}

	remains := length - l
	padLen := Len(pad)

	output := &bytes.Buffer{}
	output.Grow(len(str) + (remains/padLen+1)*len(pad))
	output.WriteString(str)
	writePadString(output, pad, padLen, remains)
	return output.String()
}

// RightJustify returns a string with pad string at left side if str's rune length is smaller than length.
// If str's rune length is larger than length, str itself will be returned.
//
// If pad is an empty string, str will be returned.
//
// Samples:
//     RightJustify("hello", 4, " ")    => "hello"
//     RightJustify("hello", 10, " ")   => "     hello"
//     RightJustify("hello", 10, "123") => "12312hello"
func RightJustify(str string, length int, pad string) string {
	l := Len(str)

	if l >= length || pad == "" {
		return str
	}

	remains := length - l
	padLen := Len(pad)

	output := &bytes.Buffer{}
	output.Grow(len(str) + (remains/padLen+1)*len(pad))
	writePadString(output, pad, padLen, remains)
	output.WriteString(str)
	return output.String()
}

// Center returns a string with pad string at both side if str's rune length is smaller than length.
// If str's rune length is larger than length, str itself will be returned.
//
// If pad is an empty string, str will be returned.
//
// Samples:
//     Center("hello", 4, " ")    => "hello"
//     Center("hello", 10, " ")   => "  hello   "
//     Center("hello", 10, "123") => "12hello123"
func Center(str string, length int, pad string) string {
	l := Len(str)

	if l >= length || pad == "" {
		return str
	}

	remains := length - l
	padLen := Len(pad)

	output := &bytes.Buffer{}
	output.Grow(len(str) + (remains/padLen+1)*len(pad))
	writePadString(output, pad, padLen, remains/2)
	output.WriteString(str)
	writePadString(output, pad, padLen, (remains+1)/2)
	return output.String()
}

func writePadString(output *bytes.Buffer, pad string, padLen, remains int) {
	var r rune
	var size int

	repeats := remains / padLen

	for i := 0; i < repeats; i++ {
		output.WriteString(pad)
	}

	remains = remains % padLen

	if remains != 0 {
		for i := 0; i < remains; i++ {
			r, size = utf8.DecodeRuneInString(pad)
			output.WriteRune(r)
			pad = pad[size:]
		}
	}
}

// Slice a string by rune.
//
// Start must satisfy 0 <= start <= rune length.
//
// End can be positive, zero or negative.
// If end >= 0, start and end must satisfy start <= end <= rune length.
// If end < 0, it means slice to the end of string.
//
// Otherwise, Slice will panic as out of range.
func Slice(str string, start, end int) string {
	var size, startPos, endPos int

	origin := str

	if start < 0 || end > len(str) || (end >= 0 && start > end) {
		panic("out of range")
	}

	if end >= 0 {
		end -= start
	}

	for start > 0 && len(str) > 0 {
		_, size = utf8.DecodeRuneInString(str)
		start--
		startPos += size
		str = str[size:]
	}

	if end < 0 {
		return origin[startPos:]
	}

	endPos = startPos

	for end > 0 && len(str) > 0 {
		_, size = utf8.DecodeRuneInString(str)
		end--
		endPos += size
		str = str[size:]
	}

	if len(str) == 0 && (start > 0 || end > 0) {
		panic("out of range")
	}

	return origin[startPos:endPos]
}

// Partition splits a string by sep into three parts.
// The return value is a slice of strings with head, match and tail.
//
// If str contains sep, for example "hello" and "l", Partition returns
//     "he", "l", "lo"
//
// If str doesn't contain sep, for example "hello" and "x", Partition returns
//     "hello", "", ""
func Partition(str, sep string) (head, match, tail string) {
	index := strings.Index(str, sep)

	if index == -1 {
		head = str
		return
	}

	head = str[:index]
	match = str[index : index+len(sep)]
	tail = str[index+len(sep):]
	return
}

// LastPartition splits a string by last instance of sep into three parts.
// The return value is a slice of strings with head, match and tail.
//
// If str contains sep, for example "hello" and "l", LastPartition returns
//     "hel", "l", "o"
//
// If str doesn't contain sep, for example "hello" and "x", LastPartition returns
//     "", "", "hello"
func LastPartition(str, sep string) (head, match, tail string) {
	index := strings.LastIndex(str, sep)

	if index == -1 {
		tail = str
		return
	}

	head = str[:index]
	match = str[index : index+len(sep)]
	tail = str[index+len(sep):]
	return
}

// Insert src into dst at given rune index.
// Index is counted by runes instead of bytes.
//
// If index is out of range of dst, panic with out of range.
func Insert(dst, src string, index int) string {
	return Slice(dst, 0, index) + src + Slice(dst, index, -1)
}

// Scrub scrubs invalid utf8 bytes with repl string.
// Adjacent invalid bytes are replaced only once.
func Scrub(str, repl string) string {
	var buf *bytes.Buffer
	var r rune
	var size, pos int
	var hasError bool

	origin := str

	for len(str) > 0 {
		r, size = utf8.DecodeRuneInString(str)

		if r == utf8.RuneError {
			if !hasError {
				if buf == nil {
					buf = &bytes.Buffer{}
				}

				buf.WriteString(origin[:pos])
				hasError = true
			}
		} else if hasError {
			hasError = false
			buf.WriteString(repl)

			origin = origin[pos:]
			pos = 0
		}

		pos += size
		str = str[size:]
	}

	if buf != nil {
		buf.WriteString(origin)
		return buf.String()
	}

	// No invalid byte.
	return origin
}

// WordSplit splits a string into words. Returns a slice of words.
// If there is no word in a string, return nil.
//
// Word is defined as a locale dependent string containing alphabetic characters,
// which may also contain but not start with `'` and `-` characters.
func WordSplit(str string) []string {
	var word string
	var words []string
	var r rune
	var size, pos int

	inWord := false

	for len(str) > 0 {
		r, size = utf8.DecodeRuneInString(str)

		switch {
		case isAlphabet(r):
			if !inWord {
				inWord = true
				word = str
				pos = 0
			}

		case inWord && (r == '\'' || r == '-'):
			// Still in word.

		default:
			if inWord {
				inWord = false
				words = append(words, word[:pos])
			}
		}

		pos += size
		str = str[size:]
	}

	if inWord {
		words = append(words, word[:pos])
	}

	return words
}

type runeRangeMap struct {
	FromLo rune // Lower bound of range map.
	FromHi rune // An inclusive higher bound of range map.
	ToLo   rune
	ToHi   rune
}

type runeDict struct {
	Dict [unicode.MaxASCII + 1]rune
}

type runeMap map[rune]rune

// Translator can translate string with pre-compiled from and to patterns.
// If a from/to pattern pair needs to be used more than once, it's recommended
// to create a Translator and reuse it.
type Translator struct {
	quickDict  *runeDict       // A quick dictionary to look up rune by index. Only availabe for latin runes.
	runeMap    runeMap         // Rune map for translation.
	ranges     []*runeRangeMap // Ranges of runes.
	mappedRune rune            // If mappedRune >= 0, all matched runes are translated to the mappedRune.
	reverted   bool            // If to pattern is empty, all matched characters will be deleted.
	hasPattern bool
}

// NewTranslator creates new Translator through a from/to pattern pair.
func NewTranslator(from, to string) *Translator {
	tr := &Translator{}

	if from == "" {
		return tr
	}

	reverted := from[0] == '^'
	deletion := len(to) == 0

	if reverted {
		from = from[1:]
	}

	var fromStart, fromEnd, fromRangeStep rune
	var toStart, toEnd, toRangeStep rune
	var fromRangeSize, toRangeSize rune
	var singleRunes []rune

	// Update the to rune range.
	updateRange := func() {
		// No more rune to read in the to rune pattern.
		if toEnd == utf8.RuneError {
			return
		}

		if toRangeStep == 0 {
			to, toStart, toEnd, toRangeStep = nextRuneRange(to, toEnd)
			return
		}

		// Current range is not empty. Consume 1 rune from start.
		if toStart != toEnd {
			toStart += toRangeStep
			return
		}

		// No more rune. Repeat the last rune.
		if to == "" {
			toEnd = utf8.RuneError
			return
		}

		// Both start and end are used. Read two more runes from the to pattern.
		to, toStart, toEnd, toRangeStep = nextRuneRange(to, utf8.RuneError)
	}

	if deletion {
		toStart = utf8.RuneError
		toEnd = utf8.RuneError
	} else {
		// If from pattern is reverted, only the last rune in the to pattern will be used.
		if reverted {
			var size int

			for len(to) > 0 {
				toStart, size = utf8.DecodeRuneInString(to)
				to = to[size:]
			}

			toEnd = utf8.RuneError
		} else {
			to, toStart, toEnd, toRangeStep = nextRuneRange(to, utf8.RuneError)
		}
	}

	fromEnd = utf8.RuneError

	for len(from) > 0 {
		from, fromStart, fromEnd, fromRangeStep = nextRuneRange(from, fromEnd)

		// fromStart is a single character. Just map it with a rune in the to pattern.
		if fromRangeStep == 0 {
			singleRunes = tr.addRune(fromStart, toStart, singleRunes)
			updateRange()
			continue
		}

		for toEnd != utf8.RuneError && fromStart != fromEnd {
			// If mapped rune is a single character instead of a range, simply shift first
			// rune in the range.
			if toRangeStep == 0 {
				singleRunes = tr.addRune(fromStart, toStart, singleRunes)
				updateRange()
				fromStart += fromRangeStep
				continue
			}

			fromRangeSize = (fromEnd - fromStart) * fromRangeStep
			toRangeSize = (toEnd - toStart) * toRangeStep

			// Not enough runes in the to pattern. Need to read more.
			if fromRangeSize > toRangeSize {
				fromStart, toStart = tr.addRuneRange(fromStart, fromStart+toRangeSize*fromRangeStep, toStart, toEnd, singleRunes)
				fromStart += fromRangeStep
				updateRange()

				// Edge case: If fromRangeSize == toRangeSize + 1, the last fromStart value needs be considered
				// as a single rune.
				if fromStart == fromEnd {
					singleRunes = tr.addRune(fromStart, toStart, singleRunes)
					updateRange()
				}

				continue
			}

			fromStart, toStart = tr.addRuneRange(fromStart, fromEnd, toStart, toStart+fromRangeSize*toRangeStep, singleRunes)
			updateRange()
			break
		}

		if fromStart == fromEnd {
			fromEnd = utf8.RuneError
			continue
		}

		fromStart, toStart = tr.addRuneRange(fromStart, fromEnd, toStart, toStart, singleRunes)
		fromEnd = utf8.RuneError
	}

	if fromEnd != utf8.RuneError {
		singleRunes = tr.addRune(fromEnd, toStart, singleRunes)
	}

	tr.reverted = reverted
	tr.mappedRune = -1
	tr.hasPattern = true

	// Translate RuneError only if in deletion or reverted mode.
	if deletion || reverted {
		tr.mappedRune = toStart
	}

	return tr
}

func (tr *Translator) addRune(from, to rune, singleRunes []rune) []rune {
	if from <= unicode.MaxASCII {
		if tr.quickDict == nil {
			tr.quickDict = &runeDict{}
		}

		tr.quickDict.Dict[from] = to
	} else {
		if tr.runeMap == nil {
			tr.runeMap = make(runeMap)
		}

		tr.runeMap[from] = to
	}

	singleRunes = append(singleRunes, from)
	return singleRunes
}

func (tr *Translator) addRuneRange(fromLo, fromHi, toLo, toHi rune, singleRunes []rune) (rune, rune) {
	var r rune
	var rrm *runeRangeMap

	if fromLo < fromHi {
		rrm = &runeRangeMap{
			FromLo: fromLo,
			FromHi: fromHi,
			ToLo:   toLo,
			ToHi:   toHi,
		}
	} else {
		rrm = &runeRangeMap{
			FromLo: fromHi,
			FromHi: fromLo,
			ToLo:   toHi,
			ToHi:   toLo,
		}
	}

	// If there is any single rune conflicts with this rune range, clear single rune record.
	for _, r = range singleRunes {
		if rrm.FromLo <= r && r <= rrm.FromHi {
			if r <= unicode.MaxASCII {
				tr.quickDict.Dict[r] = 0
			} else {
				delete(tr.runeMap, r)
			}
		}
	}

	tr.ranges = append(tr.ranges, rrm)
	return fromHi, toHi
}

func nextRuneRange(str string, last rune) (remaining string, start, end rune, rangeStep rune) {
	var r rune
	var size int

	remaining = str
	escaping := false
	isRange := false

	for len(remaining) > 0 {
		r, size = utf8.DecodeRuneInString(remaining)
		remaining = remaining[size:]

		// Parse special characters.
		if !escaping {
			if r == '\\' {
				escaping = true
				continue
			}

			if r == '-' {
				// Ignore slash at beginning of string.
				if last == utf8.RuneError {
					continue
				}

				start = last
				isRange = true
				continue
			}
		}

		escaping = false

		if last != utf8.RuneError {
			// This is a range which start and end are the same.
			// Considier it as a normal character.
			if isRange && last == r {
				isRange = false
				continue
			}

			start = last
			end = r

			if isRange {
				if start < end {
					rangeStep = 1
				} else {
					rangeStep = -1
				}
			}

			return
		}

		last = r
	}

	start = last
	end = utf8.RuneError
	return
}

// Translate str with a from/to pattern pair.
//
// See comment in Translate function for usage and samples.
func (tr *Translator) Translate(str string) string {
	if !tr.hasPattern || str == "" {
		return str
	}

	var r rune
	var size int
	var needTr bool

	orig := str

	var output *bytes.Buffer

	for len(str) > 0 {
		r, size = utf8.DecodeRuneInString(str)
		r, needTr = tr.TranslateRune(r)

		if needTr && output == nil {
			output = allocBuffer(orig, str)
		}

		if r != utf8.RuneError && output != nil {
			output.WriteRune(r)
		}

		str = str[size:]
	}

	// No character is translated.
	if output == nil {
		return orig
	}

	return output.String()
}

// TranslateRune return translated rune and true if r matches the from pattern.
// If r doesn't match the pattern, original r is returned and translated is false.
func (tr *Translator) TranslateRune(r rune) (result rune, translated bool) {
	switch {
	case tr.quickDict != nil:
		if r <= unicode.MaxASCII {
			result = tr.quickDict.Dict[r]

			if result != 0 {
				translated = true

				if tr.mappedRune >= 0 {
					result = tr.mappedRune
				}

				break
			}
		}

		fallthrough

	case tr.runeMap != nil:
		var ok bool

		if result, ok = tr.runeMap[r]; ok {
			translated = true

			if tr.mappedRune >= 0 {
				result = tr.mappedRune
			}

			break
		}

		fallthrough

	default:
		var rrm *runeRangeMap
		ranges := tr.ranges

		for i := len(ranges) - 1; i >= 0; i-- {
			rrm = ranges[i]

			if rrm.FromLo <= r && r <= rrm.FromHi {
				translated = true

				if tr.mappedRune >= 0 {
					result = tr.mappedRune
					break
				}

				if rrm.ToLo < rrm.ToHi {
					result = rrm.ToLo + r - rrm.FromLo
				} else if rrm.ToLo > rrm.ToHi {
					// ToHi can be smaller than ToLo if range is from higher to lower.
					result = rrm.ToLo - r + rrm.FromLo
				} else {
					result = rrm.ToLo
				}

				break
			}
		}
	}

	if tr.reverted {
		if !translated {
			result = tr.mappedRune
		}

		translated = !translated
	}

	if !translated {
		result = r
	}

	return
}

// HasPattern returns true if Translator has one pattern at least.
func (tr *Translator) HasPattern() bool {
	return tr.hasPattern
}

// Translate str with the characters defined in from replaced by characters defined in to.
//
// From and to are patterns representing a set of characters. Pattern is defined as following.
//
//     * Special characters
//       * '-' means a range of runes, e.g.
//         * "a-z" means all characters from 'a' to 'z' inclusive;
//         * "z-a" means all characters from 'z' to 'a' inclusive.
//       * '^' as first character means a set of all runes excepted listed, e.g.
//         * "^a-z" means all characters except 'a' to 'z' inclusive.
//       * '\' escapes special characters.
//     * Normal character represents itself, e.g. "abc" is a set including 'a', 'b' and 'c'.
//
// Translate will try to find a 1:1 mapping from from to to.
// If to is smaller than from, last rune in to will be used to map "out of range" characters in from.
//
// Note that '^' only works in the from pattern. It will be considered as a normal character in the to pattern.
//
// If the to pattern is an empty string, Translate works exactly the same as Delete.
//
// Samples:
//     Translate("hello", "aeiou", "12345")    => "h2ll4"
//     Translate("hello", "a-z", "A-Z")        => "HELLO"
//     Translate("hello", "z-a", "a-z")        => "svool"
//     Translate("hello", "aeiou", "*")        => "h*ll*"
//     Translate("hello", "^l", "*")           => "**ll*"
//     Translate("hello ^ world", `\^lo`, "*") => "he*** * w*r*d"
func Translate(str, from, to string) string {
	tr := NewTranslator(from, to)
	return tr.Translate(str)
}

// Delete runes in str matching the pattern.
// Pattern is defined in Translate function.
//
// Samples:
//     Delete("hello", "aeiou") => "hll"
//     Delete("hello", "a-k")   => "llo"
//     Delete("hello", "^a-k")  => "he"
func Delete(str, pattern string) string {
	tr := NewTranslator(pattern, "")
	return tr.Translate(str)
}

// Count how many runes in str match the pattern.
// Pattern is defined in Translate function.
//
// Samples:
//     Count("hello", "aeiou") => 3
//     Count("hello", "a-k")   => 3
//     Count("hello", "^a-k")  => 2
func Count(str, pattern string) int {
	if pattern == "" || str == "" {
		return 0
	}

	var r rune
	var size int
	var matched bool

	tr := NewTranslator(pattern, "")
	cnt := 0

	for len(str) > 0 {
		r, size = utf8.DecodeRuneInString(str)
		str = str[size:]

		if _, matched = tr.TranslateRune(r); matched {
			cnt++
		}
	}

	return cnt
}

// Squeeze deletes adjacent repeated runes in str.
// If pattern is not empty, only runes matching the pattern will be squeezed.
//
// Samples:
//     Squeeze("hello", "")             => "helo"
//     Squeeze("hello", "m-z")          => "hello"
//     Squeeze("hello   world", " ")    => "hello world"
func Squeeze(str, pattern string) string {
	var last, r rune
	var size int
	var skipSqueeze, matched bool
	var tr *Translator
	var output *bytes.Buffer

	orig := str
	last = -1

	if len(pattern) > 0 {
		tr = NewTranslator(pattern, "")
	}

	for len(str) > 0 {
		r, size = utf8.DecodeRuneInString(str)

		// Need to squeeze the str.
		if last == r && !skipSqueeze {
			if tr != nil {
				if _, matched = tr.TranslateRune(r); !matched {
					skipSqueeze = true
				}
			}

			if output == nil {
				output = allocBuffer(orig, str)
			}

			if skipSqueeze {
				output.WriteRune(r)
			}
		} else {
			if output != nil {
				output.WriteRune(r)
			}

			last = r
			skipSqueeze = false
		}

		str = str[size:]
	}

	if output == nil {
		return orig
	}

	return output.String()
}

//构建请求url
func BuildUrl(base string, params map[string]string) string {
	v := url.Values{}
	for k, p := range params {
		v.Add(k, p)
	}
	baseLen := len(base)
	if baseLen > 0 {
		if base[baseLen-1] == '/' {
			base = Substr(base, 0, baseLen-1)
		}
		base = base + "?"
	}
	return base + v.Encode()
}

//解析url
func GetUrlParam(s string) (url.Values, error) {
	//解析这个 URL 并确保解析没有出错。
	u, err := url.Parse(s)
	if err != nil {
		return nil, err
	}
	return url.ParseQuery(u.RawQuery)
}

func GetUrlHost(s string) (string, error) {
	//解析这个 URL 并确保解析没有出错。
	u, err := url.Parse(s)
	if err != nil {
		return "", err
	}
	h := strings.Split(u.Host, ":")
	return h[0], nil
}

func GetUrlBase(s string) (string, error) {
	//解析这个 URL 并确保解析没有出错。
	u, err := url.Parse(s)
	if err != nil {
		return "", err
	}
	h := strings.Split(u.Host, ":")
	return u.Scheme + "://" + h[0], nil
}

func GetUrlRootDomain(s string) (string, error) {
	//解析这个 URL 并确保解析没有出错。
	u, err := url.Parse(s)
	if err != nil {
		return "", err
	}
	h := strings.Split(u.Host, ":")
	host := strings.Split(h[0], ".")
	hostLen := len(host)
	if hostLen < 2 {
		return "", fmt.Errorf("url format error")
	}
	return host[hostLen-2] + "." + host[hostLen-1], nil
}

func IsUrl(str string) bool {
	return rxUrl.MatchString(str)
}

func IsHost(str string) bool {
	return rxHost.MatchString(str)
}

func GetUriLocation(str string) string {
	return ReplacePattern(str, `\?.*`, "")
}

//分号数据组即 1:2,3:4
func SpiltColonCommaInt(str string) (res []map[int]int, err error) {
	res = make([]map[int]int, 0)
	if str == "" {
		return res, nil
	}
	var re = regexp.MustCompile(`^(\d+:\d+,?)+$`)
	if !re.MatchString(str) {
		return nil, errors.New("spilt format error")
	}
	pairs := strings.Split(str, ",")

	for _, v := range pairs {
		if len(v) == 0 {
			continue
		}
		pair := strings.Split(v, ":")
		var i1, i2 int
		i1, _ = strconv.Atoi(pair[0])
		i2, _ = strconv.Atoi(pair[1])
		item := map[int]int{i1: i2}
		res = append(res, item)
	}
	return res, nil
}

//逗号数据组 即 1,2,3
func SpiltCommaInt(str string) ([]int, error) {
	res := make([]int, 0)
	if str == "" {
		return res, nil
	}
	var re = regexp.MustCompile(`^(\d+,?)+$`)
	if !re.MatchString(str) {
		return nil, errors.New("spilt format error")
	}
	ids := strings.Split(str, ",")
	for _, v := range ids {
		if len(v) == 0 {
			continue
		}
		id, _ := strconv.Atoi(v)
		res = append(res, id)
	}
	return res, nil
}

//字符串逗号数组即 a:1,b:2
func SpiltCommaMap(str string) (res map[string]int, err error) {
	var re = regexp.MustCompile(`^([\w_]+:[\w_]+,?)+$`)
	if !re.MatchString(str) {
		return nil, errors.New("spilt format error")
	}
	pairs := strings.Split(str, ",")
	res = make(map[string]int, 0)
	for _, v := range pairs {
		if len(v) == 0 {
			continue
		}
		pair := strings.Split(v, ":")
		var i1 string
		var i2 int
		i1 = pair[0]
		i2, _ = strconv.Atoi(pair[1])
		res[i1] = i2
	}
	return res, nil
}

func UrlEncoded(str string) (string, error) {
	u, err := url.Parse(str)
	if err != nil {
		return "", err
	}
	return u.String(), nil
}

func RandomBytes(n int) []byte {
	b := make([]byte, n)
	_, err := crrand.Read(b)
	if err != nil {
		panic(err)
	}
	return b
}

func RandomBase64(s int) string {
	return RandString(s, "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ+/")
}

//实际就是ASCII码的十六进制 注意长度是s的两倍
func RandomHex(s int) string {
	b := RandomBytes(s)
	hexstring := hex.EncodeToString(b)
	return strings.ToUpper(hexstring)
}

func RandString(s int, letters ...string) string {
	randomFactor := RandomBytes(1)
	rand.Seed(time.Now().UnixNano() * int64(randomFactor[0]))
	var letterRunes = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	if len(letters) > 0 {
		letterRunes = []rune(letters[0])
	}
	b := make([]rune, s)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func Explode(src []int, sep string) string {
	var res = ""
	for _, v := range src {
		if res == "" {
			res = strconv.Itoa(v)
		} else {
			res += sep + strconv.Itoa(v)
		}
	}
	return res
}
