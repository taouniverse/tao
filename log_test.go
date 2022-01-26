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
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLogger(t *testing.T) {
	t.Run("ObjectFunction", func(t *testing.T) {
		GetLogger(ConfigKey).Debug(2, "debug")
		GetLogger(ConfigKey).Debugf(2, "%s", "debug")
		GetLogger(ConfigKey).Info(2, "info")
		GetLogger(ConfigKey).Infof(2, "%s", "info")
		GetLogger(ConfigKey).Warn(2, "warn")
		GetLogger(ConfigKey).Warnf(2, "%s", "warn")
		GetLogger(ConfigKey).Error(2, "error")
		GetLogger(ConfigKey).Errorf(2, "%s", "error")
		// GetLogger(ConfigKey).Panic(2, "panic")
		// GetLogger(ConfigKey).Panicf(2, "%s", "panic")
		// GetLogger(ConfigKey).Fatal(2, "fatal")
		// GetLogger(ConfigKey).Fatalf(2, "%s", "fatal")
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
		assert.Nil(t, taoLogger.writers[ConfigKey])

		assert.Nil(t, SetWriter(ConfigKey, writer))
	})

	t.Run("Logger", func(t *testing.T) {
		logger := GetLogger(ConfigKey)
		assert.NotNil(t, logger)

		assert.Nil(t, DeleteLogger(ConfigKey))
		assert.Nil(t, GetLogger(ConfigKey))

		assert.Nil(t, SetLogger(ConfigKey, logger))
	})
}
