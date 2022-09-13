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
	"encoding/json"
	"log"
)

// Config interface
type Config interface {
	// Default config
	Default() Config
	// ValidSelf with some default values
	ValidSelf()
	// ToTask transform itself to Task
	ToTask() Task
	// RunAfter defines pre task names
	RunAfter() []string
}

// init config file to this interface map
var configInterfaceMap = make(map[string]interface{})

// transform interface to concrete Config type
var configMap = make(map[string]Config)

// GetConfigBytes in json schema by key of config
func GetConfigBytes(key string) ([]byte, error) {
	c, ok := configInterfaceMap[key]
	if !ok {
		return nil, NewError(ConfigNotFound, "config: %s not found", key)
	}
	bytes, err := json.Marshal(c)
	if err != nil {
		return nil, NewErrorWrapped("config: marshal failed", err)
	}
	return bytes, nil
}

// SetConfig by key & Config
func SetConfig(key string, c Config) error {
	_, ok := configMap[key]
	if ok {
		return NewError(DuplicateCall, "config: %s has been set before", key)
	}
	configMap[key] = c
	return nil
}

// ConfigKey for this repo
const ConfigKey = "tao"

// taoConfig implements Config
type taoConfig struct {
	Log    *Log    `json:"log"`
	Banner *Banner `json:"banner"`
}

// Banner config
type Banner struct {
	Hide    bool   `json:"hide"`
	Content string `json:"content"`
}

var defaultTao = &taoConfig{
	Log: &Log{
		Level:     DEBUG,
		Type:      Console | File,
		Flag:      log.LstdFlags | log.Lshortfile,
		CallDepth: 3,
		Path:      "./test.log",
		Disable:   false,
	},
	Banner: &Banner{
		Hide: false,
		Content: `
___________              
\__    ___/____    ____  
  |    |  \__  \  /  _ \ 
  |    |   / __ \(  <_> )
  |____|  (____  /\____/ 
               \/
`,
	},
}

// Default config
func (t *taoConfig) Default() Config {
	return defaultTao
}

// ValidSelf with some default values
func (t *taoConfig) ValidSelf() {
	if t.Log == nil {
		t.Log = defaultTao.Log
	} else {
		if t.Log.Level < DEBUG || t.Log.Level > FATAL {
			t.Log.Level = defaultTao.Log.Level
		}
		if t.Log.Type == 0 {
			t.Log.Type = defaultTao.Log.Type
		}
		if t.Log.CallDepth <= 0 {
			t.Log.CallDepth = defaultTao.Log.CallDepth
		}
		if t.Log.Type&File != 0 {
			if t.Log.Path == "" {
				t.Log.Path = defaultTao.Log.Path
			}
		}
	}
	if t.Banner == nil {
		t.Banner = defaultTao.Banner
	} else {
		if !t.Banner.Hide && t.Banner.Content == "" {
			t.Banner.Content = defaultTao.Banner.Content
		}
	}
}

// ToTask transform itself to Task
func (t *taoConfig) ToTask() Task {
	return nil
}

// RunAfter defines pre task names
func (t *taoConfig) RunAfter() []string {
	return nil
}
