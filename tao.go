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
	"os"
	"os/signal"
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
		// refer to defaultConfigs in init.go to get some help
		return NewError(UniverseNotInit, "none of %+v existed", defaultConfigs)
	}

	// non-block check
	select {
	case <-ctx.Done():
		return NewError(ContextCanceled, "tao: context has been canceled")
	default:
	}

	// tasks register
	for _, c := range configMap {
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
	if configPath != "" {
		Debugf("load config from %q\n", configPath)
	}
	Debugf("config data: \n%s", string(cm))

	// graceful shutdown
	gracefulShutdown()

	// tao run
	err = tao.Run(ctx, param)
	if err != nil {
		return err
	}

	// tao wait
	tao.Wait()
	return
}

// Register unit to tao universe
func Register(configKey string, config Config, setup func() error) error {
	unitSetup := func() (err error) {
		defer func() {
			if err != nil {
				return
			}

			// 3. setup unit
			if setup != nil {
				err = setup()
			}
		}()

		if config == nil {
			return
		}

		// 1. transfer config bytes to object
		bytes, err := GetConfigBytes(configKey)
		if err == nil {
			err = json.Unmarshal(bytes, &config)
			if err != nil {
				return err
			}
		}

		// 2. set object to tao
		config.ValidSelf()
		return SetConfig(configKey, config)
	}

	if configKey == ConfigKey {
		// tao init
		return unitSetup()
	}

	switch tao.universe.State() {
	case Running, Over, Closed:
		return unitSetup()
	default:
		return tao.universe.Register(NewPipeTask(NewTask(configKey, func(ctx context.Context, param Parameter) (Parameter, error) {
			select {
			case <-ctx.Done():
				return param, NewError(ContextCanceled, "universe: %s init failed", configKey)
			default:
				return param, unitSetup()
			}
		})))
	}
}

func gracefulShutdown() {
	sc := make(chan os.Signal, 1)
	signal.Notify(sc)
	go func() {
		for {
			sig := <-sc
			if _, ok := map[os.Signal]struct{}{
				syscall.SIGINT:  {},
				syscall.SIGQUIT: {},
				syscall.SIGTERM: {},
			}[sig]; ok {
				Debugf("got exiting signal now: %v\n", sig)
				if err := tao.Close(); err != nil {
					os.Exit(1)
				} else {
					os.Exit(0)
				}
			} else {
				Debugf("got non-exiting signal: %v\n", sig)
			}
		}
	}()
}
