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
	"context"
	"sync"
)

// TaskState to describe state of task
type TaskState uint8

const (
	// Runnable task
	Runnable TaskState = iota
	// Running task
	Running
	// Over task
	Over
	// Closed task
	Closed
)

// TaskRun with param
type TaskRun func(ctx context.Context, param Parameter) (Parameter, error)

// Task single Task
type Task interface {
	Name() string
	Run(ctx context.Context, param Parameter) error
	Result() Parameter
	Error() string
	Close() error
	State() TaskState
}

var _ Task = (*task)(nil)

// task implement of Task
type task struct {
	mu sync.RWMutex

	name string

	fun       TaskRun
	closeFun  func() error
	postStart TaskRun
	preStop   TaskRun

	result Parameter
	err    error
	state  TaskState
}

// NewTask constructor of Task
func NewTask(name string, fun TaskRun, options ...TaskOption) Task {
	if fun == nil {
		return nil
	}

	t := &task{
		name:  name,
		fun:   fun,
		state: Runnable,
	}

	for _, option := range options {
		option(t)
	}

	return t
}

// Name of Task
func (t *task) Name() string {
	return t.name
}

// Run Task
func (t *task) Run(ctx context.Context, param Parameter) (err error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if ctx == nil {
		ctx = context.Background()
	}

	if param == nil {
		param = NewParameter()
	}

	if t.state == Closed {
		return NewError(TaskClosed, "task: task has been closed")
	}

	if t.state != Runnable {
		return NewError(TaskRunTwice, "task: Run called twice for task "+t.name)
	}

	select {
	case <-ctx.Done():
		return NewError(ContextCanceled, "task: context has been canceled")
	default:
	}

	t.state = Running
	defer func() {
		// SPECIAL: result should be cloned param because it's just for this task
		t.result = param.Clone()
		t.err = err
		t.state = Over
	}()

	if t.postStart != nil {
		param, err = t.postStart(ctx, param)
		if err != nil {
			return
		}
	}

	param, err = t.fun(ctx, param)
	if err != nil {
		return
	}

	if t.preStop != nil {
		param, err = t.preStop(ctx, param)
	}
	return
}

// Result of Task
func (t *task) Result() Parameter {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.result
}

// Error info of Task
func (t *task) Error() string {
	t.mu.RLock()
	defer t.mu.RUnlock()
	if t.err == nil {
		return ""
	}
	return t.err.Error()
}

// Close resource of Task
func (t *task) Close() error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.state == Running {
		return NewError(TaskRunning, "task: task %s is running", t.Name())
	}

	if t.state == Closed {
		return NewError(TaskCloseTwice, "task: Close called twice for task %s", t.Name())
	}

	t.state = Closed

	if t.closeFun != nil {
		return t.closeFun()
	}
	return nil
}

// State of task
func (t *task) State() TaskState {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.state
}

// TaskOption optional function of task
type TaskOption func(t *task)

// SetClose of task
func SetClose(cf func() error) TaskOption {
	return func(t *task) {
		t.closeFun = cf
	}
}

// SetPostStart of task
func SetPostStart(tr TaskRun) TaskOption {
	return func(t *task) {
		t.postStart = tr
	}
}

// SetPreStop of task
func SetPreStop(tr TaskRun) TaskOption {
	return func(t *task) {
		t.preStop = tr
	}
}
