// Copyright 2021 huija
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
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
)

// Log config in tao
type Log struct {
	Level     LogLevel `json:"level"`
	Type      LogType  `json:"type"`
	CallDepth int      `json:"callDepth"`
	Path      string   `json:"path,omitempty"`
}

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

// String for LogLevel Config
func (l LogLevel) String() string {
	switch l {
	case DEBUG:
		return "debug"
	case INFO:
		return "info"
	case WARNING:
		return "warning"
	case ERROR:
		return "error"
	case PANIC:
		return "panic"
	case FATAL:
		return "fatal"
	default:
		return fmt.Sprintf("tao.LogLevel(%d)", l)
	}
}

// MarshalText instead of number
func (l LogLevel) MarshalText() ([]byte, error) {
	return []byte(l.String()), nil
}

// UnmarshalText to number
func (l *LogLevel) UnmarshalText(text []byte) error {
	if l == nil {
		return errors.New("log: can't unmarshal a nil *LogLevel")
	}
	switch lower := string(bytes.ToLower(text)); lower {
	case "debug":
		*l = DEBUG
	case "info":
		*l = INFO
	case "warning":
		*l = WARNING
	case "error":
		*l = ERROR
	case "panic":
		*l = PANIC
	case "fatal":
		*l = FATAL
	default:
		return fmt.Errorf("log: unrecognized LogLevel: %q", lower)
	}
	return nil
}

// LogType log's type
type LogType uint8

const (
	// Console log
	Console LogType = 1 // 0b1
	// File log
	File LogType = 2 // 0b10
)

// String for LogType Config
func (l LogType) String() string {
	switch l {
	case Console:
		return "console"
	case File:
		return "file"
	case Console | File:
		return "console|file"
	default:
		return fmt.Sprintf("tao.LOGTYPE(%d)", l)
	}
}

// MarshalText instead of number
func (l LogType) MarshalText() ([]byte, error) {
	return []byte(l.String()), nil
}

