// Copyright 2022 huija
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
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRegister(t *testing.T) {
	p := new(printConfig)
	err := Register(printConfigKey, p, nil)
	assert.Nil(t, err)

	err = Register(printConfigKey, nil, nil)
	assert.NotNil(t, err)

	err = SetConfig(printConfigKey, nil)
	assert.NotNil(t, err)

	err = universeInit()
	assert.NotNil(t, err)
}

func TestRun(t *testing.T) {
	t.Log(new(taoConfig).ToTask())
	t.Log(new(taoConfig).RunAfter())

	cancel, cancelFunc := context.WithCancel(context.Background())
	cancelFunc()
	err := Run(cancel, nil)
	assert.NotNil(t, err)

	Add(1)
	Done()
	err = Run(nil, nil)
	assert.Nil(t, err)

	err = Run(nil, nil)
	assert.NotNil(t, err)
}
