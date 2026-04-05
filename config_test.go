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
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJsonConfig(t *testing.T) {
	err := LoadConfig("Unknown", nil)
	assert.NotNil(t, err)
	assert.Equal(t, ConfigNotFound, err.(ErrorTao).Code())
}

const printConfigKey = "print"

// printConfigInstance 单实例配置
type printConfigInstance struct {
	Print     string   `json:"print"`
	Times     int      `json:"times"`
	RunAfters []string `json:"run_after"`
}

// printConfig implements MultiConfig
type printConfig struct {
	BaseMultiConfig[printConfigInstance]
}

var defaultPrint = &printConfigInstance{
	Print: "==============  hello,tao!  ==============",
	Times: 1,
}

// Name of Config
func (l *printConfig) Name() string {
	return printConfigKey
}

// ValidSelf with some default values
func (l *printConfig) ValidSelf() {
	for name, inst := range l.Instances {
		if inst.Print == "" {
			inst.Print = defaultPrint.Print
		}
		if inst.Times == 0 {
			inst.Times = defaultPrint.Times
		}
		if inst.RunAfters == nil {
			inst.RunAfters = defaultPrint.RunAfters
		}
		l.Instances[name] = inst
	}
}

// ToTask transform itself to Task
func (l *printConfig) ToTask() Task {
	return NewTask("print", func(ctx context.Context, param Parameter) (Parameter, error) {
		select {
		case <-ctx.Done():
			return param, NewError(ContextCanceled, "test: ctx already Done")
		default:
			for _, inst := range l.Instances {
				for i := 0; i < inst.Times; i++ {
					fmt.Println(inst.Print)
				}
			}
			return param, nil
		}
	})
}

// RunAfter defines pre task names
func (l *printConfig) RunAfter() []string {
	return nil
}

func TestTaoConfigValidSelf(t *testing.T) {
	t.Run("TestValidSelf_NilLog", func(t *testing.T) {
		config := &taoConfig{}
		config.SetInstances(map[string]taoInstanceConfig{"default": {}})
		config.ValidSelf()
		inst := config.Instances["default"]
		assert.NotNil(t, inst.Log)
		assert.NotNil(t, inst.Banner)
	})

	t.Run("TestValidSelf_InvalidLogLevel", func(t *testing.T) {
		config := &taoConfig{}
		config.SetInstances(map[string]taoInstanceConfig{"default": {
			Log: &Log{
				Level: 100,
			},
		}})
		config.ValidSelf()
		assert.Equal(t, defaultTaoInstance.Log.Level, config.Instances["default"].Log.Level)
	})

	t.Run("TestValidSelf_InvalidCallDepth", func(t *testing.T) {
		config := &taoConfig{}
		config.SetInstances(map[string]taoInstanceConfig{"default": {
			Log: &Log{
				CallDepth: 0,
			},
		}})
		config.ValidSelf()
		assert.Equal(t, defaultTaoInstance.Log.CallDepth, config.Instances["default"].Log.CallDepth)
	})

	t.Run("TestValidSelf_ZeroLogType", func(t *testing.T) {
		config := &taoConfig{}
		config.SetInstances(map[string]taoInstanceConfig{"default": {
			Log: &Log{
				Type: 0,
			},
		}})
		config.ValidSelf()
		assert.Equal(t, defaultTaoInstance.Log.Type, config.Instances["default"].Log.Type)
	})

	t.Run("TestValidSelf_FileTypeEmptyPath", func(t *testing.T) {
		config := &taoConfig{}
		config.SetInstances(map[string]taoInstanceConfig{"default": {
			Log: &Log{
				Type: File,
				Path: "",
			},
		}})
		config.ValidSelf()
		assert.Equal(t, defaultTaoInstance.Log.Path, config.Instances["default"].Log.Path)
	})

	t.Run("TestValidSelf_ZeroFlag", func(t *testing.T) {
		config := &taoConfig{}
		config.SetInstances(map[string]taoInstanceConfig{"default": {
			Log: &Log{
				Flag: 0,
			},
		}})
		config.ValidSelf()
		assert.Equal(t, defaultTaoInstance.Log.Flag, config.Instances["default"].Log.Flag)
	})

	t.Run("TestValidSelf_NilBanner", func(t *testing.T) {
		config := &taoConfig{}
		config.SetInstances(map[string]taoInstanceConfig{"default": {
			Log: defaultTaoInstance.Log,
		}})
		config.ValidSelf()
		assert.NotNil(t, config.Instances["default"].Banner)
	})

	t.Run("TestValidSelf_EmptyBannerContent", func(t *testing.T) {
		config := &taoConfig{}
		config.SetInstances(map[string]taoInstanceConfig{"default": {
			Log:    defaultTaoInstance.Log,
			Banner: &Banner{Hide: false, Content: ""},
		}})
		config.ValidSelf()
		assert.Equal(t, defaultTaoInstance.Banner.Content, config.Instances["default"].Banner.Content)
	})
}

