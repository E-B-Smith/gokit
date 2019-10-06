/**
@file          scanner.go
@package       scanner
@brief         Parsing and scanning functions.
@author        Edward Smith
@date          November 2014
@copyright     -©- Copyright © 2014-2016 Edward Smith, all rights reserved. -©-
*/

/*
Package scanner provides scanning and parsing functions for converting strings and file contents
into data.
*/
package scanner

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"strconv"
	"strings"
	"unicode"
)

// Scanner is used to scan a
type Scanner struct {
	filename   string
	reader     *bufio.Reader
	lineNumber int
	error      error
	token      string
}

// NewScannerWithFilename creates and returns a Scanner that reads from a file named `filename`.
func NewScannerWithFilename(filename string) *Scanner {
	inputFile, error := os.Open(filename)
	s := NewScannerWithFile(inputFile)
	if error != nil {
		s.error = fmt.Errorf("can't open file '%s' for reading: %v", filename, error)
	}
	return s
}

// NewScannerWithFile creates and returns a Scanner that reads from a file.
func NewScannerWithFile(file *os.File) *Scanner {
	if file == nil {
		return nil
	}
	scanner := new(Scanner)
	scanner.filename = file.Name()
	scanner.reader = bufio.NewReader(file)
	scanner.lineNumber = 1
	scanner.token = ""
	return scanner
}

// NewScannerWithString creates and returns a Scanner that reads from a string.
func NewScannerWithString(s string) *Scanner {
	return NewScannerWithReader(strings.NewReader(s))
}

// NewScannerWithReader returns a scanner with input from the io.Reader.
func NewScannerWithReader(r io.Reader) *Scanner {
	scanner := new(Scanner)
	scanner.reader = bufio.NewReader(r)
	scanner.lineNumber = 1
	scanner.token = ""
	return scanner
}

// FileName returns the name of the current file being scanned.
func (scanner *Scanner) FileName() string {
	return scanner.filename
}

// LineNumber returns the current line number being scanned.
func (scanner *Scanner) LineNumber() int {
	return scanner.lineNumber
}

// IsAtEnd returns true if the scanner has read all the data.
func (scanner *Scanner) IsAtEnd() bool {
	return scanner.error != nil
}

// Token returns the current scanned token.
func (scanner *Scanner) Token() string {
	return scanner.token
}

// LastError return the last error encountered.
func (scanner *Scanner) LastError() error {
	return scanner.error
}

// SetErrorMessage sets the current error message.
func (scanner *Scanner) SetErrorMessage(message string) error {
	basename := path.Base(scanner.FileName())
	message = fmt.Sprintf("%s:%d Scanned '%s'. %s",
		basename, scanner.LineNumber(), scanner.Token(), message)
	scanner.error = errors.New(message)
	return scanner.error
}

// SetError sets the current error.
func (scanner *Scanner) SetError(error error) error {
	if error == nil {
		scanner.error = nil
		return scanner.error
	}
	basename := path.Base(scanner.FileName())
	message := fmt.Sprintf("%s:%d Scanned '%s'. %v",
		basename, scanner.LineNumber(), scanner.Token(), error)
	scanner.error = errors.New(message)
	return scanner.error
}

// NextRune returns the next rune to be scanned.
func (scanner *Scanner) NextRune() rune {
	var r rune
	r, _, scanner.error = scanner.reader.ReadRune()
	scanner.reader.UnreadRune()
	return r
}

// NextRuneIsDigit returns true if the next rune of the token is a digit.
func (scanner *Scanner) NextRuneIsDigit() bool {
	return unicode.IsDigit(scanner.NextRune())
}

// NextRuneIsPunct returns true if the next rune of the token is punctuation.
func (scanner *Scanner) NextRuneIsPunct() bool {
	return unicode.IsPunct(scanner.NextRune())
}

//  Scan Routines --

// IsValidIdentifierStartRune returns true if the rune is a valid run to start an identifier.
func IsValidIdentifierStartRune(r rune) bool {
	return unicode.IsLetter(r) || r == '_'
}

