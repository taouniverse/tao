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
	_, err := GetConfigBytes("Unknown")
	assert.NotNil(t, err)
}

const printConfigKey = "print"

// printConfig implements Config
type printConfig struct {
	Print     string   `json:"print"`
	Times     int      `json:"times"`
	RunAfterT []string `json:"run_after"`
}

var defaultPrint = &printConfig{
	Print: "==============  hello,tao!  ==============",
	Times: 1,
}

// Default config
func (l *printConfig) Default() Config {
	return defaultPrint
}

// ValidSelf with some default values
func (l *printConfig) ValidSelf() {
	if l.Print == "" {
		l.Print = defaultPrint.Print
	}
	if l.Times == 0 {
		l.Times = defaultPrint.Times
	}
	if l.RunAfterT == nil {
		l.RunAfterT = defaultPrint.RunAfterT
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
	return l.RunAfterT
}
