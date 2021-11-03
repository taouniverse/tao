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
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	taskNameEmpty = NewTask("", func(ctx context.Context, param Parameter) (Parameter, error) {
		return nil, nil
	}, SetClose(func() error {
		return nil
	}))
	taskA = NewTask("A", func(ctx context.Context, param Parameter) (Parameter, error) {
		param.Set("A", "hello from A")
		return param, nil
	}, SetClose(func() error {
		return errors.New("error close A")
	}))
	taskB = NewTask("B", func(ctx context.Context, param Parameter) (Parameter, error) {
		param.Set("B", "hello from B")
		return param, errors.New("error B")
	}, SetClose(func() error {
		return errors.New("error close B")
	}))
	taskC = NewTask("C", func(ctx context.Context, param Parameter) (Parameter, error) {
		param.Set("C", "hello from C")
		return param, errors.New("error C")
	}, SetClose(func() error {
		return errors.New("error close C")
	}))
	// pipeline in pipeline
	pipeA = NewPipeline("pA")
	pipeB = NewPipeline("pB")
	pipe  = NewPipeline("hello")
)

func TestNewPipeline(t *testing.T) {
	t.Run("TestPipelineRun_Register", func(t *testing.T) {
		assert.NotEqual(t, nil, pipeA.Register(NewPipeTask(taskNameEmpty)))
		assert.Equal(t, nil, pipeA.Register(NewPipeTask(taskA)))
		assert.NotEqual(t, nil, pipeA.Register(NewPipeTask(taskA)))
		assert.Equal(t, nil, pipeB.Register(NewPipeTask(taskB)))
		assert.Equal(t, nil, pipeB.Register(NewPipeTask(taskC, "B")))
		assert.Equal(t, nil, pipe.Register(NewPipeTask(pipeA)))
		assert.Equal(t, nil, pipe.Register(NewPipeTask(pipeB, "pA")))
	})

	t.Run("TestPipelineRun_Run", func(t *testing.T) {
		assert.Equal(t, "", pipe.Error())

		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		assert.Equal(t, ContextCanceled, pipeA.Run(ctx, NewParameter()).(Error).Code())

		input := NewParameter()
		input.Set("tao", "useful")
		assert.NotEqual(t, nil, pipe.Run(context.Background(), input))
	})

	t.Run("TestPipelineRun_Run_Twice", func(t *testing.T) {
		assert.Equal(t, TaskRunTwice, pipe.Run(context.Background(), NewParameter()).(Error).Code())
	})

	t.Run("TestPipelineRun_GetName", func(t *testing.T) {
		assert.Equal(t, "hello", pipe.Name())
	})

	t.Run("TestPipelineRun_Result", func(t *testing.T) {
		assert.NotEmpty(t, pipe.Result().Get("pA"))
		assert.NotEmpty(t, pipe.Result().Get("pB"))
		assert.Equal(t, pipe.Result().Get("pA"), pipeA.Result())
		assert.Equal(t, pipe.Result().Get("pB"), pipeB.Result())

		t.Log(pipeA.Result().String())
		t.Log(pipeB.Result().String())

		assert.NotEmpty(t, pipeA.Result().Get("A"))
		assert.NotEmpty(t, pipeB.Result().Get("B"))
		assert.NotEmpty(t, pipeB.Result().Get("C"))
		assert.Equal(t, pipeA.Result().Get("A"), taskA.Result())
		assert.Equal(t, pipeB.Result().Get("B"), taskB.Result())
		assert.Equal(t, pipeB.Result().Get("C"), taskC.Result())

		t.Log(taskA.Result().String())
		t.Log(taskB.Result().String())
		t.Log(taskC.Result().String())
	})

	t.Run("TestPipelineRun_Error", func(t *testing.T) {
		assert.NotEqual(t, "", pipe.Error())
		assert.Equal(t, pipe.Error(), pipeB.Error())
		assert.Equal(t, "", pipeA.Error())
		t.Log(pipe.Error())

		t.Log(taskB.Error())
		t.Log(taskC.Error())
	})

	t.Run("TestPipelineRun_Close", func(t *testing.T) {
		assert.NotEqual(t, nil, pipe.Close())
		assert.Equal(t, TaskCloseTwice, pipe.Close().(Error).Code())
	})
}
