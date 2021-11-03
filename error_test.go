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
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	err1 = NewErrorUnWrapper("err1", nil)
	err2 = NewErrorUnWrapper("err2", err1)
	err3 = NewErrorUnWrapper("err3", err2)

	err4 = NewErrorUnWrapper("err4", nil)
	err5 = NewErrorUnWrapper("err5", nil)

	err = NewError("I", "Error I")
)

func TestNewErrorUnWrapper(t *testing.T) {
	t.Run("TestErrorUnWrap_Error", func(t *testing.T) {
		assert.Equal(t, "err1", err1.Error())
		assert.Equal(t, "err2"+errSplit+err1.Error(), err2.Error())
		assert.Equal(t, "err3"+errSplit+err2.Error(), err3.Error())
	})

	t.Run("TestErrorUnWrap_UnWrap", func(t *testing.T) {
		assert.Equal(t, nil, err1.Unwrap())
		assert.Equal(t, err1, err2.Unwrap())
		assert.Equal(t, err2, err3.Unwrap())
	})
}

func TestNewError(t *testing.T) {
	t.Run("TestError_Code", func(t *testing.T) {
		assert.Equal(t, "I", err.Code())
	})

	t.Run("TestError_Error", func(t *testing.T) {
		err.Wrap(err4)
		assert.Equal(t, "Error I"+errSplit+err4.Error(), err.Error())
	})

	t.Run("TestError_Wrap", func(t *testing.T) {
		err.Wrap(nil)
		err.Wrap(err5)
		assert.Equal(t, "Error I"+errSplit+err5.Error()+errSplit+err4.Error(), err.Error())
	})

	t.Run("TestError_Cause", func(t *testing.T) {
		assert.Equal(t, NewErrorUnWrapper(err5.Error(), err4), err.Cause())
	})
}
