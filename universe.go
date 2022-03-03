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
	"time"
)

// universe for tao
var universe = NewPipeline("universe")

func universeInit() error {
	if universe.State() != Runnable {
		return NewError(TaskRunTwice, "universe: init twice")
	}
	// universe run
	timeout, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	return universe.Run(timeout, nil)
}

// Register to tao universe
func Register(configKey string, fn func() error) error {
	switch universe.State() {
	case Running, Over, Close:
		return fn()
	default:
		return universe.Register(NewPipeTask(NewTask(configKey, func(ctx context.Context, param Parameter) (Parameter, error) {
			select {
			case <-ctx.Done():
				return param, NewError(ContextCanceled, "universe: %s init failed", configKey)
			default:
				return param, fn()
			}
		})))
	}
}
