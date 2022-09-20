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
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestInit(t *testing.T) {
	t.Run("TestBeforeInit", func(t *testing.T) {
		err := Register("before all", nil, func() error {
			t.Log("before tao universe init")
			return nil
		})
		assert.Nil(t, err)

		Fatal("fatal before all")
		Fatalf("%s before all", "fatal")
	})

	file := []byte(`
{
    "tao": {
        "log": {
            "level": "debug"
        },
        "banner": {
            "hide": false
		}
    },
    "print": {
        "print": "==============  hello,tao!  ==============",
        "times": 2,
        "run_after": []
    }
}`)

	t.Run("TestSetConfigBytesAll", func(t *testing.T) {
		err := SetConfigBytesAll(file, JSON)
		assert.Nil(t, err)

		err = SetConfigBytesAll(file, Yaml)
		assert.NotNil(t, err)

		_, err = GetConfigBytes(printConfigKey)
		assert.Nil(t, err)
	})

	t.Run("TestSetConfigPath", func(t *testing.T) {
		err := DevelopMode()
		assert.NotNil(t, err)

		err = os.WriteFile("conf.yaml", file, 0666)
		assert.Nil(t, err)

		err = SetConfigPath("conf.yaml")
		assert.NotNil(t, err)

		err = os.Rename("conf.yaml", "conf.json")
		assert.Nil(t, err)

		err = SetConfigPath("conf.json")
		assert.NotNil(t, err)

		err = os.Rename("conf.json", "conf.unknown")
		assert.Nil(t, err)

		err = SetConfigPath("conf.unknown")
		assert.NotNil(t, err)

		err = os.Remove("conf.unknown")
		assert.Nil(t, err)
	})
}
