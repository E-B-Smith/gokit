/**
@file          strings_test.go
@package       util
@brief         String utility functions.
@author        Edward Smith
@date          November, 2014
@copyright     -©- Copyright © 2014-2019 Edward Smith. All rights reserved. -©-
*/

package util

import (
	"bytes"
	"errors"
	"fmt"
	"math"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

// CleanStringPtr trims spaces and return a nil string if the string is empty.
func CleanStringPtr(s *string) *string {
	if s == nil {
		return s
	}
	temp := strings.TrimSpace(*s)
	if len(temp) == 0 {
		return nil
	}
	return &temp
}

// CleanStringPtrFromString converts a string to a clean string pointer.
func CleanStringPtrFromString(s string) *string {
	s = strings.TrimSpace(s)
	if len(s) == 0 {
		return nil
	}
	return &s
}

// StringIncludingCharactersInSet returns a string from `string` containing only characters in the `characterset`.
func StringIncludingCharactersInSet(inputstring string, characterset string) string {
	return strings.Map(
		func(r rune) rune {
			if !strings.ContainsRune(characterset, r) {
				return -1
			}
			return r
		},
		inputstring)
}

// StringExcludingCharactersInSet returns a string from `string` excluding characters in the `characterset`.
func StringExcludingCharactersInSet(inputstring string, characterset string) string {
	return strings.Map(
		func(r rune) rune {
			if !strings.ContainsRune(characterset, r) {
				return r
			}
			return -1
		},
		inputstring)
}

// ReplaceCharactersNotInSetWithRune returns a string with the characters in `characterset` replaced with `rune`.
func ReplaceCharactersNotInSetWithRune(inputstring string, characterset string, replacement rune) string {
	return strings.Map(
		func(r rune) rune {
			if !strings.ContainsRune(characterset, r) {
				return replacement
			}
			return r
		},
		inputstring)
}

// ValidatedEmailAddress does simple validation on an email address.
func ValidatedEmailAddress(inputstring string) (string, error) {
	inputstring = strings.TrimSpace(inputstring)
	re := regexp.MustCompile(`.+@.+\..+`)
	matched := re.Match([]byte(inputstring))
	if !matched {
		return "", errors.New("Invalid email address")
	}
	return inputstring, nil
}

// ValidatedPhoneNumber validates a phone number string.
func ValidatedPhoneNumber(inputstring string) (string, error) {
	inputstring = StringIncludingCharactersInSet(inputstring, "0123456789")
	if utf8.RuneCountInString(inputstring) > 10 {
		inputstring = strings.TrimLeft(inputstring, "1")
	}
	if utf8.RuneCountInString(inputstring) != 10 {
		return "", errors.New("Invalid phone number")
	}
	return inputstring, nil
}

// HumanBytes returns string that represents a number of bytes in a human readable form.
func HumanBytes(intBytes int64) string {
	suffix := []string{
		"B",
		"KB",
		"MB",
		"GB",
		"TB",
		"PB",
		"Really?",
	}
	if intBytes < 1024 {
		return fmt.Sprintf("%d B", intBytes)
	}

	idx := 0
	bytes := float64(intBytes)
	for bytes > 1024.0 && idx < len(suffix) {
		idx++
		bytes /= 1024.0
	}
	return fmt.Sprintf("%1.2f %s", bytes, suffix[idx])
}

// HumanInt returns string that represents an integer in a human readable form.
func HumanInt(i int64) string {
	var bo bytes.Buffer

	if i < 0 {
		bo.WriteByte('-')
		i *= -1
	}

	str := strconv.FormatInt(i, 10)
	bi := bytes.NewBufferString(str)
	c := utf8.RuneCountInString(str)

	for c > 0 {
		r, _, _ := bi.ReadRune()
		bo.WriteRune(r)
		if c%3 == 1 && c != 1 {
			bo.WriteString(",")
		}
		c--
	}
	return bo.String()
}

// HumanDuration returns a string representation of a duration in a human readable form.
func HumanDuration(d time.Duration) string {

	var (
		s           string
		dys, hrs, m int
		sec         float64
	)

	dys = int(math.Floor(d.Hours() / 24.0))
	hrs = int(math.Floor(math.Mod(d.Hours(), 24.0)))
	m = int(math.Floor(math.Mod(d.Minutes(), 60.0)))
	sec = math.Mod(d.Seconds(), 60.0)

	switch {
	case dys == 1:
		s = fmt.Sprintf("%d day %d:%02d:%04.1f hours", dys, hrs, m, sec)

	case dys > 0:
		s = fmt.Sprintf("%d days %d:%02d:%04.1f hours", dys, hrs, m, sec)

	case hrs > 0:
		s = fmt.Sprintf("%d:%02d:%04.1f hours", hrs, m, sec)

	case m > 0:
		s = fmt.Sprintf("%d:%04.1f minutes", m, sec)

	default:
		s = fmt.Sprintf("%1.3f seconds", sec)
	}

	return s
}

// HumanDurationBrief returns a string representation of a duration in a compact human readable form.
func HumanDurationBrief(d time.Duration) string {

	var (
		s           string
		dys, hrs, m int
		sec         float64
	)

	dys = int(math.Floor(d.Hours() / 24.0))
	hrs = int(math.Floor(math.Mod(d.Hours(), 24.0)))
	m = int(math.Floor(math.Mod(d.Minutes(), 60.0)))
	sec = math.Mod(d.Seconds(), 60.0)

	switch {
	case dys == 1:
		s = fmt.Sprintf("%d day %d:%02d", dys, hrs, m)

	case dys > 0:
		s = fmt.Sprintf("%d days %d:%02d", dys, hrs, m)

	case hrs > 0:
		s = fmt.Sprintf("%d:%02d:%02.0f", hrs, m, sec)

	case m > 0:
		s = fmt.Sprintf("%d:%02.0f", m, sec)

	default:
		s = fmt.Sprintf("0:%02.0f", sec)
	}

	return s
}

// CompareVersionStrings compares two semantic version strings and returns -1, 0, 1 indicating the result.
func CompareVersionStrings(version1 string, version2 string) int {
	v1 := strings.Split(version1, ".")
	v2 := strings.Split(version2, ".")

	for len(v1) != len(v2) {
		if len(v1) < len(v2) {
			v1 = append(v1, "0")
		} else {
			v2 = append(v2, "0")
		}
	}

	i := 0
	for i < len(v1) {
		i1, _ := strconv.Atoi(v1[i])
		i2, _ := strconv.Atoi(v2[i])
		if i1 < i2 {
			return -1
		}
		if i1 > i2 {
			return 1
		}
		i++
	}
	return 0
}

/*
func NewUUIDString() string {
	return uuid.NewV4().String()
}
*/

// IPAddressFromHTTPRequest returns a cleaned up IP address from an httpRequest header.
func IPAddressFromHTTPRequest(httpRequest *http.Request) string {
	if httpRequest == nil {
		return ""
	}
	address := httpRequest.Header.Get("x-forwarded-for")
	if address != "" {
		addressArray := strings.Split(address, ",")
		if len(addressArray) > 0 {
			address = addressArray[0]
			address = strings.TrimSpace(address)
		}
	}
	if address == "" {
		address = httpRequest.RemoteAddr
		i := strings.IndexRune(address, ']')
		if i < 0 {
			i = strings.IndexRune(address, ':')
			if i > 0 {
				address = address[:i-1]
			} else if i == 0 {
				address = ""
			}
		} else {
			address = address[:i]
		}
	}
	return address
}

// FirstNRunes returns a string up to `n` runes long.
func FirstNRunes(s string, n int) string {
	result := ""
	for _, runeValue := range s {
		if n <= 0 {
			return result
		}
		result += string(runeValue)
		n--
	}
	return result
}

// TruncateStringToLength returns a string up to `maxLength`. The string is terminated at the last full word and '...' is appended if necessary.
func TruncateStringToLength(s string, maxLength int) string {
	if len(s) < maxLength {
		return s
	}

	result := ""
	a := strings.Split(s, " ")
	for _, next := range a {
		if len(result)+len(next)+3 <= maxLength {
			result += next + " "
		} else {
			result = strings.TrimRight(result, " ")
			if len(result)+3 <= maxLength {
				result += "..."
			}
			return result
		}
	}
	return result
}
