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
	"os"
	"os/signal"
	"reflect"
	"sync"
	"syscall"
)

// Universe of tao
type Universe struct {
	sync.WaitGroup

	Pipeline
	universe Pipeline
}

// The Tao produced One; One produced Two; Two produced Three; Three produced All things.
var tao = &Universe{
	Pipeline: NewPipeline(ConfigKey),
	universe: NewPipeline("universe"),
}

// Add of tao
var Add = tao.Add

// Done of tao
var Done = tao.Done

// Run tao
func Run(ctx context.Context, param Parameter) (err error) {
	if ctx == nil {
		ctx = context.Background()
	}

	if param == nil {
		param = NewParameter()
	}

	if len(once) == 0 {
		return NewError(UniverseNotInit, "none of %+v existed", defaultConfigs)
	}

	select {
	case <-ctx.Done():
		return NewError(ContextCanceled, "tao: context has been canceled")
	default:
	}

	for _, c := range configMap {
		err = tao.Register(NewPipeTask(c.ToTask(), c.RunAfter()...))
		if err != nil {
			return NewErrorWrapped("tao: fail to register unit task", err)
		}
	}

	cm, err := json.MarshalIndent(configMap, "", "  ")
	if err != nil {
		return NewErrorWrapped("tao: fail to marshal configmap", err)
	}
	if configPath != "" {
		Debugf("load config from %q", configPath)
	}
	Debugf("config data: \n%s", string(cm))

	gracefulShutdown()

	err = tao.Run(ctx, param)
	if err != nil {
		return NewErrorWrapped("tao: fail to run", err)
	}

	tao.Wait()
	return
}

// Register unit to tao universe.
//
// Generic factory pattern registration with MultiConfig:
//   - Parses multi-instance config (single-instance config is auto-wrapped as {default: config})
//   - Calls constructor for each instance
func Register[T any, C any](
	configKey string,
	config MultiConfig[C],
	constructor func(name string, cfg C) (T, func() error, error),
) (*BaseFactory[T], error) {
	if constructor == nil {
		return nil, NewError(ParamInvalid, "tao: constructor is nil")
	}

	rv := reflect.ValueOf(config)
	if rv.Kind() != reflect.Pointer || rv.IsNil() {
		return nil, NewError(ParamInvalid, "tao: type of config should be pointer(notnull) instead of %+v", config)
	}

	factory := NewBaseFactory[T]()

	setup := func() error {
		err := LoadConfig(configKey, config)
		if err != nil {
			if e, ok := err.(ErrorTao); ok && e.Code() == ConfigNotFound {
			} else {
				return NewErrorWrapped(fmt.Sprintf("tao: fail to load config by key %q", configKey), err)
			}
		}

		if err := parseMultiConfig(configKey, config); err != nil {
			return err
		}

		config.ValidSelf()

		instances := config.GetInstances()
		if len(instances) == 0 {
			var zeroC C
			config.SetInstances([]Instance[C]{{Name: DefaultInstanceKey, Cfg: zeroC}})
			config.ValidSelf()
			instances = config.GetInstances()
		}

		var createdInstances []string
		for _, inst := range instances {
			instance, closer, err := constructor(inst.Name, inst.Cfg)
			if err != nil {
				for _, createdName := range createdInstances {
					_ = factory.Close(createdName)
				}
				return NewErrorWrapped(fmt.Sprintf("factory: failed to create instance %q", inst.Name), err)
			}
			if err := factory.RegisterWithCloser(inst.Name, instance, closer); err != nil {
				for _, createdName := range createdInstances {
					_ = factory.Close(createdName)
				}
				return err
			}
			createdInstances = append(createdInstances, inst.Name)
		}

		return SetConfig(configKey, config)
	}

	if err := registerToUniverse(configKey, config, setup); err != nil {
		return nil, err
	}

	return factory, nil
}

func registerToUniverse(configKey string, config Config, setup func() error) error {
	if config != nil && configKey != config.Name() {
		return NewError(ParamInvalid, "universe: config's name should be same as task's name")
	}

	if configKey == ConfigKey {
		return setup()
	}

	switch tao.universe.State() {
	case Running, Over, Closed:
		return setup()
	default:
		return tao.universe.Register(NewPipeTask(NewTask(configKey, func(ctx context.Context, param Parameter) (Parameter, error) {
			select {
			case <-ctx.Done():
				return param, NewError(ContextCanceled, "universe: fail to init %q", configKey)
			default:
				return param, setup()
			}
		})))
	}
}

func gracefulShutdown() {
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
	go func() {
		for sig := range sc {
			Debugf("got exiting signal now: %v", sig)
			if err := tao.Close(); err != nil {
				os.Exit(1)
			} else {
				os.Exit(0)
			}
		}
	}()
}