func TestSetConfig(t *testing.T) {
	t.Run("TestSetConfig_Duplicate", func(t *testing.T) {
		testKey := "test_config_key"
		testConfig := &taoConfig{}
		testConfig.SetInstances(map[string]taoInstanceConfig{"default": {}})
		testConfig.ValidSelf()

		err := SetConfig(testKey, testConfig)
		assert.Nil(t, err)

		err = SetConfig(testKey, testConfig)
		assert.NotNil(t, err)
		assert.Equal(t, DuplicateCall, err.(ErrorTao).Code())
	})
}

// --- BaseMultiConfig tests ---

func TestBaseMultiConfig_GetDefaultInstanceName_Empty(t *testing.T) {
	cfg := BaseMultiConfig[int]{}
	assert.Equal(t, "default", cfg.GetDefaultInstanceName())
}

func TestBaseMultiConfig_GetDefaultInstanceName_Custom(t *testing.T) {
	cfg := BaseMultiConfig[int]{Default: "master"}
	assert.Equal(t, "master", cfg.GetDefaultInstanceName())
}

func TestBaseMultiConfig_GetSetInstances(t *testing.T) {
	type dbConfig struct{ Host string }

	cfg := BaseMultiConfig[dbConfig]{}
	assert.Nil(t, cfg.GetInstances())

	instances := map[string]dbConfig{
		"master":  {Host: "10.0.0.1"},
		"replica": {Host: "10.0.0.2"},
	}
	cfg.SetInstances(instances)

	got := cfg.GetInstances()
	assert.Len(t, got, 2)
	assert.Equal(t, "10.0.0.1", got["master"].Host)
	assert.Equal(t, "10.0.0.2", got["replica"].Host)
}

func TestBaseMultiConfig_SetInstances_Nil(t *testing.T) {
	cfg := BaseMultiConfig[int]{Instances: map[string]int{"a": 1}}
	cfg.SetInstances(nil)
	assert.Nil(t, cfg.GetInstances())
}

// --- parseMultiConfig tests ---

type instanceConfig struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}

type multiConfigImpl struct {
	BaseMultiConfig[instanceConfig]
	RunAfters []string `json:"run_after,omitempty"`
}

func (c *multiConfigImpl) Name() string       { return "multi_test" }
func (c *multiConfigImpl) ValidSelf()         {}
func (c *multiConfigImpl) ToTask() Task       { return nil }
func (c *multiConfigImpl) RunAfter() []string { return c.RunAfters }

type nestedInstanceConfig struct {
	Inner struct {
		A int    `json:"a"`
		B string `json:"b"`
	} `json:"inner"`
	Items []int `json:"items"`
}

type nestedMultiConfig struct {
	BaseMultiConfig[nestedInstanceConfig]
}

func (c *nestedMultiConfig) Name() string       { return "test_nested" }
func (c *nestedMultiConfig) ValidSelf()         {}
func (c *nestedMultiConfig) ToTask() Task       { return nil }
func (c *nestedMultiConfig) RunAfter() []string { return nil }

