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
	"fmt"
	"log"
)

// Config interface
type Config interface {
	// Name of Config
	Name() string
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

// LoadConfig by key of config
func LoadConfig(configKey string, config Config) error {
	c, ok := configInterfaceMap[configKey]
	if !ok {
		return NewError(ConfigNotFound, "config: %s not found", configKey)
	}
	bytes, err := json.Marshal(c)
	if err != nil {
		return NewErrorWrapped("config: fail to marshal", err)
	}
	err = json.Unmarshal(bytes, &config)
	if err != nil {
		return NewErrorWrapped(fmt.Sprintf("config: fail to unmarshal config bytes for '%q'", configKey), err)
	}
	return nil
}

// SetConfig by key & Config
func SetConfig(configKey string, config Config) error {
	_, ok := configMap[configKey]
	if ok {
		return NewError(DuplicateCall, "config: %s has been set before", configKey)
	}
	configMap[configKey] = config
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

// Name of Config
func (t *taoConfig) Name() string {
	return ConfigKey
}

// ValidSelf with some default values
func (t *taoConfig) ValidSelf() {
	if t.Log == nil {
		t.Log = defaultTao.Log
	} else {
		if t.Log.Level < DEBUG || t.Log.Level > FATAL {
			t.Log.Level = defaultTao.Log.Level
		}
		if t.Log.CallDepth <= 0 {
			t.Log.CallDepth = defaultTao.Log.CallDepth
		}
		if t.Log.Type == 0 {
			t.Log.Type = defaultTao.Log.Type
		}
		if t.Log.Type&File != 0 {
			if t.Log.Path == "" {
				t.Log.Path = defaultTao.Log.Path
			}
		}
		if t.Log.Flag == 0 {
			t.Log.Flag = defaultTao.Log.Flag
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
