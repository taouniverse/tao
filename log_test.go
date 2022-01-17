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
	"testing"
)

func TestLogger(t *testing.T) {
	t.Run("ObjectFunction", func(t *testing.T) {
		TaoLogger.Debug(2, "debug")
		TaoLogger.Debugf(2, "%s", "debug")
		TaoLogger.Info(2, "info")
		TaoLogger.Infof(2, "%s", "info")
		TaoLogger.Warn(2, "warn")
		TaoLogger.Warnf(2, "%s", "warn")
		TaoLogger.Error(2, "error")
		TaoLogger.Errorf(2, "%s", "error")
		// TaoLogger.Panic(2, "panic")
		// TaoLogger.Panicf(2, "%s", "panic")
		// TaoLogger.Fatal(2, "fatal")
		// TaoLogger.Fatalf(2, "%s", "fatal")
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
}
