// Copyright 2021-2026 huija
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

	"github.com/mitchellh/mapstructure"
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
		return NewError(ConfigNotFound, "config: %q not found", configKey)
	}
	configDecoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Result:           config,
		TagName:          "json",
		WeaklyTypedInput: true,
		DecodeHook:       mapstructure.TextUnmarshallerHookFunc(),
	})
	if err != nil {
		return NewErrorWrapped("config: fail to create decoder", err)
	}
	if err := configDecoder.Decode(c); err != nil {
		return NewErrorWrapped(fmt.Sprintf("config: fail to decode config for %q", configKey), err)
	}
	return nil
}

func looksLikeSingleInstance[C any](rawMap map[string]interface{}) bool {
	var zeroC C
	bytes, err := json.Marshal(zeroC)
	if err != nil {
		return false
	}
	var structFields map[string]interface{}
	if err := json.Unmarshal(bytes, &structFields); err != nil {
		return false
	}
	for k := range rawMap {
		if reservedFields[k] {
			continue
		}
		if _, exists := structFields[k]; exists {
			return true
		}
	}
	return false
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

// DefaultInstanceKey is the default instance name for single-instance configs
const DefaultInstanceKey = "default"

// taoInstanceConfig 单实例配置
type taoInstanceConfig struct {
	Log    *Log    `json:"log"`
	Banner *Banner `json:"banner"`
}

// Banner config
type Banner struct {
	Hide    bool   `json:"hide"`
	Content string `json:"content"`
}

var defaultTaoInstance = &taoInstanceConfig{
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

// taoConfig implements MultiConfig
type taoConfig struct {
	BaseMultiConfig[taoInstanceConfig]
}

// Name of Config
func (t *taoConfig) Name() string {
	return ConfigKey
}

// ValidSelf with some default values
func (t *taoConfig) ValidSelf() {
	for name, inst := range t.Instances {
		if inst.Log == nil {
			inst.Log = defaultTaoInstance.Log
		} else {
			if inst.Log.Level < DEBUG || inst.Log.Level > FATAL {
				inst.Log.Level = defaultTaoInstance.Log.Level
			}
			if inst.Log.CallDepth <= 0 {
				inst.Log.CallDepth = defaultTaoInstance.Log.CallDepth
			}
			if inst.Log.Type == 0 {
				inst.Log.Type = defaultTaoInstance.Log.Type
			}
			if inst.Log.Type&File != 0 {
				if inst.Log.Path == "" {
					inst.Log.Path = defaultTaoInstance.Log.Path
				}
			}
			if inst.Log.Flag == 0 {
				inst.Log.Flag = defaultTaoInstance.Log.Flag
			}
		}
		if inst.Banner == nil {
			inst.Banner = defaultTaoInstance.Banner
		} else {
			if !inst.Banner.Hide && inst.Banner.Content == "" {
				inst.Banner.Content = defaultTaoInstance.Banner.Content
			}
		}
		t.Instances[name] = inst
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

// MultiConfig 支持多实例的配置接口
type MultiConfig[C any] interface {
	Config
	GetInstances() map[string]C
	SetInstances(instances map[string]C)
	GetDefaultInstanceName() string
	SetDefaultInstanceName(name string)
}

// BaseMultiConfig 多实例配置基础实现
type BaseMultiConfig[C any] struct {
	Instances map[string]C `json:"-" yaml:"-"`
	Default   string       `json:"default_instance" yaml:"default_instance"`
}

// GetInstances returns the map of instances
func (b *BaseMultiConfig[C]) GetInstances() map[string]C {
	return b.Instances
}

// SetInstances sets the map of instances
func (b *BaseMultiConfig[C]) SetInstances(instances map[string]C) {
	b.Instances = instances
}

// GetDefaultInstanceName returns the default instance name
func (b *BaseMultiConfig[C]) GetDefaultInstanceName() string {
	if b.Default == "" {
		return DefaultInstanceKey
	}
	return b.Default
}

// SetDefaultInstanceName sets the default instance name
func (b *BaseMultiConfig[C]) SetDefaultInstanceName(name string) {
	b.Default = name
}

// MarshalJSON implements json.Marshaler to restore original config file structure.
// - Single instance (only "default" key): outputs flat config for easy editing
// - Multi-instance: outputs named instance keys with default_instance field
// Includes ValidSelf supplemented defaults
func (b *BaseMultiConfig[C]) MarshalJSON() ([]byte, error) {
	if len(b.Instances) == 0 {
		return []byte(`{}`), nil
	}

	if len(b.Instances) == 1 {
		if single, ok := b.Instances[DefaultInstanceKey]; ok {
			return json.Marshal(single)
		}
	}

	result := make(map[string]interface{}, len(b.Instances)+1)
	for name, inst := range b.Instances {
		result[name] = inst
	}
	if b.Default != "" {
		result["default_instance"] = b.Default
	}

	return json.Marshal(result)
}

var reservedFields = map[string]bool{
	"run_after":        true,
	"default_instance": true,
}

func parseMultiConfig[C any](configKey string, config MultiConfig[C]) error {
	rawConfig, ok := configInterfaceMap[configKey]
	if !ok {
		return nil
	}

	rawMap, ok := rawConfig.(map[string]interface{})
	if !ok {
		var single C
		if err := decodeMapToStruct(rawConfig, &single); err != nil {
			return NewErrorWrapped("factory: fail to decode single config", err)
		}
		config.SetInstances(map[string]C{DefaultInstanceKey: single})
		return nil
	}

	if looksLikeSingleInstance[C](rawMap) {
		var single C
		if err := decodeMapToStruct(rawMap, &single); err != nil {
			return NewErrorWrapped("factory: fail to decode single config", err)
		}
		config.SetInstances(map[string]C{DefaultInstanceKey: single})
		return nil
	}

	instances := make(map[string]C)
	for name, v := range rawMap {
		if reservedFields[name] {
			continue
		}

		var instance C
		if err := decodeMapToStruct(v, &instance); err != nil {
			return NewErrorWrapped(fmt.Sprintf("factory: fail to decode instance %q", name), err)
		}
		instances[name] = instance
	}

	if len(instances) == 0 {
		var zeroC C
		config.SetInstances(map[string]C{DefaultInstanceKey: zeroC})
		return nil
	}

	config.SetInstances(instances)

	// Parse and validate default_instance
	if defaultVal, ok := rawMap["default_instance"]; ok {
		if defaultStr, ok := defaultVal.(string); ok && defaultStr != "" {
			if _, ok := instances[defaultStr]; !ok {
				return NewError(ParamInvalid, "factory: default_instance %q not found in instances", defaultStr)
			}
			config.SetDefaultInstanceName(defaultStr)
		}
	}

	return nil
}

// decodeMapToStruct decodes input map to output struct using mapstructure
func decodeMapToStruct(input interface{}, output interface{}) error {
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Result:           output,
		TagName:          "json",
		WeaklyTypedInput: true,
		DecodeHook:       mapstructure.TextUnmarshallerHookFunc(),
	})
	if err != nil {
		return err
	}
	return decoder.Decode(input)
}
