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
	"flag"
	"io/ioutil"
	"path"

	"gopkg.in/yaml.v3"
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

// Config Type
type ConfigType uint8

const (
	Yaml ConfigType = iota
	Json
)

// default yaml config
const defaultYamlConfig = "./conf/config.yaml"

func init() {
	// xxx -f conf/config.yaml
	confPath := flag.String("f", defaultYamlConfig, "config file path")

	flag.Parse()

	data, err := ioutil.ReadFile(*confPath)
	if err != nil {
		// 1. config in code
		// 2. use default config(default yaml not existed)
		return
	}

	// 1. config in file
	// 2. use default config(default yaml is existed but empty)
	switch t := path.Ext(*confPath); t {
	case ".yaml", ".yml":
		err = SetConfigBytesAll(data, Yaml)
	case ".json":
		err = SetConfigBytesAll(data, Json)
	default:
		panic(NewError(ParamInvalid, "%s file not supported", t))
	}
	if err != nil {
		panic(err)
	}
}

// SetConfigBytesAll & taoInit can only be called once
var once = make(chan struct{}, 1)

// SetConfigBytesAll from config file
func SetConfigBytesAll(data []byte, configType ConfigType) (err error) {
	select {
	case once <- struct{}{}:
		switch configType {
		case Yaml:
			err = yaml.Unmarshal(data, &configInterfaceMap)
		case Json:
			err = json.Unmarshal(data, &configInterfaceMap)
		default:
		}
		if err == nil {
			// init tao with config
			taoInit()
		}
	default:
		err = NewError(DuplicateCall, "config: SetConfigBytes has been called before")
	}
	return
}

// GetConfigBytes by key of config
func GetConfigBytes(key string) ([]byte, error) {
	c, ok := configInterfaceMap[key]
	if !ok {
		return nil, NewError(ConfigNotFound, "config: %s not found", key)
	}
	bytes, err := json.Marshal(c)
	if err != nil {
		return nil, NewErrorUnWrapper("config: marshal failed", err)
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
