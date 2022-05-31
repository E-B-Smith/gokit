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
	// Like: 2020-02-03-16-02-44-0800 for January 3, 2020 4:02:44 PST
	t, _ := time.Parse("2006-01-02-15-04-05Z0700", compileTime)
	return t
}