func TestParseMultiConfig_NotFoundKey(t *testing.T) {
	cfg := &multiConfigImpl{}
	err := parseMultiConfig("nonexistent_key_xyz", cfg)
	assert.NoError(t, err)
	assert.Nil(t, cfg.Instances)
}

func TestParseMultiConfig_FlatSingleInstance(t *testing.T) {
	configInterfaceMap["test_flat_single"] = map[string]interface{}{
		"name":  "single",
		"value": 42,
	}

	cfg := &multiConfigImpl{}
	err := parseMultiConfig("test_flat_single", cfg)
	assert.NoError(t, err)

	instances := cfg.GetInstances()
	assert.Len(t, instances, 1)
	assert.Contains(t, instances, "default")
	assert.Equal(t, "single", instances["default"].Name)
	assert.Equal(t, 42, instances["default"].Value)

	delete(configInterfaceMap, "test_flat_single")
}

func TestParseMultiConfig_MultipleInstances(t *testing.T) {
	configInterfaceMap["test_multi"] = map[string]interface{}{
		"master": map[string]interface{}{
			"name":  "master",
			"value": 1,
		},
		"replica": map[string]interface{}{
			"name":  "replica",
			"value": 2,
		},
	}

	cfg := &multiConfigImpl{}
	err := parseMultiConfig("test_multi", cfg)
	assert.NoError(t, err)

	instances := cfg.GetInstances()
	assert.Len(t, instances, 2)
	assert.Equal(t, "master", instances["master"].Name)
	assert.Equal(t, 1, instances["master"].Value)
	assert.Equal(t, "replica", instances["replica"].Name)
	assert.Equal(t, 2, instances["replica"].Value)

	delete(configInterfaceMap, "test_multi")
}

func TestParseMultiConfig_WithReservedFields(t *testing.T) {
	configInterfaceMap["test_reserved"] = map[string]interface{}{
		"master": map[string]interface{}{
			"name":  "master",
			"value": 99,
		},
		"run_after":        []string{"other"},
		"default_instance": "master",
	}

	cfg := &multiConfigImpl{}
	err := parseMultiConfig("test_reserved", cfg)
	assert.NoError(t, err)

	instances := cfg.GetInstances()
	assert.Len(t, instances, 1)
	assert.Contains(t, instances, "master")
	assert.NotContains(t, instances, "run_after")
	assert.NotContains(t, instances, "default_instance")
	// Verify default_instance is correctly parsed and assigned
	assert.Equal(t, "master", cfg.GetDefaultInstanceName())

	delete(configInterfaceMap, "test_reserved")
}

func TestParseMultiConfig_DefaultInstanceNotExists(t *testing.T) {
	configInterfaceMap["test_default_not_exists"] = map[string]interface{}{
		"master": map[string]interface{}{
			"name":  "master",
			"value": 1,
		},
		"default_instance": "nonexistent",
	}

	cfg := &multiConfigImpl{}
	err := parseMultiConfig("test_default_not_exists", cfg)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "default_instance")

	delete(configInterfaceMap, "test_default_not_exists")
}

func TestParseMultiConfig_FlatObjectNotMap(t *testing.T) {
	configInterfaceMap["test_flat_obj"] = "just a string value"

	cfg := &multiConfigImpl{}
	err := parseMultiConfig("test_flat_obj", cfg)
	assert.Error(t, err)

	delete(configInterfaceMap, "test_flat_obj")
}

func TestParseMultiConfig_SingleInstanceFlatFields(t *testing.T) {
	configInterfaceMap["test_single_flat"] = map[string]interface{}{
		"name":  "primary",
		"value": 100,
		"extra": "ignored_field",
	}

	cfg := &multiConfigImpl{}
	err := parseMultiConfig("test_single_flat", cfg)
	assert.NoError(t, err)

	instances := cfg.GetInstances()
	assert.Len(t, instances, 1)
	assert.Equal(t, "primary", instances["default"].Name)
	assert.Equal(t, 100, instances["default"].Value)

	delete(configInterfaceMap, "test_single_flat")
}

