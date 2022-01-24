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
	"fmt"
	"log"
	"os"
)

// The Tao produced One; One produced Two; Two produced Three; Three produced All things.
var tao Pipeline

// T global config of tao
var T *TaoConfig

// Run tao
func Run(ctx context.Context, param Parameter) (err error) {
	if ctx == nil {
		ctx = context.Background()
	}

	if param == nil {
		param = NewParameter()
	}

	// fallback
	if tao == nil {
		// init default tao
		taoInit()
		// warning print
		Warnf("%s not existed\n", defaultYamlConfig)
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

// TaoConfig implements Config
type TaoConfig struct {
	*Log       `json:"log"`
	HideBanner bool `json:"hide_banner"`
}

var defaultTao = &TaoConfig{
	Log: &Log{
		Level: DEBUG,
		Type:  COMMAND,
		Path:  "./test.log",
	},
}

// Default config
func (t *TaoConfig) Default() Config {
	return defaultTao
}

// ValidSelf with some default values
func (t *TaoConfig) ValidSelf() {
	if t.Log == nil {
		t.Log = defaultTao.Log
	} else {
		if t.Level < DEBUG || t.Level > FATAL {
			t.Level = defaultTao.Level
		}
		if t.Type == "" {
			t.Type = defaultTao.Type
		}
		if t.Type == File {
			if t.Path == "" {
				t.Path = defaultTao.Path
			}
		}
	}
}

// ToTask transform itself to Task
func (t *TaoConfig) ToTask() Task {
	return nil
}

// RunAfter defines pre task names
func (t *TaoConfig) RunAfter() []string {
	return nil
}

// ConfigKey for this repo
const ConfigKey = "tao"

func taoInit() {
	// transfer config bytes to object
	T = new(TaoConfig)
	bytes, err := GetConfigBytes(ConfigKey)
	if err != nil {
		T = T.Default().(*TaoConfig)
	} else {
		err = json.Unmarshal(bytes, &T)
		if err != nil {
			panic(err)
		}
	}

	// tao config
	T.ValidSelf()

	switch T.Type {
	case File:
		TaoWriter, err = os.OpenFile(T.Path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			panic(err)
		}
	default:
		TaoWriter = os.Stdout
	}
	TaoLogger = &logger{log.New(TaoWriter, "", log.LstdFlags|log.Lshortfile)}

	tao = NewPipeline(ConfigKey)

	// print banner
	banner := `
___________              
\__    ___/____    ____  
  |    |  \__  \  /  _ \ 
  |    |   / __ \(  <_> )
  |____|  (____  /\____/ 
               \/
`
	if !T.HideBanner {
		fmt.Print(banner)
	}
}
