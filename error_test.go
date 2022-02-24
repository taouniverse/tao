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
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	err1 = NewErrorWrapped("err1", nil)
	err2 = NewErrorWrapped("err2", err1)
	err3 = NewErrorWrapped("err3", err2)

	err4 = NewErrorWrapped("err4", nil)
	err5 = NewErrorWrapped("err5", nil)

	err = NewError(Unknown, "Error I")
)

func TestNewErrorWrapped(t *testing.T) {
	t.Run("TestErrorWrapped_Error", func(t *testing.T) {
		assert.Equal(t, "err1", err1.Error())
		assert.Equal(t, "err2"+errSplit+err1.Error(), err2.Error())
		assert.Equal(t, "err3"+errSplit+err2.Error(), err3.Error())
	})

	t.Run("TestErrorWrapped_UnWrap", func(t *testing.T) {
		assert.Equal(t, nil, err1.Unwrap())
		assert.Equal(t, err1, err2.Unwrap())
		assert.Equal(t, err2, err3.Unwrap())
	})

	t.Run("TestErrorWrapped_Is", func(t *testing.T) {
		assert.Equal(t, true, errors.Is(err2, err1))
		assert.Equal(t, true, errors.Is(err3, err2))
	})

	t.Run("TestErrorWrapped_As", func(t *testing.T) {
		var wrapped = new(errorWrapped)
		assert.Equal(t, true, errors.As(err3, &wrapped))
		assert.Equal(t, err3, wrapped)
	})
}

func TestNewError(t *testing.T) {
	t.Run("TestError_Code", func(t *testing.T) {
		assert.Equal(t, Unknown, err.Code())
	})

	t.Run("TestError_Error", func(t *testing.T) {
		err.Wrap(err4)
		assert.Equal(t, fmt.Sprintf("<%s>Error I%s", Unknown, errSplit+err4.Error()), err.Error())
	})

	t.Run("TestError_Wrap", func(t *testing.T) {
		err.Wrap(nil)
		err.Wrap(err5)
		assert.Equal(t, fmt.Sprintf("<%s>Error I%s", Unknown, errSplit+err5.Error()+errSplit+err4.Error()), err.Error())
	})

	t.Run("TestError_Cause", func(t *testing.T) {
		assert.Equal(t, NewErrorWrapped(err5.Error(), err4), err.Cause())
	})

	t.Run("TestError_Is", func(t *testing.T) {
		assert.Equal(t, true, errors.Is(err, err4))
		assert.Equal(t, true, errors.Is(err, NewErrorWrapped(err5.Error(), err4)))
	})

	t.Run("TestError_As", func(t *testing.T) {
		var wrapped = new(errorWrapped)
		assert.Equal(t, true, errors.As(err, &wrapped))
		assert.Equal(t, NewErrorWrapped(err5.Error(), err4), wrapped)
	})
}
