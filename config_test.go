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
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// printConfig implements Config
type printConfig struct {
	Print     string   `json:"print"`
	RunAfter_ []string `json:"run_after"`
}

func (l *printConfig) ToTask() Task {
	return NewTask("print", func(ctx context.Context, param Parameter) (Parameter, error) {
		select {
		case <-ctx.Done():
			return param, NewError(ContextCanceled, "test: ctx already Done")
		default:
			fmt.Println(l.Print)
			return param, nil
		}
	})
}

func (l *printConfig) ValidSelf() {
}

func (l *printConfig) RunAfter() []string {
	return l.RunAfter_
}

func TestJsonConfig(t *testing.T) {
	file := `
{
    "print": {
        "print": "==============  hello,tao!  ==============",
        "run_after": []
    }
}`
	err := json.Unmarshal([]byte(file), &configInterfaceMap)
	assert.Nil(t, err)
	bytes, err := GetConfigBytes("print")
	assert.Nil(t, err)
	c := new(printConfig)
	err = json.Unmarshal(bytes, &c)
	assert.Nil(t, err)
	err = SetConfig("print", c)
	assert.Nil(t, err)
	err = SetConfig("print", c)
	assert.NotNil(t, err)
	t.Log(configMap["print"])

	_, err = GetConfigBytes("Unknown")
	assert.NotNil(t, err)
}
