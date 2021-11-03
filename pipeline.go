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
	"context"
	"sync"
)

// pipeTask task in Pipeline
type pipeTask struct {
	Task
	runAfter []string
}

// NewPipeTask constructor of pipeTask
func NewPipeTask(task Task, runAfter ...string) *pipeTask {
	return &pipeTask{
		Task:     task,
		runAfter: runAfter,
	}
}

// Pipeline to run tasks in order
// pipeline is also a task
type Pipeline interface {
	Task
	Register(task *pipeTask) error
}

var _ Pipeline = (*pipeline)(nil)

// pipeline implement of Pipeline
type pipeline struct {
	wg sync.WaitGroup
	mu sync.RWMutex

	name string

	tasks     []*pipeTask
	signals   map[string]chan struct{}
	closeChan chan func() error

	results Parameter
	err     Error
	state   TaskState
}

// NewPipeline constructor of Pipeline
func NewPipeline(name string) Pipeline {
	return &pipeline{
		name:    name,
		tasks:   make([]*pipeTask, 0),
		signals: make(map[string]chan struct{}),
		err:     nil,
		state:   Runnable,
	}
}

// Name of Pipeline
func (p *pipeline) Name() string {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.name
}

// Register task to Pipeline
func (p *pipeline) Register(task *pipeTask) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	tName := task.Name()
	if tName == "" {
		return NewError(ParamInvalid, "pipeline: Register task name is empty")
	}
	if _, dup := p.signals[tName]; dup {
		return NewError(ParamInvalid, "pipeline: Register called twice for task "+tName)
	}

	p.tasks = append(p.tasks, task)
	p.signals[tName] = make(chan struct{}, 1)
	return nil
}

// Run Pipeline
func (p *pipeline) Run(ctx context.Context, param Parameter) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.state == Close {
		return NewError(TaskClosed, "pipeline: pipeline has been closed")
	}

	if p.state != Runnable {
		return NewError(TaskRunTwice, "pipeline: Run called twice for task "+p.name)
	}

	select {
	case <-ctx.Done():
		return NewError(ContextCanceled, "pipeline: context has been canceled")
	default:
	}

	p.state = Running
	// init closeChan & results when run
	p.closeChan = make(chan func() error, len(p.tasks))
	p.results = NewParameter()

	for _, task := range p.tasks {
		p.wg.Add(1)
		go p.taskRun(ctx, task, param)
	}
	p.wg.Wait()

	p.state = Over
	return p.err
}

func (p *pipeline) taskRun(ctx context.Context, task *pipeTask, param Parameter) {
	defer p.wg.Done()
	var err error

	// waiting...
	for _, pre := range task.runAfter {
		if signal, ok := p.signals[pre]; ok {
			<-signal
		}
	}

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
	// signal
	close(p.signals[task.Name()])
	// close fun
	p.closeChan <- task.Close
}

// String context of Pipeline
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
func (p *pipeline) Close() error {
	p.mu.Lock()
	defer p.mu.Unlock()
	var (
		err        error
		closeSlice = make([]func() error, 0, len(p.tasks))
	)

	if p.state == Running {
		return NewError(TaskRunning, "pipeline: pipeline is running")
	}

	if p.state == Close {
		return NewError(TaskCloseTwice, "pipeline: Close called twice for pipeline "+p.name)
	}

	// close chan before for range
	close(p.closeChan)
	for closeFn := range p.closeChan {
		closeSlice = append(closeSlice, closeFn)
	}

	for i := len(p.tasks) - 1; i >= 0; i-- {
		if e := closeSlice[i](); e != nil {
			err = NewErrorUnWrapper(e.Error(), err)
		}
	}

	p.state = Close
	return err
}
