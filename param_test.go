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
	parameter = NewParameter()
)

func TestNewParameter(t *testing.T) {
	t.Run("TestParam_Get", func(t *testing.T) {
		assert.Equal(t, nil, parameter.Get("key"))
	})

	t.Run("TestParam_Set", func(t *testing.T) {
		parameter.Set("key", "value")
		assert.Equal(t, "value", parameter.Get("key"))
	})

	t.Run("TestParam_Delete", func(t *testing.T) {
		parameter.Delete("key")
		assert.Equal(t, nil, parameter.Get("key"))
		parameter.Set("key", "value")
		assert.Equal(t, "value", parameter.Get("key"))
	})

	t.Run("TestParam_String", func(t *testing.T) {
		param1 := NewParameter()
		param1.Set("p", parameter)
		t.Log(param1.String())
	})

	t.Run("TestParam_Clone", func(t *testing.T) {
		param2 := parameter.Clone()
		assert.Equal(t, param2, parameter)
		assert.NotSame(t, param2, parameter)
	})
}
