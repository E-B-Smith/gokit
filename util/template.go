/**
@file          template.go
@package       util
@brief         Web template helpers.
@author        Edward Smith
@date          July 2019
@copyright     -©- Copyright © 2019 Edward Smith. All rights reserved. -©-
*/

package util

import (
	"fmt"
	"html"
	"html/template"
	"math"
	"strings"
	"time"

	"github.com/E-B-Smith/gokit/log"
)

// UnescapeHTMLString un-escapes an HTML string in a template.
func UnescapeHTMLString(args ...interface{}) string {
	ok := false
	var s string
	if len(args) == 1 {
		s, ok = args[0].(string)
		s = html.UnescapeString(s)
	}
	if !ok {
		s = fmt.Sprint(args...)
	}
	return s
}

// EscapeHTMLString escape an HTML string in a template.
func EscapeHTMLString(args ...interface{}) string {
	ok := false
	var s string
	if len(args) == 1 {
		s, ok = args[0].(string)
		s = html.EscapeString(s)
	}
	if !ok {
		s = fmt.Sprint(args...)
	}
	return s
}

// BoolPtr safely returns a bool from a bool pointer in a template.
func BoolPtr(b *bool) bool {
	if b != nil && *b {
		return true
	}
	return false
}

// StringPtr safely returns a string from a string pointer.
func StringPtr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

const monthYearFormat string = "1/2006"

// TimeFromDouble returns time.Time from an epoch double.
func TimeFromDouble(timestamp float64) time.Time {
	i, f := math.Modf(timestamp)
	var sec = int64(math.Floor(i))
	var nsec = int64(f * 1000000.0)
	return time.Unix(sec, nsec)
}

// MonthYearStringFromEpochPtr format a date ptr.
func MonthYearStringFromEpochPtr(epoch *float64) string {
	if epoch == nil || *epoch <= 0.0 {
		return ""
	}
	t := TimeFromDouble(*epoch)
	return t.Format(monthYearFormat)
}

// ParseMonthYearString parses a date format.
func ParseMonthYearString(s string) time.Time {
	s = StringIncludingCharactersInSet(s, "0123456789/")
	t, _ := time.Parse(monthYearFormat, s)
	return t
}

// LoadTemplates loads the templates in a given path.
func LoadTemplates(templatesPath string) (*template.Template, error) {
	var error error
	templatesPath = strings.TrimSpace(templatesPath)
	if len(templatesPath) == 0 {
		templatesPath = ResourcePath("templates")
	}
	log.Infof("Loading templates from '%s'.", templatesPath)

	path := templatesPath + "/*"
	templates := template.New("Base")
	templates = templates.Funcs(template.FuncMap{
		"UnescapeHTMLString": UnescapeHTMLString,
		"EscapeHTMLString":   EscapeHTMLString,
		"BoolPtr":            BoolPtr,
		"StringPtr":          StringPtr,
		"MonthYearString":    MonthYearStringFromEpochPtr,
		"HumanDurationBrief": HumanDurationBrief,
	})
	templates, error = templates.ParseGlob(path)
	if error != nil || templates == nil {
		if error == nil {
			error = fmt.Errorf("no files")
		}
		error = fmt.Errorf("can't parse template files: %v", error)
		return nil, error
	}
	return templates, nil
}
