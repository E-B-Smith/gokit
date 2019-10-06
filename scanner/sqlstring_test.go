/**
@file          sqlstring_test.go
@package       scanner
@brief         Tests for sqlstring.go.
@author        Edward Smith
@date          November 2014
@copyright     -©- Copyright © 2014-2016 Edward Smith, all rights reserved. -©-
*/

package scanner

import (
	"testing"
)

type TestCaseType struct {
	Test    string
	Result  string
	Success bool
}

func TestScanSQLString(t *testing.T) {

	TestCases := []TestCaseType{
		{"Unquoted  ", "Unquoted", true},
		{"'Quoted'  ", "Quoted", true},
		{"'Don''t look!'  ", "Don't look!", true},
		{"", "", false},
		{"''", "", true},
		{"'fail", "", false},
	}

	for _, tc := range TestCases {
		scanner := NewScannerWithString(tc.Test)
		result, error := scanner.ScanSQLString()
		if error != nil {
			if tc.Success {
				t.Errorf("Tested '%s' expected '%s' got '%s' error: %v.", tc.Test, tc.Result, result, error)
			}
		}
		if result != tc.Result {
			t.Errorf("Tested '%s' expected '%s' got '%s' error: %v.", tc.Test, tc.Result, result, error)
		}
	}
}
