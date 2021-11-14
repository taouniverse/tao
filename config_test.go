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
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

type logConfig struct {
	Type      string   `json:"type"`
	Level     string   `json:"level"`
	Path      string   `json:"path"`
	RunAfter_ []string `json:"run_after"`
}

func (l *logConfig) ToTask() Task {
	return nil
}

func (l *logConfig) ValidSelf() {
}

func (l *logConfig) RunAfter() []string {
	return l.RunAfter_
}

func TestJsonConfig(t *testing.T) {
	file := `
{
    "log": {
        "type": "file",
        "level": "debug",
        "path": "./test.log",
        "run_after": []
    }
}`
	err := json.Unmarshal([]byte(file), &configInterfaceMap)
	assert.Nil(t, err)
	bytes, err := GetConfigBytes("log")
	assert.Nil(t, err)
	c := new(logConfig)
	err = json.Unmarshal(bytes, &c)
	assert.Nil(t, err)
	err = SetConfig("log", c)
	assert.Nil(t, err)
	t.Log(configMap["log"])
}
