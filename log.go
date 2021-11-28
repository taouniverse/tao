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
	Debug(calldepth int, v ...interface{})
	Debugf(calldepth int, format string, v ...interface{})
	Info(calldepth int, v ...interface{})
	Infof(calldepth int, format string, v ...interface{})
	Warn(calldepth int, v ...interface{})
	Warnf(calldepth int, format string, v ...interface{})
	Error(calldepth int, v ...interface{})
	Errorf(calldepth int, format string, v ...interface{})
	Panic(calldepth int, v ...interface{})
	Panicf(calldepth int, format string, v ...interface{})
	Fatal(calldepth int, v ...interface{})
	Fatalf(calldepth int, format string, v ...interface{})
}

var _ Logger = (*logger)(nil)

// logger implements Logger
type logger struct {
	*log.Logger
}

// DEBUG logs info in debug level
func (l *logger) Debug(calldepth int, v ...interface{}) {
	l.Output(calldepth, LevelPrefix[DEBUG]+fmt.Sprintln(v...))
}

// Debugf logs info in debug level
func (l *logger) Debugf(calldepth int, format string, v ...interface{}) {
	l.Output(calldepth, LevelPrefix[DEBUG]+fmt.Sprintf(format, v...))
}

// Info logs info in info level
func (l *logger) Info(calldepth int, v ...interface{}) {
	l.Output(calldepth, LevelPrefix[INFO]+fmt.Sprintln(v...))
}

// Infof logs info in info level
func (l *logger) Infof(calldepth int, format string, v ...interface{}) {
	l.Output(calldepth, LevelPrefix[INFO]+fmt.Sprintf(format, v...))
}

// Warn logs info in warn level
func (l *logger) Warn(calldepth int, v ...interface{}) {
	l.Output(calldepth, LevelPrefix[WARNING]+fmt.Sprintln(v...))
}

// Warnf logs info in warn level
func (l *logger) Warnf(calldepth int, format string, v ...interface{}) {
	l.Output(calldepth, LevelPrefix[WARNING]+fmt.Sprintf(format, v...))
}

// Error logs info in error level
func (l *logger) Error(calldepth int, v ...interface{}) {
	l.Output(calldepth, LevelPrefix[ERROR]+fmt.Sprintln(v...))
}

// Errorf logs info in error level
func (l *logger) Errorf(calldepth int, format string, v ...interface{}) {
	l.Output(calldepth, LevelPrefix[ERROR]+fmt.Sprintf(format, v...))
}

// Panic logs info in panic level
func (l *logger) Panic(calldepth int, v ...interface{}) {
	s := LevelPrefix[PANIC] + fmt.Sprintln(v...)
	l.Output(calldepth, s)
	panic(s)
}

// Panicf logs info in panic level
func (l *logger) Panicf(calldepth int, format string, v ...interface{}) {
	s := LevelPrefix[PANIC] + fmt.Sprintf(format, v...)
	l.Output(calldepth, s)
	panic(s)
}

// Fatal logs info in fatal level
func (l *logger) Fatal(calldepth int, v ...interface{}) {
	l.Output(calldepth, LevelPrefix[FATAL]+fmt.Sprintln(v...))
	os.Exit(1)
}

// Fatalf logs info in fatal level
func (l *logger) Fatalf(calldepth int, format string, v ...interface{}) {
	l.Output(calldepth, LevelPrefix[FATAL]+fmt.Sprintf(format, v...))
	os.Exit(1)
}

// TaoLogger in global
// default to provide based log print
var TaoLogger Logger

// TaoWriter for TaoLogger
var TaoWriter io.Writer

// Debug function wrap of TaoLogger
func Debug(v ...interface{}) {
	TaoLogger.Debug(3, v...)
}

// Debugf function wrap of TaoLogger
func Debugf(format string, v ...interface{}) {
	TaoLogger.Debugf(3, format, v...)
}

// Info function wrap of TaoLogger
func Info(v ...interface{}) {
	TaoLogger.Info(3, v...)
}

// Infof function wrap of TaoLogger
func Infof(format string, v ...interface{}) {
	TaoLogger.Infof(3, format, v...)
}

// Warn function wrap of TaoLogger
func Warn(v ...interface{}) {
	TaoLogger.Warn(3, v...)
}

// Warnf function wrap of TaoLogger
func Warnf(format string, v ...interface{}) {
	TaoLogger.Warnf(3, format, v...)
}

// Error function wrap of TaoLogger
func Error(v ...interface{}) {
	TaoLogger.Error(3, v...)
}

// Errorf function wrap of TaoLogger
func Errorf(format string, v ...interface{}) {
	TaoLogger.Errorf(3, format, v...)
}

// Panic function wrap of TaoLogger
func Panic(v ...interface{}) {
	TaoLogger.Panic(3, v...)
}

// Panicf function wrap of TaoLogger
func Panicf(format string, v ...interface{}) {
	TaoLogger.Panicf(3, format, v...)
}

// Fatal function wrap of TaoLogger
func Fatal(v ...interface{}) {
	TaoLogger.Fatal(3, v...)
}

// Fatalf function wrap of TaoLogger
func Fatalf(format string, v ...interface{}) {
	TaoLogger.Fatalf(3, format, v...)
}

/**
TODO implements Config interface
*/
func taoLogger() {
	TaoWriter = os.Stdout
	TaoLogger = &logger{log.New(TaoWriter, "", log.LstdFlags|log.Lshortfile)}
}
