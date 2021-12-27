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
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

// The Tao produced One; One produced Two; Two produced Three; Three produced All things.
var tao Pipeline

// init tao
func init() {
	// default logger & writer
	TaoWriter = os.Stdout
	TaoLogger = &logger{log.New(TaoWriter, "", log.LstdFlags|log.Lshortfile)}

	loadConfig()

	taoConfig()
}

// Run tao
func Run(ctx context.Context, param Parameter) (err error) {
	if ctx == nil {
		ctx = context.Background()
	}

	if param == nil {
		param = NewParameter()
	}

	cm, err := json.MarshalIndent(configMap, "", "  ")
	if err != nil {
		return err
	}
	Debugf("config data: \n%s", string(cm))

	for _, c := range configMap {
		c.ValidSelf()
		err = tao.Register(NewPipeTask(c.ToTask(), c.RunAfter()...))
		if err != nil {
			return err
		}
	}

	return tao.Run(ctx, param)
}

// TaoConfig implements Config
type TaoConfig struct {
	LogLevel `json:"log_level"`
}

var defaultTao = &TaoConfig{
	DEBUG,
}

// Default config
func (t *TaoConfig) Default() Config {
	return defaultTao
}

// ValidSelf with some default values
func (t *TaoConfig) ValidSelf() {
	if t.LogLevel < DEBUG || t.LogLevel > FATAL {
		t.LogLevel = defaultTao.LogLevel
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

func taoConfig() {
	// transfer config bytes to object
	t := new(TaoConfig)
	bytes, err := GetConfigBytes(ConfigKey)
	if err != nil {
		t = t.Default().(*TaoConfig)
	} else {
		err = json.Unmarshal(bytes, &t)
		if err != nil {
			panic(err)
		}
	}

	// tao config
	t.ValidSelf()

	TaoLevel = t.LogLevel
	TaoWriter = os.Stdout
	TaoLogger = &logger{log.New(TaoWriter, "", log.LstdFlags|log.Lshortfile)}

	tao = NewPipeline("tao")

	// print banner
	banner := `
___________              
\__    ___/____    ____  
  |    |  \__  \  /  _ \ 
  |    |   / __ \(  <_> )
  |____|  (____  /\____/ 
               \/
`
	fmt.Print(banner)
}
