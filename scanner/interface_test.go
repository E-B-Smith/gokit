/**
@file          interface_test.go
@package       scanner
@brief         Test the `interface` parser.
@author        Edward Smith
@date          March 2016
@copyright     -©- Copyright © 2014-2016 Edward Smith, all rights reserved. -©-
*/

package scanner

import (
	"testing"
)

type Level int32

const (
	// LevelInvalid indicates that an invalid level was set.
	LevelInvalid Level = iota

	// LevelAll indicates all messages will be logged.
	LevelAll

	// LevelDebug indicates a debug message.
	LevelDebug

	// LevelInfo indicates an Info level message.
	LevelInfo

	// LevelMax is the max value
	LevelMax
)

// SubStruct for test.
type SubStruct struct {
	S1 string
	I1 int
}

// TestStruct for test.
type TestStruct struct {
	B          bool
	I8         int8
	I16        int16
	I32        int32
	I64        int64
	I          int
	F32        float32
	F64        float64
	S          string
	EnumTest   Level `enum:"LevelInvalid,LevelAll,LevelDebug,LevelInfo,LevelMax"`
	SubStruct  SubStruct
	LastString string
}

const testString = `
b   true
i8  1
i16 2
i32 3
i64 4
i   5
f32 0.6
f64 0.7
s   "This is a string."
enum_test LevelDebug
sub_struct {
    s1  "s1 string"
    i1  42
}
last-string Final
`

func TestInterface(t *testing.T) {
	var ts TestStruct
	var ss = SubStruct{
		S1: "s1 string",
		I1: 42,
	}
	scanner := NewScannerWithString(testString)
	error := scanner.ScanInterface(&ts)
	if error != nil {
		t.Errorf("Error %v.", error)
	}
	if ts.B &&
		ts.I8 == 1 &&
		ts.I16 == 2 &&
		ts.I32 == 3 &&
		ts.I64 == 4 &&
		ts.I == 5 &&
		ts.F32 == 0.6 &&
		ts.F64 == 0.7 &&
		ts.S == "This is a string." &&
		ts.EnumTest == LevelDebug &&
		ts.SubStruct == ss &&
		ts.LastString == "Final" {
	} else {
		t.Errorf("TestStruct failed: %+v", ts)
	}
}
