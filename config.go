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