// UnmarshalText to number
func (l *LogType) UnmarshalText(text []byte) error {
	if l == nil {
		return errors.New("log: can't unmarshal a nil *LogType")
	}
	switch lower := string(bytes.ToLower(text)); lower {
	case "console":
		*l = Console
	case "file":
		*l = File
	case "console|file", "file|console":
		*l = File | Console
	default:
		return fmt.Errorf("log: unrecognized LogType: %q", lower)
	}
	return nil
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

// logger implements Logger using standard lib
type logger struct {
	*log.Logger

	calldepth int
}

// levelPrefix to define log prefix of log level
var levelPrefix = map[LogLevel]string{
	DEBUG:   "[D] ",
	INFO:    "[I] ",
	WARNING: "[W] ",
	ERROR:   "[E] ",
	PANIC:   "[P] ",
	FATAL:   "[F] ",
}

// Debug logs info in debug level
func (l *logger) Debug(v ...interface{}) {
	if t.Level > DEBUG {
		return
	}
	l.Output(l.calldepth, levelPrefix[DEBUG]+fmt.Sprintln(v...))
}

// Debugf logs info in debug level
func (l *logger) Debugf(format string, v ...interface{}) {
	if t.Level > DEBUG {
		return
	}
	l.Output(l.calldepth, levelPrefix[DEBUG]+fmt.Sprintf(format, v...))
}

// Info logs info in info level
func (l *logger) Info(v ...interface{}) {
	if t.Level > INFO {
		return
	}
	l.Output(l.calldepth, levelPrefix[INFO]+fmt.Sprintln(v...))
}

// Infof logs info in info level
func (l *logger) Infof(format string, v ...interface{}) {
	if t.Level > INFO {
		return
	}
	l.Output(l.calldepth, levelPrefix[INFO]+fmt.Sprintf(format, v...))
}

// Warn logs info in warn level
func (l *logger) Warn(v ...interface{}) {
	if t.Level > WARNING {
		return
	}
	l.Output(l.calldepth, levelPrefix[WARNING]+fmt.Sprintln(v...))
}

// Warnf logs info in warn level
func (l *logger) Warnf(format string, v ...interface{}) {
	if t.Level > WARNING {
		return
	}
	l.Output(l.calldepth, levelPrefix[WARNING]+fmt.Sprintf(format, v...))
}

// Error logs info in error level
func (l *logger) Error(v ...interface{}) {
	if t.Level > ERROR {
		return
	}
	l.Output(l.calldepth, levelPrefix[ERROR]+fmt.Sprintln(v...))
}

// Errorf logs info in error level
func (l *logger) Errorf(format string, v ...interface{}) {
	if t.Level > ERROR {
		return
	}
	l.Output(l.calldepth, levelPrefix[ERROR]+fmt.Sprintf(format, v...))
}

// Panic logs info in panic level
func (l *logger) Panic(v ...interface{}) {
	if t.Level > PANIC {
		return
	}
	s := levelPrefix[PANIC] + fmt.Sprintln(v...)
	l.Output(l.calldepth, s)
	panic(s)
}

// Panicf logs info in panic level
func (l *logger) Panicf(format string, v ...interface{}) {
	if t.Level > PANIC {
		return
	}
	s := levelPrefix[PANIC] + fmt.Sprintf(format, v...)
	l.Output(l.calldepth, s)
	panic(s)
}

// Fatal logs info in fatal level
func (l *logger) Fatal(v ...interface{}) {
	if t.Level > FATAL {
		return
	}
	l.Output(l.calldepth, levelPrefix[FATAL]+fmt.Sprintln(v...))
	os.Exit(1)
}

// Fatalf logs info in fatal level
func (l *logger) Fatalf(format string, v ...interface{}) {
	if t.Level > FATAL {
		return
	}
	l.Output(l.calldepth, levelPrefix[FATAL]+fmt.Sprintf(format, v...))
	os.Exit(1)
}

// Closed this logger
func (l *logger) Close() error {
	return nil
}

// taoLogger for tao
type taoLogger struct {
	mu sync.Mutex

	loggers map[string]Logger
	writers map[string]io.Writer
}

// globalLogger which default to provide based log print
var globalLogger = new(taoLogger)

// GetWriter in tao
func GetWriter(configKey string) io.Writer {
	return globalLogger.writers[configKey]
}

// SetWriter to tao
func SetWriter(configKey string, w io.Writer) error {
	globalLogger.mu.Lock()
	defer globalLogger.mu.Unlock()

	if globalLogger.writers == nil {
		globalLogger.writers = make(map[string]io.Writer)
	}

	if _, ok := globalLogger.writers[configKey]; ok {
		return NewError(DuplicateCall, "log: %s's writer has been set before", configKey)
	}

	globalLogger.writers[configKey] = w
	return nil
}

// DeleteWriter of tao
func DeleteWriter(configKey string) error {
	globalLogger.mu.Lock()
	defer globalLogger.mu.Unlock()

	writer, ok := globalLogger.writers[configKey]
	if !ok {
		return NewError(ParamInvalid, "log: %s's writer not set", configKey)
	}
	delete(globalLogger.writers, configKey)

	// writer close
	if l, ok := writer.(io.Closer); ok {
		return l.Close()
	}
	return nil
}

// GetLogger in tao
func GetLogger(configKey string) Logger {
	return globalLogger.loggers[configKey]
}

// SetLogger to tao
func SetLogger(configKey string, logger Logger) error {
	globalLogger.mu.Lock()
	defer globalLogger.mu.Unlock()

	if globalLogger.loggers == nil {
		globalLogger.loggers = make(map[string]Logger)
	}

	if _, ok := globalLogger.loggers[configKey]; ok {
		return NewError(DuplicateCall, "log: %s's logger has been set before", configKey)
	}

	globalLogger.loggers[configKey] = logger
	return nil
}

// DeleteLogger of tao
func DeleteLogger(configKey string) error {
	globalLogger.mu.Lock()
	defer globalLogger.mu.Unlock()

	logger, ok := globalLogger.loggers[configKey]
	if !ok {
		return NewError(ParamInvalid, "log: %s's logger not set", configKey)
	}
	delete(globalLogger.loggers, configKey)

	// logger close
	if l, ok := logger.(io.Closer); ok {
		return l.Close()
	}
	return nil
}

// Debug function wrap of taoLogger
func Debug(v ...interface{}) {
	for _, l := range globalLogger.loggers {
		l.Debug(v...)
	}
}

// Debugf function wrap of taoLogger
func Debugf(format string, v ...interface{}) {
	for _, l := range globalLogger.loggers {
		l.Debugf(format, v...)
	}
}

// Info function wrap of taoLogger
func Info(v ...interface{}) {
	for _, l := range globalLogger.loggers {
		l.Info(v...)
	}
}

// Infof function wrap of taoLogger
func Infof(format string, v ...interface{}) {
	for _, l := range globalLogger.loggers {
		l.Infof(format, v...)
	}
}

// Warn function wrap of taoLogger
func Warn(v ...interface{}) {
	for _, l := range globalLogger.loggers {
		l.Warn(v...)
	}
}

// Warnf function wrap of taoLogger
func Warnf(format string, v ...interface{}) {
	for _, l := range globalLogger.loggers {
		l.Warnf(format, v...)
	}
}

// Error function wrap of taoLogger
func Error(v ...interface{}) {
	for _, l := range globalLogger.loggers {
		l.Error(v...)
	}
}

// Errorf function wrap of taoLogger
func Errorf(format string, v ...interface{}) {
	for _, l := range globalLogger.loggers {
		l.Errorf(format, v...)
	}
}

// Panic function wrap of taoLogger
func Panic(v ...interface{}) {
	for _, l := range globalLogger.loggers {
		l.Panic(v...)
	}
}

// Panicf function wrap of taoLogger
func Panicf(format string, v ...interface{}) {
	for _, l := range globalLogger.loggers {
		l.Panicf(format, v...)
	}
}

// Fatal function wrap of taoLogger
func Fatal(v ...interface{}) {
	for _, l := range globalLogger.loggers {
		l.Fatal(v...)
	}
}

// Fatalf function wrap of taoLogger
func Fatalf(format string, v ...interface{}) {
	for _, l := range globalLogger.loggers {
		l.Fatalf(format, v...)
	}
}
