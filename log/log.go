/**
@file          log.go
@package       log
@brief         Basic log services.
@author        Edward Smith
@date          November 2014
@copyright     -©- Copyright © 2014-2016 Edward Smith, all rights reserved. -©-
*/

/*
Package log provides simple logging with a standardized message format to a log file and optionally Stdout.

The log files are automatically rotated once a day and the oldest log file is removed.
*/
package log

import (
	"fmt"
	"io"
	"math"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"sync"
	"syscall"
	"time"
	"unicode"
)

func debugMessagef(format string, args ...interface{}) {
	// For debugging the log.
	//var message = fmt.Sprintf(format, args...)
	//fmt.Fprintf(os.Stderr, "%s\n", message)
}

// Level is the type used for indicating log message severity.
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

	// LevelStart indicates an app start message.
	LevelStart

	// LevelExit indicates an app exit message.
	LevelExit

	// LevelWarning indicates a warning message.
	LevelWarning

	// LevelError indicates an error message.
	LevelError

	// LevelNone indicates no log messages will be written to the log.
	LevelNone
)

var levelNames = []string{
	"LevelInvalid",
	"LevelAll",
	"LevelDebug",
	"LevelInfo",
	"LevelStart",
	"LevelExit",
	"LevelWarning",
	"LevelError",
	"LevelNone",
}

// LevelFromString returns a log Level const from a string representing the const name.
func LevelFromString(s string) Level {
	for index := range levelNames {
		if s == levelNames[index] {
			return Level(index)
		}
	}
	return LevelInvalid
}

// StringFromLevel returns a string representing the passed log Level const.
func StringFromLevel(level Level) string {
	if level < LevelInvalid || level > LevelNone {
		return levelNames[LevelInvalid]
	}
	return levelNames[level]
}

var (
	mutex     = &sync.RWMutex{}
	logLevel  = LevelInfo
	teeStderr bool

	// LogRotationInterval sets how often the log file will be rotated.
	logRotationInterval = time.Hour * 24.0

	// Number of old log files to keep
	logRetentionCount = 1

	logWriter       io.WriteCloser = os.Stderr
	logFilename     string
	logRotationTime time.Time
)

// SetLogLevel sets the minimum log severity level written to the log.
func SetLogLevel(level Level) {
	mutex.Lock()
	defer mutex.Unlock()
	logLevel = level
}

// LogLevel returns the current log level severity level being written to the log.
func LogLevel() Level {
	mutex.RLock()
	defer mutex.RUnlock()
	return logLevel
}

// SetTeeStderr when set to true log messages are output to Stderr as well as the log file.
func SetTeeStderr(value bool) {
	mutex.Lock()
	defer mutex.Unlock()
	teeStderr = value
}

// TeeStderr returns the current state of TeeStderr. If TeeStderr is true log messages are output to Stderr as well as the log file.
func TeeStderr() bool {
	mutex.RLock()
	defer mutex.RUnlock()
	return teeStderr
}

func closeLogFile() {
	if logWriter != os.Stderr && logWriter != os.Stdout {
		logWriter.Close()
	}
}

func openLogFile() {
	logRotationTime = time.Unix(math.MaxInt64-1000, 0) //  Distant future

	defer func() {
		if reason := recover(); reason != nil {
			logFilename = ""
			logWriter = os.Stderr
			Errorf("%s", reason)
		}
		name := logFilename
		if name == "" {
			name = "Stderr"
		}
		// fmt.Fprintf(os.Stderr, "Log file is '%s'.\n", name)
	}()

	debugMessagef("Logfile: '%s'.", logFilename)

	logFilename = strings.TrimSpace(logFilename)
	if logFilename == "" {
		logWriter = os.Stderr
		return
	}

	logFilename = absolutePath(logFilename)
	if len(logFilename) <= 0 {
		logWriter = os.Stderr
		return
	}

	debugMessagef("Logfile: '%s'.", logFilename)

	var error error
	pathname := filepath.Dir(logFilename)
	if len(pathname) > 0 {
		if error = os.MkdirAll(pathname, 0700); error != nil {
			logWriter = os.Stderr
			panic(fmt.Sprintf("Can't create directory for log file '%s': %v.", logFilename, error))
		}
	}

	debugMessagef("Logfile: '%s'.", pathname)

	var flags = syscall.O_APPEND | syscall.O_CREAT | syscall.O_WRONLY
	var mode = os.ModeAppend | 0700

	debugMessagef("Logfile: '%s'.", logFilename)

	logWriter, error = os.OpenFile(logFilename, flags, mode)
	if error != nil {
		logWriter = os.Stderr
		panic(fmt.Sprintf("Can't open log file '%s' for writing: %v.", logFilename, error))
	}

	if logRotationInterval.Seconds() > 0 {
		var nextTime = (int64(time.Now().Unix()) / int64(logRotationInterval.Seconds())) + 1
		nextTime *= int64(logRotationInterval.Seconds())
		logRotationTime = time.Unix(nextTime, 0)
	}
}

