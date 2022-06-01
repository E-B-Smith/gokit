/**
@file          log_test.go
@package       log
@brief         Test the log services.
@author        Edward Smith
@date          November 2014
@copyright     -©- Copyright © 2014-2016 Edward Smith, all rights reserved. -©-
*/

package log

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

var tempLogBaseName = ""

func testLogFileName(t *testing.T) string {
	if tempLogBaseName != "" {
		return tempLogBaseName
	}
	tempLogBaseName, err := ioutil.TempDir("", "log_test")
	if err != nil {
		t.Error("Can't make a temporary directory 'log_test'.")
	}
	tempLogBaseName += "/TestLog.log"
	return tempLogBaseName
}

func removeFileExtension(fileName string) string {
	if pos := strings.LastIndexByte(fileName, '.'); pos != -1 {
		return fileName[:pos]
	}
	return fileName
}

func TestLogFileName(t *testing.T) {
	filename := testLogFileName(t)
	if !strings.HasPrefix(filename, "/") || !strings.HasSuffix(filename, "/TestLog.log") {
		t.Errorf("Invalid file name %s.", filename)
	}
}

func TestLogRotation(t *testing.T) {
	SetTeeStderr(false)
	SetLogLevel(LevelAll)
	logRotationInterval = time.Second * 2
	logRetentionCount = 100
	logfilename := testLogFileName(t)
	SetFilename(logfilename)

	Infof("Starting test.")
	sleepTime, _ := time.ParseDuration("0.2s")
	for i := 0; i < 50; i++ {
		Infof("Message %d.", i)
		time.Sleep(sleepTime)
	}
	SetTeeStderr(false)
	SetFilename("")
	globname := removeFileExtension(logfilename)
	logfiles, _ := filepath.Glob(globname + "*")
	if len(logfiles) != 6 {
		t.Errorf("Expected 6 files, found %d.", len(logfiles))
	}
	RemoveLogFiles(t)
}

var channel chan string
var pipeRead, pipeWrite, tempOut *os.File

func BeginPipeToString() {
	tempOut = os.Stderr
	pipeRead, pipeWrite, _ = os.Pipe()
	channel = make(chan string)
	// Copy the output in a separate goroutine so printing can't block indefinitely
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, pipeRead)
		channel <- buf.String()
	}()
	os.Stderr = pipeWrite
}

func EndPipeToString() string {
	os.Stderr = tempOut
	pipeWrite.Close()
	return <-channel
}

func TestLogStackWithError(t *testing.T) {
	SetFilename("log/TestLog.log")
	SetTeeStderr(true)
	BeginPipeToString()
	LogStackWithError(io.ErrUnexpectedEOF)
	s := EndPipeToString()
	SetTeeStderr(false)
	if len(s) < 180 {
		t.Errorf("Expected a stack trace of at least 180 characters. Found %d.", len(s))
		return
	}
	t1 := "Error: 'unexpected EOF'."
	s1 := s[58 : 58+len(t1)]
	if t1 != s1 {
		t.Errorf("Expected '%s' but found '%s'.", t1, s1)
	}
	t1 = "Error: Stack of "
	s1 = s[141 : 141+len(t1)]
	if t1 != s1 {
		t.Errorf("Expected '%s' but found '%s'.", t1, s1)
	}
	RemoveLogFiles(t)
}

func TestPrettyStackString(t *testing.T) {
	s := PrettyStackString(0)
	r := "log.go:340\nlog_test.go:122\ntesting.go:1439\nasm_amd64.s:1571\n"
	if s != r {
		t.Errorf("Expected\n%s\nbut found\n%s.", r, s)
	}
}

func RemoveLogFiles(t *testing.T) {
	SetTeeStderr(false)
	SetFilename("")

	logfilename := testLogFileName(t)
	logfilename = filepath.Dir(logfilename)
	logfiles, _ := filepath.Glob(logfilename + "/*")
	for _, file := range logfiles {
		error := os.Remove(file)
		if error != nil {
			t.Errorf("Can't remove log file '%s': %v.", file, error)
		}
	}
	error := os.Remove(logfilename)
	if error != nil {
		t.Errorf("Can't remove log directory '%s': %v.", logfilename, error)
	}
}