func TestParseMultiConfig_ManyInstances(t *testing.T) {
	raw := make(map[string]interface{})
	for i := 0; i < 50; i++ {
		name := fmt.Sprintf("node_%d", i)
		raw[name] = map[string]interface{}{
			"name":  name,
			"value": i,
		}
	}
	configInterfaceMap["test_many"] = raw

	cfg := &multiConfigImpl{}
	err := parseMultiConfig("test_many", cfg)
	assert.NoError(t, err)

	instances := cfg.GetInstances()
	assert.Len(t, instances, 50)
	for i := 0; i < 50; i++ {
		name := fmt.Sprintf("node_%d", i)
		assert.Contains(t, instances, name)
		assert.Equal(t, i, instances[name].Value)
	}

	delete(configInterfaceMap, "test_many")
}

func TestParseMultiConfig_AllReservedFieldsOnly(t *testing.T) {
	configInterfaceMap["test_only_reserved"] = map[string]interface{}{
		"run_after":        []string{"task_a"},
		"default_instance": "custom_default",
	}

	cfg := &multiConfigImpl{}
	err := parseMultiConfig("test_only_reserved", cfg)
	assert.NoError(t, err)

	assert.NotNil(t, cfg.Instances)
	assert.Len(t, cfg.Instances, 1)
	assert.Contains(t, cfg.Instances, "default")

	delete(configInterfaceMap, "test_only_reserved")
}

func TestParseMultiConfig_StructWithNestedFields(t *testing.T) {
	nestedCfg := &nestedMultiConfig{}

	configInterfaceMap["test_nested"] = map[string]interface{}{
		"db1": map[string]interface{}{
			"inner": map[string]interface{}{
				"a": 42,
				"b": "hello",
			},
			"items": []interface{}{1, 2, 3},
		},
	}

	err := parseMultiConfig("test_nested", nestedCfg)
	assert.NoError(t, err)

	instances := nestedCfg.GetInstances()
	assert.Len(t, instances, 1)
	assert.Equal(t, 42, instances["db1"].Inner.A)
	assert.Equal(t, "hello", instances["db1"].Inner.B)
	assert.Equal(t, []int{1, 2, 3}, instances["db1"].Items)

	delete(configInterfaceMap, "test_nested")
}

func TestReservedFields_IsComplete(t *testing.T) {
	assert.True(t, reservedFields["run_after"])
	assert.True(t, reservedFields["default_instance"])
	assert.Equal(t, 2, len(reservedFields))
}

func TestBaseMultiConfig_MarshalJSON_SingleInstance(t *testing.T) {
	cfg := BaseMultiConfig[instanceConfig]{}
	cfg.SetInstances(map[string]instanceConfig{
		"default": {Name: "primary", Value: 100},
	})

	data, err := json.MarshalIndent(&cfg, "", "  ")
	assert.NoError(t, err)

	var result instanceConfig
	err = json.Unmarshal(data, &result)
	assert.NoError(t, err)
	assert.Equal(t, "primary", result.Name)
	assert.Equal(t, 100, result.Value)

	assert.NotContains(t, string(data), "default")
	assert.NotContains(t, string(data), "instances")
}

func TestBaseMultiConfig_MarshalJSON_MultiInstance(t *testing.T) {
	cfg := BaseMultiConfig[instanceConfig]{}
	cfg.SetInstances(map[string]instanceConfig{
		"master":  {Name: "master", Value: 1},
		"replica": {Name: "replica", Value: 2},
		"read":    {Name: "read", Value: 3},
	})
	cfg.SetDefaultInstanceName("master")

	data, err := json.MarshalIndent(&cfg, "", "  ")
	assert.NoError(t, err)
	str := string(data)

	assert.Contains(t, str, `"master"`)
	assert.Contains(t, str, `"replica"`)
	assert.Contains(t, str, `"read"`)
	assert.Contains(t, str, `"default_instance"`)
	assert.Contains(t, str, `"master"`)

	var raw map[string]interface{}
	err = json.Unmarshal(data, &raw)
	assert.NoError(t, err)
	assert.Len(t, raw, 4) // 3 instances + default_instance
}
