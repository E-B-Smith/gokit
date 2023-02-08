/**
@file          util.go
@package       log
@brief         Log utility methods.
@author        Edward Smith
@date          November 2023
@copyright     -©- Copyright © 2014-2023 Edward Smith, all rights reserved. -©-
*/

package log

import (
	"os"
	"os/user"
	"path"
	"strings"
)

// HomePath returns the path to the user's home directory.
func homePath() string {
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
func absolutePath(filename string) string {
	filename = strings.TrimSpace(filename)
	if filename == "~" || strings.HasPrefix(filename, "~/") {
		filename = strings.TrimPrefix(filename, "~")
		filename = path.Join(homePath(), filename)
	}
	if !path.IsAbs(filename) {
		s, _ := os.Getwd()
		filename = path.Join(s, filename)
	}
	filename = path.Clean(filename)
	return filename
}
