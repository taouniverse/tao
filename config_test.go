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
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJsonConfig(t *testing.T) {
	err := LoadConfig("Unknown", nil)
	assert.NotNil(t, err)
	assert.Equal(t, ConfigNotFound, err.(ErrorTao).Code())
}

const printConfigKey = "print"

// printConfig implements Config
type printConfig struct {
	Print     string   `json:"print"`
	Times     int      `json:"times"`
	RunAfters []string `json:"run_after"`
}

var defaultPrint = &printConfig{
	Print: "==============  hello,tao!  ==============",
	Times: 1,
}

// Name of Config
func (l *printConfig) Name() string {
	return printConfigKey
}

// ValidSelf with some default values
func (l *printConfig) ValidSelf() {
	if l.Print == "" {
		l.Print = defaultPrint.Print
	}
	if l.Times == 0 {
		l.Times = defaultPrint.Times
	}
	if l.RunAfters == nil {
		l.RunAfters = defaultPrint.RunAfters
	}
}

// ToTask transform itself to Task
func (l *printConfig) ToTask() Task {
	return NewTask("print", func(ctx context.Context, param Parameter) (Parameter, error) {
		select {
		case <-ctx.Done():
			return param, NewError(ContextCanceled, "test: ctx already Done")
		default:
			for i := 0; i < l.Times; i++ {
				fmt.Println(l.Print)
			}
			return param, nil
		}
	})
}

// RunAfter defines pre task names
func (l *printConfig) RunAfter() []string {
	return l.RunAfters
}
