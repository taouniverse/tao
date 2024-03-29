// Copyright 2022 huija
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
	"gopkg.in/yaml.v3"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
	"time"
)

// ConfigType of config file
type ConfigType uint8

const (
	// None of config
	None ConfigType = iota
	// Yaml config
	Yaml
	// JSON config
	JSON
)

// List of default config files, traverse all until one is found
// if no one found, you can config in go code.(SetConfigPath || SetAllConfigBytes)
// if you do not want to config anything, call DevelopMode() to use default configs.
var defaultConfigs = []string{
	"./conf/config.yaml",
	"./conf/config.json",
	"./conf/config.yml",
}

func init() {
	for _, confPath := range defaultConfigs {
		_ = SetConfigPath(confPath)
	}
}

var configPath = ""

// SetConfigPath in your project's init()
func SetConfigPath(confPath string) error {
	data, err := ioutil.ReadFile(confPath)
	if err != nil {
		return NewErrorWrapped("init: fail to read config file", err)
	}

	switch t := path.Ext(confPath); t {
	case ".yaml", ".yml":
		err = SetAllConfigBytes(data, Yaml)
	case ".json":
		err = SetAllConfigBytes(data, JSON)
	default:
		return NewError(ParamInvalid, "%s file not supported", t)
	}
	if err != nil {
		return NewErrorWrapped("init: fail to set config path", err)
	}
	configPath = confPath
	return err
}

// DevelopMode called to enable default configs for all
func DevelopMode() error {
	if len(once) != 0 {
		return NewError(DuplicateCall, "tao: init twice")
	}

	return SetAllConfigBytes(nil, None)
}

// SetAllConfigBytes & taoInit can only be called once
var once = make(chan struct{}, 1)

// SetAllConfigBytes from config file or code
func SetAllConfigBytes(data []byte, configType ConfigType) (err error) {
	select {
	case once <- struct{}{}:
		switch configType {
		case Yaml:
			err = yaml.Unmarshal(data, &configInterfaceMap)
		case JSON:
			err = json.Unmarshal(data, &configInterfaceMap)
		default:
		}
		if err == nil {
			// init tao with config
			err = Register(ConfigKey, t, taoInit)
		}
	default:
		// caused by duplicate config(file & code)
		err = NewError(DuplicateCall, "config: SetConfigBytes has been called before")
	}
	return
}

// t global config of tao
var t = new(taoConfig)

// taoInit can only be called once before tao.Run
func taoInit() (err error) {
	// SetLogger
	if !t.Log.Disable {
		writers := make([]io.Writer, 0)

		if t.Log.Type&Console != 0 {
			writers = append(writers, os.Stdout)
		}

		if t.Log.Type&File != 0 {
			file, err := os.OpenFile(t.Log.Path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
			if err != nil {
				return NewErrorWrapped("init: fail to open log file", err)
			}
			writers = append(writers, file)
		}

		writer := io.MultiWriter(writers...)
		err = SetWriter(ConfigKey, writer)
		if err != nil {
			return NewErrorWrapped("init: fail to set writer for 'tao'", err)
		}

		err = SetLogger(ConfigKey, &logger{Logger: log.New(writer, "", int(t.Log.Flag)), calldepth: t.Log.CallDepth})
		if err != nil {
			return NewErrorWrapped("init: fail to set logger for 'tao'", err)
		}
	}

	// print banner
	if !t.Banner.Hide {
		w := GetWriter(ConfigKey)
		if w == nil {
			w = os.Stdout
		}
		_, err = w.Write([]byte(strings.TrimSpace(t.Banner.Content) + "\n"))
		if err != nil {
			return NewErrorWrapped("init: fail to write banner of tao", err)
		}
	}

	// init universe after tao
	return universeInit()
}

func universeInit() error {
	if tao.universe.State() != Runnable {
		return NewError(TaskRunTwice, "universe: init twice")
	}
	// universe run
	timeout, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	return tao.universe.Run(timeout, nil)
}
