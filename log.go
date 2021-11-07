// Copyright 2021
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package tao

import (
	"fmt"
	"io"
	"log"
	"os"
)

// LogLevel log's level
type LogLevel uint8

const (
	// DEBUG (usually) is used in development env to print track info but disabled in production env to avoid overweight logs
	DEBUG LogLevel = iota
	// INFO (usually) is default level to print some core infos
	INFO
	// WARNING should be mentioned, it's more important than INFO
	WARNING
	// ERROR must be solved, program shouldn't generate any error-level logs.
	ERROR
	// PANIC logs a message, then panics.
	PANIC
	// FATAL logs a message, then calls os.Exit(1).
	FATAL
)

// LevelPrefix to define log prefix of log level
var LevelPrefix = map[LogLevel]string{
	DEBUG:   "[D] ",
	INFO:    "[I] ",
	WARNING: "[W] ",
	ERROR:   "[E] ",
	PANIC:   "[P] ",
	FATAL:   "[F] ",
}

// Logger in tao
type Logger interface {
	Debug(v ...interface{})
	Debugf(format string, v ...interface{})
	Info(v ...interface{})
	Infof(format string, v ...interface{})
	Warn(v ...interface{})
	Warnf(format string, v ...interface{})
	Error(v ...interface{})
	Errorf(format string, v ...interface{})
	Panic(v ...interface{})
	Panicf(format string, v ...interface{})
	Fatal(v ...interface{})
	Fatalf(format string, v ...interface{})
}

var _ Logger = (*logger)(nil)

// logger implements Logger
type logger struct {
	*log.Logger
}

// DEBUG logs info in debug level
func (l *logger) Debug(v ...interface{}) {
	l.Output(2, LevelPrefix[DEBUG]+fmt.Sprintln(v...))
}

// Debugf logs info in debug level
func (l *logger) Debugf(format string, v ...interface{}) {
	l.Output(2, LevelPrefix[DEBUG]+fmt.Sprintf(format, v...))
}

// Info logs info in info level
func (l *logger) Info(v ...interface{}) {
	l.Output(2, LevelPrefix[INFO]+fmt.Sprintln(v...))
}

// Infof logs info in info level
func (l *logger) Infof(format string, v ...interface{}) {
	l.Output(2, LevelPrefix[INFO]+fmt.Sprintf(format, v...))
}

// Warn logs info in warn level
func (l *logger) Warn(v ...interface{}) {
	l.Output(2, LevelPrefix[WARNING]+fmt.Sprintln(v...))
}

// Warnf logs info in warn level
func (l *logger) Warnf(format string, v ...interface{}) {
	l.Output(2, LevelPrefix[WARNING]+fmt.Sprintf(format, v...))
}

// Error logs info in error level
func (l *logger) Error(v ...interface{}) {
	l.Output(2, LevelPrefix[ERROR]+fmt.Sprintln(v...))
}

// Errorf logs info in error level
func (l *logger) Errorf(format string, v ...interface{}) {
	l.Output(2, LevelPrefix[ERROR]+fmt.Sprintf(format, v...))
}

// Panic logs info in panic level
func (l *logger) Panic(v ...interface{}) {
	s := LevelPrefix[PANIC] + fmt.Sprintln(v...)
	l.Output(2, s)
	panic(s)
}

// Panicf logs info in panic level
func (l *logger) Panicf(format string, v ...interface{}) {
	s := LevelPrefix[PANIC] + fmt.Sprintf(format, v...)
	l.Output(2, s)
	panic(s)
}

// Fatal logs info in fatal level
func (l *logger) Fatal(v ...interface{}) {
	l.Output(2, LevelPrefix[FATAL]+fmt.Sprintln(v...))
	os.Exit(1)
}

// Fatalf logs info in fatal level
func (l *logger) Fatalf(format string, v ...interface{}) {
	l.Output(2, LevelPrefix[FATAL]+fmt.Sprintf(format, v...))
	os.Exit(1)
}

// TaoLogger in global
// default to provide based log print
var TaoLogger Logger

// Writer for TaoLogger
var Writer io.Writer

// init default logger & writer
func init() {
	Writer = os.Stdout
	TaoLogger = &logger{log.New(Writer, " ", log.LstdFlags|log.Lshortfile)}
}

/**
SPECIAL:
return function instead of calling function
to keep the fileName & lineNumber same
*/

// Debug function of TaoLogger
func Debug() func(v ...interface{}) {
	return TaoLogger.Debug
}

// Debugf function of TaoLogger
func Debugf() func(format string, v ...interface{}) {
	return TaoLogger.Debugf
}

// Info function of TaoLogger
func Info() func(v ...interface{}) {
	return TaoLogger.Info
}

// Infof function of TaoLogger
func Infof() func(format string, v ...interface{}) {
	return TaoLogger.Infof
}

// Warn function of TaoLogger
func Warn() func(v ...interface{}) {
	return TaoLogger.Warn
}

// Warnf function of TaoLogger
func Warnf() func(format string, v ...interface{}) {
	return TaoLogger.Warnf
}

// Error function of TaoLogger
func Error() func(v ...interface{}) {
	return TaoLogger.Error
}

// Errorf function of TaoLogger
func Errorf() func(format string, v ...interface{}) {
	return TaoLogger.Errorf
}

// Panic function of TaoLogger
func Panic() func(v ...interface{}) {
	return TaoLogger.Panic
}

// Panicf function of TaoLogger
func Panicf() func(format string, v ...interface{}) {
	return TaoLogger.Panicf
}

// Fatal function of TaoLogger
func Fatal() func(v ...interface{}) {
	return TaoLogger.Fatal
}

// Fatalf function of TaoLogger
func Fatalf() func(format string, v ...interface{}) {
	return TaoLogger.Fatalf
}
