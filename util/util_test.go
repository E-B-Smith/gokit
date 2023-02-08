/**
@file          util_test.go
@package       util
@brief         Utility tests.
@author        Edward Smith
@date          February, 2023
@copyright     -©- Copyright © 2014-2023 Edward Smith. All rights reserved. -©-
*/

package util

import (
	"os"
	"testing"
)

func TestHomePath(t *testing.T) {
	home := os.Getenv("HOME")
	if len(home) <= 0 {
		t.Errorf("Home environment variable is empty.")
		return
	}
	if home != HomePath() {
		t.Errorf("Got '%s' but wanted '%s'.", HomePath(), home)
	}
}

func TestAbsolutePath(t *testing.T) {
	workingPath, _ := os.Getwd()
	homePath := HomePath()
	tests := []struct {
		testin, testout string
	}{
		{" /Absolute/File/Name  	", "/Absolute/File/Name"},
		{"~", HomePath()},
		{"~/", HomePath()},
		{"~/Home", homePath + "/Home"},
		{"  ~/Home  ", homePath + "/Home"},
		{"  ~ Home  ", workingPath + "/~ Home"},
		{"/1/~/Home  ", "/1/~/Home"},
		{"Relative", workingPath + "/Relative"},
		{"/a/b/c", "/a/b/c"},
		{"", workingPath},
	}

	for i, test := range tests {
		result := AbsolutePath(test.testin)
		if result != test.testout {
			t.Errorf("Case %d: Got %s but want %s.", i, result, test.testout)
		}
	}
}
