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
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
	p := &printConfig{}
	_, err := Register[struct{}, printConfigInstance](printConfigKey, p, func(name string, cfg printConfigInstance) (struct{}, func() error, error) {
		return struct{}{}, nil, nil
	})
	assert.Nil(t, err)

	_, err = Register[struct{}, struct{}](printConfigKey, nil, nil)
	assert.NotNil(t, err)

	err = SetConfig(printConfigKey, nil)
	assert.NotNil(t, err)

	err = universeInit()
	assert.NotNil(t, err)
}

func TestRun(t *testing.T) {
	t.Log(new(taoConfig).ToTask())
	t.Log(new(taoConfig).RunAfter())

	cancel, cancelFunc := context.WithCancel(context.Background())
	cancelFunc()
	err := Run(cancel, nil)
	assert.NotNil(t, err)

	Add(1)
	Done()

	ctx := context.Background()
	err = Run(ctx, nil)
	assert.Nil(t, err)

	err = Run(ctx, nil)
	assert.NotNil(t, err)
}

const singleInstanceTestKey = "single_instance_test"

type singleInstanceCfgInstance struct {
	Port int    `json:"port"`
	Host string `json:"host"`
}

type singleInstanceCfg struct {
	BaseMultiConfig[singleInstanceCfgInstance]
}

func (s *singleInstanceCfg) Name() string       { return singleInstanceTestKey }
func (s *singleInstanceCfg) ValidSelf()         {}
func (s *singleInstanceCfg) ToTask() Task       { return nil }
func (s *singleInstanceCfg) RunAfter() []string { return nil }

func TestRegister_SingleInstancePassesLoadedConfig(t *testing.T) {
	if len(once) == 0 {
		err := SetAllConfigBytes([]byte(`{}`), None)
		assert.Nil(t, err)
	}

	var capturedConfig singleInstanceCfgInstance
	var captureOnce sync.Once

	cfg := &singleInstanceCfg{}

	factory, err := Register(singleInstanceTestKey, cfg, func(name string, c singleInstanceCfgInstance) (string, func() error, error) {
		captureOnce.Do(func() { capturedConfig = c })
		return "ok", func() error { return nil }, nil
	})
	assert.Nil(t, err)
	assert.NotNil(t, factory)

	instance, err := factory.Get("default")
	assert.Nil(t, err)
	assert.Equal(t, "ok", instance)

	assert.NotNil(t, capturedConfig, "constructor should have received the config, not zero value")
}
