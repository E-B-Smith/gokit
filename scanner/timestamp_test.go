/**
@file          timestamp_test.go
@package       scanner
@brief         Tests for the timestamp parser.
@author        Edward Smith
@date          November 2014
@copyright     -©- Copyright © 2014-2016 Edward Smith, all rights reserved. -©-
*/

package scanner

import (
	"testing"
	"time"
)

/*
Timezone support is kind of broken in GO so I'll drop those tests for now.

https://stackoverflow.com/questions/49084316/why-doesnt-gos-time-parse-parse-the-timezone-identifier
*/

// const testScanTimestampString = `
// 0 2006-01-02T15:04:00-08:00
// 1 Jan 2 2006, 15:04
// 2 Jan 2 2006, 3:04PM
// 3 Monday, 02-Jan-06 16:04:00 EST
// 4 02 Jan 06 16:04 MST
// 5 02 Jan 06 15:04 -0800
// 6 Mon Jan 2 15:04:00 2006
// 7 Mon Jan 2 16:04:00 MST 2006
// 8 Mon Jan 02 15:04:00 -0800 2006
// 9 Mon, 02 Jan 2006 16:04:00 MST
// 10 Mon, 02 Jan 2006 15:04:00 -0800

// 11             this is an error string
// `

const testScanTimestampString = `
0 2006-01-02T15:04:00-08:00
1 Jan 2 2006, 15:04
2 Jan 2 2006, 3:04PM
5 02 Jan 06 15:04 -0800
6 Mon Jan 2 15:04:00 2006
8 Mon Jan 02 15:04:00 -0800 2006
10 Mon, 02 Jan 2006 15:04:00 -0800

11             this is an error string
`

const testCount = 7

/*
func TestParseTimeZone(t *testing.T) {
	tm, error := time.ParseInLocation("2006-01-02T15:04:05Z07:00", "2006-01-02T15:04:00-08:00", nil)
	if error != nil {
		t.Errorf("Error %v", error)
	}
	testTime, error := time.ParseInLocation("3:04:05 01/02/2006 MST", "3:04:05 01/02/2006 MST", nil)
	if error != nil {
		t.Errorf("Error %v", error)
	}
	if !tm.Equal(testTime) {
		t.Errorf("\n  tm: %v\ntime: %v", tm, testTime)
	}
}
*/

func TestScanTimestamp(t *testing.T) {

	var error error
	var ts time.Time
	// fmt.Printf("%s", testScanTimestampString)
	testTime, error := time.Parse("2006-01-02T15:04:05Z07:00", "2006-01-02T15:04:00-08:00")
	if error != nil {
		panic(error)
	}
	scanner := NewScannerWithString(testScanTimestampString)

	for i := 0; i < testCount; i++ {
		scanner.ScanInt()
		ts, error = scanner.ScanTimestamp()
		if error != nil {
			t.Errorf("Test %d: Error %v.", i, error)
		} else if !ts.Equal(testTime) {
			t.Errorf(
				"Test %d:\nScanned: %s\n Wanted: %s\n  Input: %s\n   Diff: %v",
				i,
				ts.Format(time.RFC3339),
				testTime.Format(time.RFC3339),
				scanner.Token(),
				ts.Sub(testTime),
			)
		}
	}

	i, error := scanner.ScanInt()
	if i != 11 || error != nil {
		t.Errorf("Test mismatch.")
	}
	ts, error = scanner.ScanTimestamp()
	if error == nil {
		t.Errorf("Expected an error but got %v.", ts)
	}
}