func rotateLogFile() {
	if len(logFilename) <= 0 {
		return
	}

	defer func() {
		if reason := recover(); reason != nil {
			logFilename = ""
			logWriter = os.Stderr
			Errorf("%s", reason)
		}
		name := logFilename
		if name == "" {
			name = "Stderr"
		}
		// fmt.Fprintf(os.Stderr, "Log file is '%s'.\n", name)
	}()

	//  Create a new file for the log --

	baseName := filepath.Base(logFilename)
	ext := filepath.Ext(baseName)
	if len(ext) != 0 {
		baseName = strings.TrimSuffix(baseName, ext)
	}
	replacePunct := func(r rune) rune {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			return r
		}
		return '-'
	}
	timeString := strings.Map(replacePunct, logRotationTime.Format(time.RFC3339)) + ext
	newBase := fmt.Sprintf("%s-%s", baseName, timeString)
	newPath := filepath.Join(filepath.Dir(logFilename), newBase)
	closeLogFile()
	error := os.Rename(logFilename, newPath)
	if error != nil {
		panic(error)
	}
	openLogFile()
	Infof("Log rotated to '%s'.", newPath)
	Infof("Log continues in '%s'.", logFilename)

	//  Delete the oldest --

	globPath := filepath.Join(filepath.Dir(logFilename), baseName+"-*")
	logfiles, error := filepath.Glob(globPath)
	debugMessagef("Log files: %+v.", logfiles)
	if error != nil {
		Error(error)
		return
	}

	//  Keep the newest logRetentionCount --

	sortedLogfiles := sort.StringSlice(logfiles)
	sortedLogfiles.Sort()
	for i := 0; i < len(sortedLogfiles)-logRetentionCount; i++ {
		Infof("Removing old log '%s'.", sortedLogfiles[i])
		error = os.Remove(sortedLogfiles[i])
		if error != nil {
			Errorf("Can't remove log file '%s': %v.", sortedLogfiles[i], error)
		}
	}
}

// SetFilename sets the filename for the log file. If filename is the empty string log goes to stdout.
func SetFilename(filename string) {
	mutex.Lock()
	if len(filename) > 0 {
		filename = absolutePath(filename)
	}
	if filename == logFilename {
		mutex.Unlock()
		return
	}
	closeLogFile()
	logFilename = filename
	mutex.Unlock()
	openLogFile()
}

// Filename returns the current log file name.
func Filename() string {
	mutex.RLock()
	defer mutex.RUnlock()
	return logFilename
}

// LogStackWithError writes an error and a stack trace to the log.
func LogStackWithError(error interface{}) {
	trace := make([]byte, 64000)
	count := runtime.Stack(trace, false)
	s := trace[:count]
	logRaw(LevelError, 2, "'%v'.", error)
	logRaw(LevelError, 2, "Stack of %d bytes: %s.", count, s)
}

// PrettyStackString returns a prettyfied string of the current stack. The `skip` parameter indicates
// the number of frames to skip before reporting.
func PrettyStackString(skip int) string {
	var result string
	_, filename, linenumber, ok := runtime.Caller(skip)
	for ok {
		filename = path.Base(filename)
		i := len(filename)
		if i > 26 {
			i = 26
		}
		result += fmt.Sprintf("%s:%d\n", filename[:i], linenumber)
		skip++
		_, filename, linenumber, ok = runtime.Caller(skip)
	}
	return result
}

// CurrentFunctionName returns the current function name.
func CurrentFunctionName() string {
	pc, _, _, _ := runtime.Caller(1)
	funcname := path.Base(runtime.FuncForPC(pc).Name())
	return funcname
}

// LogFunctionName logs the current function name to the log.
func LogFunctionName() {
	pc, _, _, _ := runtime.Caller(1)
	funcname := path.Base(runtime.FuncForPC(pc).Name())
	logRaw(LevelDebug, 2, "Function %s.", funcname)
}

// FlushMessages flushes all outstanding log messages to the log file.
func FlushMessages() {
	mutex.Lock()
	defer mutex.Unlock()
	closeLogFile()
	openLogFile()
}

// logRaw logs a raw messgae.
// stackDepth is the depth in the stack to where the calling source code / line number should be billed.
func logRaw(logLevel Level, stackDepth int, format string, args ...interface{}) {

	LevelNames := []string{
		"Inval",
		"  All",
		"Debug",
		" Info",
		"Start",
		" Exit",
		" Warn",
		"Error",
		" None",
	}

	if logLevel < LogLevel() {
		return
	}
	if logLevel < LevelDebug || logLevel > LevelError {
		logLevel = LevelError
	}

	var dirname string
	_, filename, linenumber, _ := runtime.Caller(stackDepth)
	dirname, filename = path.Split(filename)
	dirname = path.Base(dirname)
	i := len(filename)
	if i > 26 {
		i = 26
	}
	filename = dirname + "/" + filename[:i]

	itemTime := time.Now()
	if itemTime.After(logRotationTime) {
		rotateLogFile()
	}

	var message = fmt.Sprintf(format, args...)
	message = strings.Replace(message, "\n", "|", -1)
	message = strings.Replace(message, "\r", "|", -1)
	message = fmt.Sprintf(
		"%s %26s:%-4d %s: %s\n",
		itemTime.Format(time.RFC3339),
		filename,
		linenumber,
		LevelNames[logLevel],
		message,
	)
	fmt.Fprint(logWriter, message)
	if TeeStderr() {
		fmt.Fprint(os.Stderr, message)
	}
}

// Debugf writes a debug level message to the log.
func Debugf(format string, args ...interface{}) { logRaw(LevelDebug, 2, format, args...) }

// Startf writes a start level message to the log.
func Startf(format string, args ...interface{}) { logRaw(LevelStart, 2, format, args...) }

// Exitf writes an exit level message to the log.
func Exitf(format string, args ...interface{}) { logRaw(LevelExit, 2, format, args...) }

// Infof writes an info level message to the log.
func Infof(format string, args ...interface{}) { logRaw(LevelInfo, 2, format, args...) }

// Warningf writes a warning level message to the log.
func Warningf(format string, args ...interface{}) { logRaw(LevelWarning, 2, format, args...) }

// Errorf writes a error level message to the log.
func Errorf(format string, args ...interface{}) { logRaw(LevelError, 2, format, args...) }

// Error Writes a error message to the log.
func Error(error error) { logRaw(LevelError, 2, "%v.", error) }
