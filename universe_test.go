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
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRegister(t *testing.T) {
	err := Register(printConfigKey, func() error {
		p := new(printConfig)
		// 1. transfer config bytes to object
		bytes, err := GetConfigBytes(printConfigKey)
		if err != nil {
			return err
		}
		err = json.Unmarshal(bytes, &p)
		if err != nil {
			return err
		}

		p.ValidSelf()

		// 2. set object to tao
		return SetConfig(printConfigKey, p)
	})
	assert.Nil(t, err)

	err = SetConfig(printConfigKey, nil)
	assert.NotNil(t, err)
}