// IsValidIdentifierRune returns true if the rune is a valid identifier rune.
func IsValidIdentifierRune(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r) || r == '-' || r == '_'
}

// IsOctalDigit returns true if the rune is a valid octal digit.
func IsOctalDigit(r rune) bool {
	return unicode.IsDigit(r) && r != '8' && r != '9'
}

// ZIsSpace returns true if the rune is a valid space-like character.
func ZIsSpace(r rune) bool {
	return unicode.IsSpace(r) || r == '#'
}

// ZIsLineFeed returns true if the rune is a linefeed-like character.
func ZIsLineFeed(r rune) bool {
	return r == '\n' || r == '\u0085'
}

// ScanSpaces scans forward through whitespace characters.
func (scanner *Scanner) ScanSpaces() error {
	scanner.token = ""
	for !scanner.IsAtEnd() {
		var r rune
		r, _, scanner.error = scanner.reader.ReadRune()

		if r == '#' {
			for !scanner.IsAtEnd() && !ZIsLineFeed(r) {
				r, _, scanner.error = scanner.reader.ReadRune()
			}
		}
		if ZIsLineFeed(r) {
			scanner.lineNumber++
			continue
		}
		if ZIsSpace(r) {
			continue
		}

		scanner.reader.UnreadRune()
		return nil
	}

	return scanner.error
}

// IsValidStringRune returns true if the character is a valid string character.
func IsValidStringRune(r rune) bool {
	if r == ';' || r == ',' || ZIsSpace(r) {
		return false
	}
	return unicode.IsGraphic(r)
}

// ScanString returns a scanned string.
func (scanner *Scanner) ScanString() (next string, error error) {
	error = scanner.ScanSpaces()

	var (
		r      rune
		buffer bytes.Buffer
	)
	r, _, scanner.error = scanner.reader.ReadRune()

	for IsValidStringRune(r) {
		buffer.WriteRune(r)
		r, _, scanner.error = scanner.reader.ReadRune()
	}
	scanner.reader.UnreadRune()

	scanner.token = buffer.String()
	return scanner.token, nil
}

// ScanInt64 scans an int64 integer.
func (scanner *Scanner) ScanInt64() (int int64, error error) {
	scanner.ScanSpaces()
	var r rune
	r, _, scanner.error = scanner.reader.ReadRune()

	if !unicode.IsDigit(r) {
		scanner.reader.UnreadRune()
		scanner.token, _ = scanner.ScanNext()
		return 0, scanner.SetErrorMessage("Integer expected")
	}

	var buffer bytes.Buffer
	for unicode.IsDigit(r) {
		buffer.WriteRune(r)
		r, _, scanner.error = scanner.reader.ReadRune()
	}
	scanner.reader.UnreadRune()

	scanner.token = buffer.String()
	var i int64
	i, scanner.error = strconv.ParseInt(scanner.token, 10, 64)
	return i, scanner.error
}

// ScanInt32 scans an int32 integer.
func (scanner *Scanner) ScanInt32() (int int32, error error) {
	i, error := scanner.ScanInt64()
	return int32(i), error
}

// ScanInt scans an int integer.
func (scanner *Scanner) ScanInt() (int, error) {
	i64, error := scanner.ScanInt64()
	return int(i64), error
}

// ScanFloat64 scans a float64.
func (scanner *Scanner) ScanFloat64() (float64, error) {
	scanner.ScanSpaces()

	var r rune
	var buffer bytes.Buffer
	r, _, scanner.error = scanner.reader.ReadRune()
	for scanner.error == nil && (unicode.IsDigit(r) || r == '-' || r == '.') {
		buffer.WriteRune(r)
		r, _, scanner.error = scanner.reader.ReadRune()
	}
	scanner.reader.UnreadRune()

	scanner.token = buffer.String()
	var f float64
	f, scanner.error = strconv.ParseFloat(scanner.token, 64)
	return f, scanner.error
}

