/**
@file          sqlstring.go
@package       scanner
@brief         Scans a formatted SQL string.
@author        Edward Smith
@date          November 2014
@copyright     -©- Copyright © 2014-2016 Edward Smith, all rights reserved. -©-
*/

package scanner

import (
	"bytes"
	"errors"
)

// ScanSQLString scans an SQL string.
func (scanner *Scanner) ScanSQLString() (string, error) {
	scanner.error = scanner.ScanSpaces()

	var r rune
	r, _, scanner.error = scanner.reader.ReadRune()
	if scanner.error != nil {
		return "", scanner.error
	}

	switch r {

	case '"':
		scanner.reader.UnreadRune()
		return scanner.ScanQuotedString()

	case '\'':
		break

	default:
		scanner.reader.UnreadRune()
		return scanner.ScanString()
	}

	var buffer bytes.Buffer
	for {
		r, _, scanner.error = scanner.reader.ReadRune()
		if scanner.error != nil {
			scanner.token = ""
			scanner.error = errors.New("quote error")
			return scanner.token, scanner.error
		}

		switch r {
		case '\'':
			r, _, scanner.error = scanner.reader.ReadRune()
			if scanner.error != nil {
				scanner.token = buffer.String()
				scanner.error = nil
				return scanner.token, scanner.error
			}
			if r == '\'' {
				buffer.WriteRune(r)
			} else {
				scanner.token = buffer.String()
				return scanner.token, nil
			}
		default:
			buffer.WriteRune(r)
		}
	}
}
