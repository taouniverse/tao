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
	"sync"
)

// Parameter describe function input or output
type Parameter interface {
	Get(key string) (val interface{})
	Set(key string, val interface{})
	Clone() Parameter
	Delete(key string)
	String() string
}

var _ Parameter = (*param)(nil)

// param store params
// implements Parameter
type param struct {
	mu sync.RWMutex

	params map[string]interface{}
}

// NewParameter constructor of Parameter
func NewParameter() Parameter {
	return &param{
		params: make(map[string]interface{}),
	}
}

// Get value with key
func (p *param) Get(key string) (value interface{}) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.params[key]
}

// Set value with key
func (p *param) Set(key string, value interface{}) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.params[key] = value
}

// Delete value with key
func (p *param) Delete(key string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	delete(p.params, key)
}

// Clone param
func (p *param) Clone() Parameter {
	p.mu.RLock()
	defer p.mu.RUnlock()
	m := make(map[string]interface{}, len(p.params))
	for k, v := range p.params {
		m[k] = v
	}
	return &param{
		params: m,
	}
}

// String of param
func (p *param) String() string {
	marshal, err := json.Marshal(p)
	if err != nil {
		return ""
	}
	return string(marshal)
}

// MarshalJSON to marshal param
func (p *param) MarshalJSON() ([]byte, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return json.Marshal(p.params)
}
