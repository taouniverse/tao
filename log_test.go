// Copyright 2021 huija
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

func TestLogger(t *testing.T) {
	t.Run("LogLevelMarshal", func(t *testing.T) {
		t.Log(DEBUG.String())
		t.Log(INFO.String())
		t.Log(WARNING.String())
		t.Log(ERROR.String())
		t.Log(PANIC.String())
		t.Log(FATAL.String())
		t.Log((FATAL + 1).String())

		marshal, err := json.Marshal(DEBUG)
		assert.Nil(t, err)
		t.Log(string(marshal))
		assert.Equal(t, "\"debug\"", string(marshal))

		err = json.Unmarshal(marshal, nil)
		assert.NotNil(t, err)

		var l = new(LogLevel)
		err = json.Unmarshal([]byte("\""+DEBUG.String()+"\""), &l)
		assert.Nil(t, err)
		assert.Equal(t, DEBUG, *l)

		err = json.Unmarshal([]byte("\""+INFO.String()+"\""), &l)
		assert.Nil(t, err)
		assert.Equal(t, INFO, *l)

		err = json.Unmarshal([]byte("\""+WARNING.String()+"\""), &l)
		assert.Nil(t, err)
		assert.Equal(t, WARNING, *l)

		err = json.Unmarshal([]byte("\""+ERROR.String()+"\""), &l)
		assert.Nil(t, err)
		assert.Equal(t, ERROR, *l)

		err = json.Unmarshal([]byte("\""+PANIC.String()+"\""), &l)
		assert.Nil(t, err)
		assert.Equal(t, PANIC, *l)

		err = json.Unmarshal([]byte("\""+FATAL.String()+"\""), &l)
		assert.Nil(t, err)
		assert.Equal(t, FATAL, *l)

		err = json.Unmarshal([]byte("\"unknown\""), &l)
		assert.NotNil(t, err)
	})

	t.Run("LogTypeMarshal", func(t *testing.T) {
		t.Log((Console - 1).String())
		t.Log(Console.String())
		t.Log(File.String())
		t.Log((Console | File).String())

		marshal, err := json.Marshal(Console)
		assert.Nil(t, err)
		t.Log(string(marshal))
		assert.Equal(t, "\"console\"", string(marshal))

		err = json.Unmarshal(marshal, nil)
		assert.NotNil(t, err)

		var l = new(LogType)
		err = json.Unmarshal([]byte("\""+Console.String()+"\""), &l)
		assert.Nil(t, err)
		assert.Equal(t, Console, *l)

		err = json.Unmarshal([]byte("\""+File.String()+"\""), &l)
		assert.Nil(t, err)
		assert.Equal(t, File, *l)

		err = json.Unmarshal([]byte("\""+(Console|File).String()+"\""), &l)
		assert.Nil(t, err)
		assert.Equal(t, Console|File, *l)

		err = json.Unmarshal([]byte("\"file|console\""), &l)
		assert.Nil(t, err)
		assert.Equal(t, Console|File, *l)

		err = json.Unmarshal([]byte("\"unknown\""), &l)
		assert.NotNil(t, err)
	})

	t.Run("ObjectFunction", func(t *testing.T) {
		GetLogger(ConfigKey).Debug("debug")
		GetLogger(ConfigKey).Debugf("%s", "debug")
		GetLogger(ConfigKey).Info("info")
		GetLogger(ConfigKey).Infof("%s", "info")
		GetLogger(ConfigKey).Warn("warn")
		GetLogger(ConfigKey).Warnf("%s", "warn")
		GetLogger(ConfigKey).Error("error")
		GetLogger(ConfigKey).Errorf("%s", "error")
		// GetLogger(ConfigKey).Panic("panic")
		// GetLogger(ConfigKey).Panicf("%s", "panic")
		// GetLogger(ConfigKey).Fatal("fatal")
		// GetLogger(ConfigKey).Fatalf("%s", "fatal")
	})

	t.Run("PackageFunction", func(t *testing.T) {
		Debug("debug")
		Debugf("%s", "debug")
		Info("info")
		Infof("%s", "info")
		Warn("warn")
		Warnf("%s", "warn")
		Error("error")
		Errorf("%s", "error")
		// Panic("panic")
		// Panicf("%s", "panic")
		// Fatal("fatal")
		// Fatalf("%s", "fatal")
	})

	t.Run("Writer", func(t *testing.T) {
		writer := GetWriter(ConfigKey)
		assert.NotNil(t, writer)

		assert.Nil(t, DeleteWriter(ConfigKey))
		assert.Nil(t, globalLogger.writers[ConfigKey])
		assert.NotNil(t, DeleteWriter(ConfigKey))

		assert.Nil(t, SetWriter(ConfigKey, writer))
		assert.NotNil(t, SetWriter(ConfigKey, writer))
	})

	t.Run("Logger", func(t *testing.T) {
		logger := GetLogger(ConfigKey)
		assert.NotNil(t, logger)

		assert.Nil(t, DeleteLogger(ConfigKey))
		assert.Nil(t, GetLogger(ConfigKey))
		assert.NotNil(t, DeleteLogger(ConfigKey))

		assert.Nil(t, SetLogger(ConfigKey, logger))
		assert.NotNil(t, SetLogger(ConfigKey, logger))
	})
}
