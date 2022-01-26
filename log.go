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
	Level LogLevel `json:"level"`
	Type  LogType  `json:"type"`
	Path  string   `json:"path,omitempty"`
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
	Console LogType = 1 // 0b1
	File    LogType = 2 // 0b10
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

// logger implements Logger using standard lib
type logger struct {
	*log.Logger
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
func (l *logger) Debug(calldepth int, v ...interface{}) {
	if t.Level > DEBUG {
		return
	}
	l.Output(calldepth, levelPrefix[DEBUG]+fmt.Sprintln(v...))
}

// Debugf logs info in debug level
func (l *logger) Debugf(calldepth int, format string, v ...interface{}) {
	if t.Level > DEBUG {
		return
	}
	l.Output(calldepth, levelPrefix[DEBUG]+fmt.Sprintf(format, v...))
}

// Info logs info in info level
func (l *logger) Info(calldepth int, v ...interface{}) {
	if t.Level > INFO {
		return
	}
	l.Output(calldepth, levelPrefix[INFO]+fmt.Sprintln(v...))
}

// Infof logs info in info level
func (l *logger) Infof(calldepth int, format string, v ...interface{}) {
	if t.Level > INFO {
		return
	}
	l.Output(calldepth, levelPrefix[INFO]+fmt.Sprintf(format, v...))
}

// Warn logs info in warn level
func (l *logger) Warn(calldepth int, v ...interface{}) {
	if t.Level > WARNING {
		return
	}
	l.Output(calldepth, levelPrefix[WARNING]+fmt.Sprintln(v...))
}

// Warnf logs info in warn level
func (l *logger) Warnf(calldepth int, format string, v ...interface{}) {
	if t.Level > WARNING {
		return
	}
	l.Output(calldepth, levelPrefix[WARNING]+fmt.Sprintf(format, v...))
}

// Error logs info in error level
func (l *logger) Error(calldepth int, v ...interface{}) {
	if t.Level > ERROR {
		return
	}
	l.Output(calldepth, levelPrefix[ERROR]+fmt.Sprintln(v...))
}

// Errorf logs info in error level
func (l *logger) Errorf(calldepth int, format string, v ...interface{}) {
	if t.Level > ERROR {
		return
	}
	l.Output(calldepth, levelPrefix[ERROR]+fmt.Sprintf(format, v...))
}

// Panic logs info in panic level
func (l *logger) Panic(calldepth int, v ...interface{}) {
	if t.Level > PANIC {
		return
	}
	s := levelPrefix[PANIC] + fmt.Sprintln(v...)
	l.Output(calldepth, s)
	panic(s)
}

// Panicf logs info in panic level
func (l *logger) Panicf(calldepth int, format string, v ...interface{}) {
	if t.Level > PANIC {
		return
	}
	s := levelPrefix[PANIC] + fmt.Sprintf(format, v...)
	l.Output(calldepth, s)
	panic(s)
}

// Fatal logs info in fatal level
func (l *logger) Fatal(calldepth int, v ...interface{}) {
	if t.Level > FATAL {
		return
	}
	l.Output(calldepth, levelPrefix[FATAL]+fmt.Sprintln(v...))
	os.Exit(1)
}

// Fatalf logs info in fatal level
func (l *logger) Fatalf(calldepth int, format string, v ...interface{}) {
	if t.Level > FATAL {
		return
	}
	l.Output(calldepth, levelPrefix[FATAL]+fmt.Sprintf(format, v...))
	os.Exit(1)
}

// Close this logger
func (l *logger) Close() error {
	return nil
}

// TaoLogger for tao
type TaoLogger struct {
	mu sync.Mutex

	loggers map[string]Logger
	writers map[string]io.Writer
}

// taoLogger in global
// default to provide based log print
var taoLogger = new(TaoLogger)

// GetWriter in tao
func GetWriter(configKey string) io.Writer {
	return taoLogger.writers[configKey]
}

// SetWriter to tao
func SetWriter(configKey string, w io.Writer) error {
	taoLogger.mu.Lock()
	defer taoLogger.mu.Unlock()

	if taoLogger.writers == nil {
		taoLogger.writers = make(map[string]io.Writer)
	}

	if _, ok := taoLogger.writers[configKey]; ok {
		return NewError(DuplicateCall, "log: %s's writer has been set before", configKey)
	}

	taoLogger.writers[configKey] = w
	return nil
}

// DeleteWriter of tao
func DeleteWriter(configKey string) error {
	taoLogger.mu.Lock()
	defer taoLogger.mu.Unlock()

	writer, ok := taoLogger.writers[configKey]
	if !ok {
		return NewError(ParamInvalid, "log: %s's writer not set", configKey)
	}
	delete(taoLogger.writers, configKey)

	// writer close
	if l, ok := writer.(io.Closer); ok {
		return l.Close()
	}
	return nil
}

// GetLogger in tao
func GetLogger(configKey string) Logger {
	return taoLogger.loggers[configKey]
}

// SetLogger to tao
func SetLogger(configKey string, logger Logger) error {
	taoLogger.mu.Lock()
	defer taoLogger.mu.Unlock()

	if taoLogger.loggers == nil {
		taoLogger.loggers = make(map[string]Logger)
	}

	if _, ok := taoLogger.loggers[configKey]; ok {
		return NewError(DuplicateCall, "log: %s's logger has been set before", configKey)
	}

	taoLogger.loggers[configKey] = logger
	return nil
}

// DeleteLogger of tao
func DeleteLogger(configKey string) error {
	taoLogger.mu.Lock()
	defer taoLogger.mu.Unlock()

	logger, ok := taoLogger.loggers[configKey]
	if !ok {
		return NewError(ParamInvalid, "log: %s's logger not set", configKey)
	}
	delete(taoLogger.loggers, configKey)

	// logger close
	if l, ok := logger.(io.Closer); ok {
		return l.Close()
	}
	return nil
}

// Debug function wrap of TaoLogger
func Debug(v ...interface{}) {
	for _, l := range taoLogger.loggers {
		l.Debug(3, v...)
	}
}

// Debugf function wrap of TaoLogger
func Debugf(format string, v ...interface{}) {
	for _, l := range taoLogger.loggers {
		l.Debugf(3, format, v...)
	}
}

// Info function wrap of TaoLogger
func Info(v ...interface{}) {
	for _, l := range taoLogger.loggers {
		l.Info(3, v...)
	}
}

// Infof function wrap of TaoLogger
func Infof(format string, v ...interface{}) {
	for _, l := range taoLogger.loggers {
		l.Infof(3, format, v...)
	}
}

// Warn function wrap of TaoLogger
func Warn(v ...interface{}) {
	for _, l := range taoLogger.loggers {
		l.Warn(3, v...)
	}
}

// Warnf function wrap of TaoLogger
func Warnf(format string, v ...interface{}) {
	for _, l := range taoLogger.loggers {
		l.Warnf(3, format, v...)
	}
}

// Error function wrap of TaoLogger
func Error(v ...interface{}) {
	for _, l := range taoLogger.loggers {
		l.Error(3, v...)
	}
}

// Errorf function wrap of TaoLogger
func Errorf(format string, v ...interface{}) {
	for _, l := range taoLogger.loggers {
		l.Errorf(3, format, v...)
	}
}

// Panic function wrap of TaoLogger
func Panic(v ...interface{}) {
	for _, l := range taoLogger.loggers {
		l.Panic(3, v...)
	}
}

// Panicf function wrap of TaoLogger
func Panicf(format string, v ...interface{}) {
	for _, l := range taoLogger.loggers {
		l.Panicf(3, format, v...)
	}
}

// Fatal function wrap of TaoLogger
func Fatal(v ...interface{}) {
	for _, l := range taoLogger.loggers {
		l.Fatal(3, v...)
	}
}

// Fatalf function wrap of TaoLogger
func Fatalf(format string, v ...interface{}) {
	for _, l := range taoLogger.loggers {
		l.Fatalf(3, format, v...)
	}
}