// ScanBool scans a bool.
func (scanner *Scanner) ScanBool() (value bool, error error) {
	var s string
	s, scanner.error = scanner.ScanNext()
	if scanner.error != nil {
		return false, scanner.error
	}

	s = strings.ToLower(s)
	if s == "true" || s == "yes" || s == "t" || s == "y" || s == "1" {
		return true, nil
	}
	if s == "false" || s == "no" || s == "f" || s == "n" || s == "0" {
		return false, nil
	}
	scanner.SetErrorMessage("Expected a boolean value")
	return false, scanner.error
}

// ScanIdentifier scans an identifier.
func (scanner *Scanner) ScanIdentifier() (identifier string, error error) {
	scanner.ScanSpaces()
	var r rune
	r, _, scanner.error = scanner.reader.ReadRune()
	if scanner.error != nil {
		return "", scanner.error
	}

	if !IsValidIdentifierStartRune(r) {
		scanner.reader.UnreadRune()
		scanner.ScanNext()
		return "", scanner.SetErrorMessage("Expected an identifier")
	}

	var buffer bytes.Buffer
	for IsValidIdentifierRune(r) {
		buffer.WriteRune(r)
		r, _, scanner.error = scanner.reader.ReadRune()
	}
	scanner.reader.UnreadRune()

	scanner.token = buffer.String()
	return scanner.token, nil
}

// ScanOctal scans an octal number.
func (scanner *Scanner) ScanOctal() (Integer int, error error) {
	scanner.ScanSpaces()
	var r rune
	r, _, scanner.error = scanner.reader.ReadRune()

	if !IsOctalDigit(r) {
		scanner.reader.UnreadRune()
		scanner.token, _ = scanner.ScanNext()
		return 0, scanner.SetErrorMessage("Octal number expected")
	}

	var buffer bytes.Buffer
	for IsOctalDigit(r) {
		buffer.WriteRune(r)
		r, _, scanner.error = scanner.reader.ReadRune()
	}
	scanner.reader.UnreadRune()

	scanner.token = buffer.String()
	val, error := strconv.ParseInt(scanner.token, 8, 0)
	return int(val), error
}

// ScanQuotedString scans a quoted string.
func (scanner *Scanner) ScanQuotedString() (string, error) {
	scanner.ScanSpaces()
	var r rune
	r, _, scanner.error = scanner.reader.ReadRune()

	if r != '"' {
		scanner.token = string(r)
		return "", scanner.SetErrorMessage("Quoted string expected")
	}
	scanner.reader.UnreadRune()

	parseCount, error := fmt.Fscanf(scanner.reader, "%q", &scanner.token)
	if error != nil {
		return "", scanner.SetError(error)
	}
	if parseCount != 1 {
		return "", scanner.SetErrorMessage("Quoted string expected")
	}

	return scanner.token, nil
}

// ScanNext scans the next item.
func (scanner *Scanner) ScanNext() (next string, error error) {
	scanner.ScanSpaces()
	var r rune
	r, _, scanner.error = scanner.reader.ReadRune()

	if r == '"' {
		scanner.reader.UnreadRune()
		return scanner.ScanQuotedString()
	}
	if unicode.IsPunct(r) {
		var buffer bytes.Buffer
		buffer.WriteRune(r)
		scanner.token = buffer.String()
		return scanner.token, nil
	}
	if unicode.IsDigit(r) {
		scanner.reader.UnreadRune()
		scanner.ScanInt64()
		return scanner.token, scanner.error
	}
	scanner.reader.UnreadRune()
	return scanner.ScanString()
}

// ScanToEOL scans to the end of the current line.
func (scanner *Scanner) ScanToEOL() (string, error) {
	var (
		r      rune
		buffer bytes.Buffer
	)
	r, _, scanner.error = scanner.reader.ReadRune()

	for !scanner.IsAtEnd() && !ZIsLineFeed(r) {
		buffer.WriteRune(r)
		r, _, scanner.error = scanner.reader.ReadRune()
	}
	scanner.reader.UnreadRune()

	scanner.token = buffer.String()
	scanner.token = strings.TrimSpace(scanner.token)
	return scanner.token, nil
}
