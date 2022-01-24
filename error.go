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
	"fmt"
	"sync"
)

const (
	errSplit = "\n==> "
)

// ErrorUnWrapper extension of error which can be unwrapped
type ErrorUnWrapper interface {
	error
	Unwrap() error
}

var _ ErrorUnWrapper = (*errorWrapped)(nil)

// errorWrapped implements ErrorUnWrapper
type errorWrapped struct {
	msg string
	err error
}

// NewErrorWrapped constructor of errorWrapped
func NewErrorWrapped(format string, e error) ErrorUnWrapper {
	if e != nil {
		return &errorWrapped{format + errSplit + e.Error(), e}
	}
	return &errorWrapped{format, nil}
}

// Error string
func (e *errorWrapped) Error() string { return e.msg }

// Unwrap e self
func (e *errorWrapped) Unwrap() error { return e.err }

// Is a same error
func (e *errorWrapped) Is(err error) bool { return e.Error() == err.Error() }

// ErrorTao extension of error, wrap of error
type ErrorTao interface {
	error
	Code() string
	Wrap(err error)
	Cause() error
}

var _ ErrorTao = (*errorTao)(nil)
var _ ErrorUnWrapper = (*errorTao)(nil)

// errorTao with code & message
// code for computer
// message for user
// implements ErrorTao & ErrorUnWrapper
type errorTao struct {
	mutex sync.RWMutex

	code    string
	message string

	cause ErrorUnWrapper
}

// NewError constructor of Error
func NewError(code, message string, a ...interface{}) ErrorTao {
	return &errorTao{
		code:    code,
		message: fmt.Sprintf(message, a...),
	}
}

// Code string
func (e *errorTao) Code() string {
	e.mutex.RLock()
	defer e.mutex.RUnlock()
	return e.code
}

// Error string
func (e *errorTao) Error() string {
	e.mutex.RLock()
	defer e.mutex.RUnlock()
	if e.cause != nil {
		return e.message + errSplit + e.cause.Error()
	}
	return e.message
}

// Wrap error into errorTao
func (e *errorTao) Wrap(err error) {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	if err == nil {
		return
	}
	// wrap cause & err
	if e.cause != nil {
		e.cause = NewErrorWrapped(err.Error(), e.cause)
		return
	}
	// init cause
	if cause, ok := err.(ErrorUnWrapper); ok {
		e.cause = cause
	} else {
		e.cause = NewErrorWrapped(err.Error(), nil)
	}
}

// Cause of error
func (e *errorTao) Cause() error {
	e.mutex.RLock()
	defer e.mutex.RUnlock()
	return e.cause
}

// Unwrap e self
func (e *errorTao) Unwrap() error {
	e.mutex.RLock()
	defer e.mutex.RUnlock()
	return e.cause
}

/**
ErrorCode
*/
const (
	Unknown         = "Unknown"
	ParamInvalid    = "ParamInvalid"
	ContextCanceled = "ContextCanceled"
	DuplicateCall   = "DuplicateCall"
	TaskRunTwice    = "TaskRunTwice"
	TaskCloseTwice  = "TaskCloseTwice"
	TaskClosed      = "TaskClosed"
	TaskRunning     = "TaskRunning"
	ConfigNotFound  = "ConfigNotFound"
)
