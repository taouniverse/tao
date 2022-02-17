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
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v3"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
)

// ConfigType of config file
type ConfigType uint8

const (
	None ConfigType = iota
	Yaml
	Json
)

// List of default config files, traverse all until one is found
// if no one found, you can config in go code.(SetConfigPath || SetConfigBytesAll)
// if you do not want to config anything, call DevelopMode() to use default configs.
var defaultConfigs = []string{
	"./conf/config.yaml",
	"./conf/config.json",
	"./conf/config.yml",
}

func init() {
	for _, confPath := range defaultConfigs {
		SetConfigPath(confPath)
	}
}

// SetConfigPath in init of your project
func SetConfigPath(confPath string) {
	data, err := ioutil.ReadFile(confPath)
	if err != nil {
		return
	}

	switch t := path.Ext(confPath); t {
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

// DevelopMode called to enable default configs for all
func DevelopMode() {
	if tao != nil {
		return
	}

	err := SetConfigBytesAll(nil, None)
	if err != nil {
		panic(err)
	}
}

// SetConfigBytesAll & taoInit can only be called once
var once = make(chan struct{}, 1)

// SetConfigBytesAll from config file or code
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
		// caused by duplicate config(file & code)
		err = NewError(DuplicateCall, "config: SetConfigBytes has been called before")
	}
	return
}

// taoInit can only be called once before tao.Run
func taoInit() {
	// transfer config bytes to object
	t = new(TaoConfig)
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

	err = SetConfig(ConfigKey, t)
	if err != nil {
		panic(err)
	}

	writers := make([]io.Writer, 0)

	if t.Type&Console != 0 {
		writers = append(writers, os.Stdout)
	}

	if t.Type&File != 0 {
		file, err := os.OpenFile(t.Path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			panic(err)
		}
		writers = append(writers, file)
	}

	writer := io.MultiWriter(writers...)
	err = SetWriter(ConfigKey, writer)
	if err != nil {
		panic(err)
	}

	err = SetLogger(ConfigKey, &logger{log.New(writer, "", log.LstdFlags|log.Lshortfile)})
	if err != nil {
		panic(err)
	}

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
	if !t.HideBanner {
		fmt.Print(banner)
	}

	// init universe after tao
	universeInit()
}
