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
	"testing"
)

func TestLogger(t *testing.T) {
	t.Run("ObjectFunction", func(t *testing.T) {
		TaoLogger.Debug("debug")
		TaoLogger.Debugf("%s", "debug")
		TaoLogger.Info("info")
		TaoLogger.Infof("%s", "info")
		TaoLogger.Warn("warn")
		TaoLogger.Warnf("%s", "warn")
		TaoLogger.Error("error")
		TaoLogger.Errorf("%s", "error")
		// TaoLogger.Panic("panic")
		// TaoLogger.Panicf("%s", "panic")
		// TaoLogger.Fatal("fatal")
		// TaoLogger.Fatalf("%s", "fatal")
	})

	t.Run("PackageFunction", func(t *testing.T) {
		Debug()("debug")
		Debugf()("%s", "debug")
		Info()("info")
		Infof()("%s", "info")
		Warn()("warn")
		Warnf()("%s", "warn")
		Error()("error")
		Errorf()("%s", "error")
		// Panic()("panic")
		// Panicf()("%s", "panic")
		// Fatal()("fatal")
		// Fatalf()("%s", "fatal")
	})
}
