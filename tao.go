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

// taoConfig implements Config
type taoConfig struct {
	*Log       `json:"log"`
	HideBanner bool `json:"hide_banner"`
}

var defaultTao = &taoConfig{
	Log: &Log{
		Level:     DEBUG,
		Type:      Console | File,
		CallDepth: 3,
		Path:      "./test.log",
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
		if t.Level < DEBUG || t.Level > FATAL {
			t.Level = defaultTao.Level
		}
		if t.Type == 0 {
			t.Type = defaultTao.Type
		}
		if t.CallDepth <= 0 {
			t.CallDepth = defaultTao.CallDepth
		}
		if t.Type&File != 0 {
			if t.Path == "" {
				t.Path = defaultTao.Path
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
