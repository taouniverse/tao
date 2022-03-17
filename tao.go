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
	"encoding/json"
)

// The Tao produced One; One produced Two; Two produced Three; Three produced All things.
var tao Pipeline

// t global config of tao
var t *taoConfig

// Run tao
func Run(ctx context.Context, param Parameter) (err error) {
	if ctx == nil {
		ctx = context.Background()
	}

	if param == nil {
		param = NewParameter()
	}

	if tao == nil {
		// refer to defaultConfigs in init.go to get some help
		return NewError(UniverseNotInit, "none of %+v existed", defaultConfigs)
	}

	// non-block check
	select {
	case <-ctx.Done():
		return NewError(ContextCanceled, "tao: context has been canceled")
	default:
	}

	// tasks run
	for _, c := range configMap {
		c.ValidSelf()
		err = tao.Register(NewPipeTask(c.ToTask(), c.RunAfter()...))
		if err != nil {
			return err
		}
	}

	// debug print
	cm, err := json.MarshalIndent(configMap, "", "  ")
	if err != nil {
		return err
	}
	Debugf("config data: \n%s", string(cm))

	return tao.Run(ctx, param)
}

// ConfigKey for this repo
const ConfigKey = "tao"

// banner of tao
const banner = `
___________              
\__    ___/____    ____  
  |    |  \__  \  /  _ \ 
  |    |   / __ \(  <_> )
  |____|  (____  /\____/ 
               \/
`

// taoConfig implements Config
type taoConfig struct {
	Log        *Log `json:"log"`
	HideBanner bool `json:"hide_banner"`
}

var defaultTao = &taoConfig{
	Log: &Log{
		Level:     DEBUG,
		Type:      Console | File,
		CallDepth: 3,
		Path:      "./test.log",
		Disable:   false,
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
}

// ToTask transform itself to Task
func (t *taoConfig) ToTask() Task {
	return nil
}

// RunAfter defines pre task names
func (t *taoConfig) RunAfter() []string {
	return nil
}
