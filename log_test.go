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
		taoLogger.loggers[ConfigKey].Debug(2, "debug")
		taoLogger.loggers[ConfigKey].Debugf(2, "%s", "debug")
		taoLogger.loggers[ConfigKey].Info(2, "info")
		taoLogger.loggers[ConfigKey].Infof(2, "%s", "info")
		taoLogger.loggers[ConfigKey].Warn(2, "warn")
		taoLogger.loggers[ConfigKey].Warnf(2, "%s", "warn")
		taoLogger.loggers[ConfigKey].Error(2, "error")
		taoLogger.loggers[ConfigKey].Errorf(2, "%s", "error")
		// taoLogger.loggers[ConfigKey].Panic(2, "panic")
		// taoLogger.loggers[ConfigKey].Panicf(2, "%s", "panic")
		// taoLogger.loggers[ConfigKey].Fatal(2, "fatal")
		// taoLogger.loggers[ConfigKey].Fatalf(2, "%s", "fatal")
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
