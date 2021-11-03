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
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	taskHello = NewTask("hello", func(ctx context.Context, param Parameter) (Parameter, error) {
		param.Set("message", "hello run")
		return param.Clone(), nil
	})
	taskError = NewTask("error", func(ctx context.Context, param Parameter) (Parameter, error) {
		param.Set("message", "error run")
		return param, NewError("I", "error run")
	}, SetClose(func() error {
		return NewError("II", "error close")
	}))
)

func TestNewTask(t *testing.T) {
	t.Run("TestTaskRun_Run", func(t *testing.T) {
		assert.Equal(t, nil, taskHello.Run(context.Background(), NewParameter()))
		assert.NotEqual(t, nil, taskHello.Run(context.Background(), NewParameter()))

		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		assert.Equal(t, ContextCanceled, taskError.Run(ctx, NewParameter()).(Error).Code())

		assert.Equal(t, "I", taskError.Run(context.Background(), NewParameter()).(Error).Code())
	})

	t.Run("TestTaskRun_GetName", func(t *testing.T) {
		assert.Equal(t, "hello", taskHello.Name())
		assert.Equal(t, "error", taskError.Name())
	})

	t.Run("TestTaskRun_Result", func(t *testing.T) {
		assert.Equal(t, "hello run", taskHello.Result().Get("message"))
		assert.Equal(t, "error run", taskError.Result().Get("message"))
	})

	t.Run("TestTaskRun_Error", func(t *testing.T) {
		assert.Equal(t, "", taskHello.Error())
		assert.NotEqual(t, "", taskError.Error())
	})

	t.Run("TestTaskRun_Close", func(t *testing.T) {
		assert.Equal(t, nil, taskHello.Close())
		assert.Equal(t, "II", taskError.Close().(Error).Code())
		assert.Equal(t, TaskCloseTwice, taskError.Close().(Error).Code())
	})
}
