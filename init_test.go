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
	"testing"
)

func TestSetConfigBytesAll(t *testing.T) {
	file := `
{
    "tao": {
        "log": {
            "level": "debug",
            "type": "console"
        },
        "hide_banner": false
    },
    "print": {
        "print": "==============  hello,tao!  ==============",
        "times": 2,
        "run_after": []
    }
}`
	err := SetConfigBytesAll([]byte(file), Json)
	assert.Nil(t, err)
	_, err = GetConfigBytes(printConfigKey)
	assert.Nil(t, err)

	// no use
	DevelopMode()
}
