/**
@file          util.go
@package       util
@brief         Grab bag of utility functions.
@author        Edward Smith
@date          November, 2014
@copyright     -©- Copyright © 2014-2019 Edward Smith. All rights reserved. -©-
*/

/*
Package util is a grab bag of frequently needed utility functions that aren't weighty enough to merit
a whole package but are useful to have.
*/
package util

import (
	"os"
	"os/user"
	"path"
	"path/filepath"
	"strings"
)

// HomePath returns the path to the user's home directory.
func HomePath() string {
	homepath := ""
	u, error := user.Current()
	if error == nil {
		homepath = u.HomeDir
	}
	if len(homepath) <= 0 {
		homepath = os.Getenv("HOME")
	}
	return homepath
}

// AbsolutePath expands the filename to the absolute path of the file.
func AbsolutePath(filename string) string {
	filename = strings.TrimSpace(filename)
	if filename == "~" || strings.HasPrefix(filename, "~/") {
		filename = strings.TrimPrefix(filename, "~")
		filename = path.Join(HomePath(), filename)
	}
	if !path.IsAbs(filename) {
		s, _ := os.Getwd()
		filename = path.Join(s, filename)
	}
	filename = path.Clean(filename)
	return filename
}

// ResourcePath returns a path to an app resource.
func ResourcePath(pathTo ...string) string {
	p, error := os.Executable()
	if error != nil || len(p) == 0 {
		panic(error)
	}
	p = filepath.Dir(p)
	fullpath := append(make([]string, 0), p)
	fullpath = append(fullpath, pathTo...)
	p = filepath.Join(fullpath...)
	return path.Clean(p)
}
