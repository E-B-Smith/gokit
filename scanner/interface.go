/**
@file          interface.go
@package       scanner
@brief         Parses an interface from a file by reflecting the variable names from an interface structure.
@author        Edward Smith
@date          March 2016
@copyright     -©- Copyright © 2014-2016 Edward Smith, all rights reserved. -©-
*/

package scanner

import (
	"fmt"
	"io"
	"reflect"
	"strings"
	"unicode"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// ScanInterface scans the input into a the passed interface.
func (scanner *Scanner) ScanInterface(config interface{}) error {
	//  Scan the input, finding fields by reflection  --

	configPtr := reflect.ValueOf(config)
	if configPtr.Kind() != reflect.Ptr {
		panic(fmt.Errorf("Pointer to struct expected"))
	}
	configPtrValue := reflect.ValueOf(config).Elem()
	if configPtrValue.Kind() != reflect.Struct {
		panic(fmt.Errorf("Pointer to struct expected"))
	}

	for !scanner.IsAtEnd() {
		var error error

		var identifier string
		identifier, error = scanner.ScanIdentifier()
		// log.Debugf("Scanned '%s'.", scanner.Token())

		if error == io.EOF {
			break
		}
		if error != nil {
			return error
		}

		//  Find the identifier --

		fieldName := CamelCaseFromIdentifier(identifier)
		//Log.Debugf("FieldName: '%s'.", fieldName)
		field := configPtrValue.FieldByName(fieldName)
		if !field.IsValid() {
			return scanner.SetErrorMessage("Configuration identifier expected")
		}
		structField, _ := configPtrValue.Type().FieldByName(fieldName)

		var (
			i int64
			s string
			b bool
			f float64
		)

		switch field.Type().Kind() {

		case reflect.Bool:
			b, error = scanner.ScanBool()
			if error != nil {
				return error
			}
			field.SetBool(b)

		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:

			enumValues := structField.Tag.Get("enum")
			if len(enumValues) > 0 {

				s, error = scanner.ScanNext()
				if error != nil {
					scanner.SetError(error)
					return scanner.LastError()
				}

				i, error = enumFromString(s, enumValues)
				if error != nil {
					return error
				}

			} else {

				i, error = scanner.ScanInt64()
				if error != nil {
					return error
				}

			}
			field.SetInt(i)

		case reflect.Float32, reflect.Float64:
			f, error = scanner.ScanFloat64()
			if error != nil {
				return error
			}
			field.SetFloat(f)

		case reflect.String:
			s, error = scanner.ScanNext()
			if error != nil {
				return error
			}
			field.SetString(s)

		case reflect.Struct:
			if !field.CanSet() {
				scanner.SetError(fmt.Errorf("struct '%s' has unset-able fields", fieldName))
				return scanner.error
			}
			s, error = scanner.ScanString()
			if error != nil || s != "{" {
				scanner.error = fmt.Errorf("expected '{'")
				return scanner.error
			}
			// TODO: This copies to a temp buffer then copies to the struct. Copy directly?
			p := reflect.New(field.Type())
			error = scanner.ScanInterface(p.Interface())
			if error == nil {
				scanner.error = fmt.Errorf("expected '}'")
				return scanner.error
			}
			if !strings.Contains(error.Error(), "Scanned '}'") {
				return scanner.error
			}
			scanner.SetError(nil)
			field.Set(p.Elem())

		default:
			return fmt.Errorf("Error: '%s' unhandled type: %s", identifier, field.Type().Name())
		}
	}
	return nil
}

func enumFromString(s string, enumValues string) (int64, error) {

	enumArray := make([]string, 0)
	a := strings.Split(enumValues, ",")
	for _, enum := range a {
		enum = strings.TrimSpace(enum)
		if len(enum) > 0 {
			enumArray = append(enumArray, enum)
		}
	}

	for i, val := range enumArray {
		if val == s {
			return int64(i), nil
		}
	}

	return -1, fmt.Errorf("Invalid enum '%s'", s)
}

func titleCased(s string) string {
	return cases.Title(language.Und, cases.NoLower).String(s)
}

// CamelCaseFromIdentifier transforms a snake case identifier into a camel case identifier.
func CamelCaseFromIdentifier(s string) string {

	lastWasUpper := false
	words := make([]string, 0, 5)
	var word []rune
	for _, r := range s {

		switch {

		case r == '-' || r == '_':
			words = append(words, string(word))
			word = word[:0]
			lastWasUpper = false

		case unicode.IsUpper(r):
			if !lastWasUpper {
				words = append(words, string(word))
				word = word[:0]
			}
			word = append(word, r)
			lastWasUpper = true

		default:
			word = append(word, r)
			lastWasUpper = false
		}
	}
	if len(word) > 0 {
		words = append(words, string(word))
	}

	//  String together the parts.  Upper-case any special words:

	upperWords := map[string]bool{
		"http":  true,
		"https": true,
		"url":   true,
		"uri":   true,
		"urn":   true,
		"smtp":  true,
		"xml":   true,
		"json":  true,
		"id":    true,
	}

	var camelString string
	for _, part := range words {
		part = strings.ToLower(part)
		if _, ok := upperWords[part]; ok {
			part = strings.ToUpper(part)
		} else {
			part = titleCased(part)
		}
		camelString += part
	}

	return camelString
}
