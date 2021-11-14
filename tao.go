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
	"log"
	"os"
)

// The Tao produced One; One produced Two; Two produced Three; Three produced All things.
var tao Pipeline

// init default logger & writer
func init() {
	TaoWriter = os.Stdout
	TaoLogger = &logger{log.New(TaoWriter, " ", log.LstdFlags|log.Lshortfile)}

	tao = NewPipeline("tao")
}

// Run tao
func Run(ctx context.Context, param Parameter) (err error) {
	if ctx == nil {
		ctx = context.Background()
	}

	if param == nil {
		param = NewParameter()
	}

	for _, c := range configMap {
		c.ValidSelf()
		err = tao.Register(NewPipeTask(c.ToTask(), c.RunAfter()...))
		if err != nil {
			return err
		}
	}

	return tao.Run(ctx, param)
}
