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
	"context"
	"testing"
)

// flag.Parse() used in init would lead to testing failed.
// https://github.com/golang/go/issues/31859
var _ = func() bool {
	testing.Init()
	return true
}()

func TestMain(m *testing.M) {
	err := Run(context.Background(), nil)
	if err != nil {
		Panic(err)
	}

	m.Run()
}
