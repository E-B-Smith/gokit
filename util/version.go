/**
@file          version.go
@package       util
@brief         Compile version and time functions.
@author        Edward Smith
@date          June 2016
@copyright     -©- Copyright © 2014-2016 Edward Smith, all rights reserved. -©-
*/

package util

import "time"

var compileVersion = "0.0.0"
var compileTime = "compile time not set"

// CompileVersion returns the version of the code when it was compiled. This is set in an extra link step.
func CompileVersion() string { return compileVersion }

// CompileTimeString returns the timestamp of the code when the code was compiled. This is set in an extra link step.
func CompileTimeString() string { return compileTime }

// CompileTime returns the timestamp of the code when the code was compiled. This is set in an extra link step.
func CompileTime() time.Time {
	// Mon-Aug-19-13:22:52-PDT-2019
	t, _ := time.Parse("Mon-Jan-2-15:04:05-MST-2006", compileTime)
	return t
}
