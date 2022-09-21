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

// PipeTask task in Pipeline
type PipeTask struct {
	Task
	runAfter []string
}

// NewPipeTask constructor of PipeTask
func NewPipeTask(task Task, runAfter ...string) *PipeTask {
	return &PipeTask{
		Task:     task,
		runAfter: runAfter,
	}
}

// Pipeline to run tasks in order
// pipeline is also a task
type Pipeline interface {
	Task
	Register(task *PipeTask) error
}

var _ Pipeline = (*pipeline)(nil)

// pipeline implement of Pipeline
type pipeline struct {
	wg sync.WaitGroup
	mu sync.RWMutex

	name string

	tasks     []*PipeTask
	signals   map[string]chan struct{}
	closeChan chan func() error
	postStart *PipeTask
	preStop   *PipeTask

	results Parameter
	err     ErrorTao
	state   TaskState
}

// NewPipeline constructor of Pipeline
func NewPipeline(name string, options ...PipelineOption) Pipeline {
	p := &pipeline{
		name:    name,
		tasks:   make([]*PipeTask, 0),
		signals: make(map[string]chan struct{}),
		err:     nil,
		state:   Runnable,
	}

	for _, option := range options {
		option(p)
	}

	return p
}

// Name of Pipeline
func (p *pipeline) Name() string {
	return p.name
}

// Register task to Pipeline
func (p *pipeline) Register(task *PipeTask) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if task == nil {
		return NewError(ParamInvalid, "pipeline: Register task is null")
	}

	// nothing to do
	if task.Task == nil {
		return nil
	}

	tName := task.Name()
	if tName == "" {
		return NewError(ParamInvalid, "pipeline: Register task name is empty")
	}

	if p.state == Running {
		return NewError(TaskRunning, "pipeline: pipeline is running")
	}
	if _, dup := p.signals[tName]; dup {
		return NewError(ParamInvalid, "pipeline: Register called twice for task %q", tName)
	}

	p.tasks = append(p.tasks, task)
	p.signals[tName] = make(chan struct{}, 1)
	return nil
}

// Run Pipeline
func (p *pipeline) Run(ctx context.Context, param Parameter) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if ctx == nil {
		ctx = context.Background()
	}

	if param == nil {
		param = NewParameter()
	}

	if p.state == Closed {
		return NewError(TaskClosed, "pipeline: pipeline has been closed")
	}

	if p.state != Runnable {
		return NewError(TaskRunTwice, "pipeline: Run called twice for pipeline "+p.name)
	}

	select {
	case <-ctx.Done():
		return NewError(ContextCanceled, "pipeline: context has been canceled")
	default:
	}

	p.state = Running
	defer func() {
		p.state = Over
	}()

	// init closeChan & results when run
	p.closeChan = make(chan func() error, len(p.tasks)+2)
	p.results = NewParameter()

	if p.postStart != nil {
		p.taskRun(ctx, p.postStart, param, false)
		if p.err != nil {
			return p.err
		}
	}

	for _, task := range p.tasks {
		p.wg.Add(1)
		go p.taskRun(ctx, task, param, true)
	}
	p.wg.Wait()

	if p.preStop != nil {
		p.taskRun(ctx, p.preStop, param, false)
	}

	return p.err
}

func (p *pipeline) taskRun(ctx context.Context, task *PipeTask, param Parameter, async bool) {
	if async {
		defer p.wg.Done()
	}
	var err error

	// waiting...
	for _, pre := range task.runAfter {
		if signal, ok := p.signals[pre]; ok {
			<-signal
		}
	}

	// register close before run
	p.closeChan <- task.Close

	// run & wrap cause
	err = task.Run(ctx, param)
	if err != nil {
		if p.err == nil {
			p.err = NewError(Unknown, err.Error())
		} else {
			p.err.Wrap(err)
		}
	}

	// result
	p.results.Set(task.Name(), task.Result())
	if async {
		// signal
		close(p.signals[task.Name()])
	}
}

// Result of Pipeline
func (p *pipeline) Result() Parameter {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.results
}

// Error info of Pipeline
func (p *pipeline) Error() string {
	p.mu.RLock()
	defer p.mu.RUnlock()
	if p.err == nil {
		return ""
	}
	return p.err.Error()
}

// Close resource of Pipeline
func (p *pipeline) Close() (err error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	var (
		closeSlice = make([]func() error, 0, len(p.tasks)+2)
	)

	if p.state == Running {
		return NewError(TaskRunning, "pipeline: pipeline %s is running", p.Name())
	}

	if p.state == Closed {
		return NewError(TaskCloseTwice, "pipeline: Close called twice for pipeline %s", p.Name())
	}

	// close chan before for range
	close(p.closeChan)
	if p.postStart != nil {
		closeSlice = append(closeSlice, p.postStart.Close)
	}

	for closeFn := range p.closeChan {
		closeSlice = append(closeSlice, closeFn)
	}

	if p.preStop != nil {
		closeSlice = append(closeSlice, p.preStop.Close)
	}

	// close functions are called in reverse order
	for i := len(closeSlice) - 1; i >= 0; i-- {
		if e := closeSlice[i](); e != nil {
			err = NewErrorWrapped(e.Error(), err)
		}
	}

	p.state = Closed
	return
}

// State of pipeline
func (p *pipeline) State() TaskState {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.state
}

// PipelineOption optional function of pipeline
type PipelineOption func(p *pipeline)

// SetPostStartTask of pipeline
func SetPostStartTask(t *PipeTask) PipelineOption {
	return func(p *pipeline) {
		p.postStart = t
	}
}

// SetPreStopTask of pipeline
func SetPreStopTask(t *PipeTask) PipelineOption {
	return func(p *pipeline) {
		p.preStop = t
	}
}
