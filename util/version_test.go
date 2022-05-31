/**
@file          version_test.go
@package       util
@brief         Test the compile version and time functions.
@author        Edward Smith
@date          November 2014
@copyright     -©- Copyright © 2014-2016 Edward Smith, all rights reserved. -©-
*/

package util

import (
	"testing"
	"time"
)

func TestCompileTime(t *testing.T) {
	if CompileVersion() == "0.0.0" {
		t.Errorf("CompileVersion not set.")
	}
	if CompileTimeString() == "compile time not set" {
		t.Errorf("CompileTime not set.")
	}
	compileTime = "2019-08-19-13-22-52Z"
	ct := time.Date(2019, 8, 19, 13, 22, 52, 0, time.UTC)
	if !CompileTime().Equal(ct) {
		t.Errorf("Dates not equal\n%+v\n%s\n%+v.", CompileTime(), compileTime, ct)
	}
}
