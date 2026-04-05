// Copyright 2021-2026 huija
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
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewBaseFactory(t *testing.T) {
	factory := NewBaseFactory[int]()
	assert.NotNil(t, factory)
	assert.Equal(t, 0, factory.Count())
}

func TestFactoryRegisterAndGet(t *testing.T) {
	factory := NewBaseFactory[string]()

	err := factory.Register("default", "hello")
	assert.NoError(t, err)

	got, err := factory.Get("default")
	assert.NoError(t, err)
	assert.Equal(t, "hello", got)

	assert.True(t, factory.Has("default"))
	assert.False(t, factory.Has("not_exist"))
}

func TestFactoryGetNotFound(t *testing.T) {
	factory := NewBaseFactory[int]()

	_, err := factory.Get("not_exist")
	assert.Error(t, err)
}

func TestFactoryDuplicateRegister(t *testing.T) {
	factory := NewBaseFactory[int]()

	err := factory.Register("default", 42)
	assert.NoError(t, err)

	err = factory.Register("default", 99)
	assert.Error(t, err)
	assert.Equal(t, DuplicateCall, err.(ErrorTao).Code())
}

func TestFactoryNames(t *testing.T) {
	factory := NewBaseFactory[int]()

	err := factory.Register("a", 1)
	assert.NoError(t, err)

	err = factory.Register("b", 2)
	assert.NoError(t, err)

	err = factory.Register("c", 3)
	assert.NoError(t, err)

	names := factory.Names()
	assert.Len(t, names, 3)
	assert.Contains(t, names, "a")
	assert.Contains(t, names, "b")
	assert.Contains(t, names, "c")
}

func TestFactoryCount(t *testing.T) {
	factory := NewBaseFactory[int]()
	assert.Equal(t, 0, factory.Count())

	err := factory.Register("a", 1)
	assert.NoError(t, err)
	assert.Equal(t, 1, factory.Count())

	err = factory.Register("b", 2)
	assert.NoError(t, err)
	assert.Equal(t, 2, factory.Count())
}

func TestFactoryRegisterWithCloser(t *testing.T) {
	factory := NewBaseFactory[int]()
	closed := false

	err := factory.RegisterWithCloser("default", 42, func() error {
		closed = true
		return nil
	})
	assert.NoError(t, err)
	assert.False(t, closed)

	err = factory.Close("default")
	assert.NoError(t, err)
	assert.True(t, closed)
	assert.Equal(t, 0, factory.Count())
}

func TestFactoryCloseNotFound(t *testing.T) {
	factory := NewBaseFactory[int]()

	err := factory.Close("not_exist")
	assert.Error(t, err)
	assert.Equal(t, ParamInvalid, err.(ErrorTao).Code())
}

func TestFactoryCloseAll(t *testing.T) {
	factory := NewBaseFactory[int]()
	order := make([]string, 0)
	mu := sync.Mutex{}

	err := factory.RegisterWithCloser("a", 1, func() error {
		mu.Lock()
		defer mu.Unlock()
		order = append(order, "a")
		return nil
	})
	assert.NoError(t, err)

	err = factory.RegisterWithCloser("b", 2, func() error {
		mu.Lock()
		defer mu.Unlock()
		order = append(order, "b")
		return nil
	})
	assert.NoError(t, err)

	err = factory.RegisterWithCloser("c", 3, func() error {
		mu.Lock()
		defer mu.Unlock()
		order = append(order, "c")
		return nil
	})
	assert.NoError(t, err)

	err = factory.CloseAll()
	assert.NoError(t, err)
	assert.Equal(t, 0, factory.Count())
	assert.Len(t, order, 3)
}

func TestFactoryCloseAllWithError(t *testing.T) {
	factory := NewBaseFactory[int]()

	err := factory.RegisterWithCloser("a", 1, func() error {
		return fmt.Errorf("close a failed")
	})
	assert.NoError(t, err)

	err = factory.RegisterWithCloser("b", 2, func() error {
		return fmt.Errorf("close b failed")
	})
	assert.NoError(t, err)

	err = factory.CloseAll()
	assert.Error(t, err)
	assert.Equal(t, InstancesNotClosed, err.(ErrorTao).Code())
}

func TestFactoryCloseIdempotent(t *testing.T) {
	factory := NewBaseFactory[int]()
	closed := 0

	err := factory.RegisterWithCloser("default", 42, func() error {
		closed++
		return nil
	})
	assert.NoError(t, err)

	err = factory.Close("default")
	assert.NoError(t, err)
	assert.Equal(t, 1, closed)

	err = factory.Close("default")
	assert.Error(t, err)
	assert.Equal(t, ParamInvalid, err.(ErrorTao).Code())
	assert.Equal(t, 1, closed)
}

func TestFactoryConcurrentAccess(t *testing.T) {
	factory := NewBaseFactory[int]()
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			name := fmt.Sprintf("instance_%d", i)
			_ = factory.Register(name, i)
		}(i)
	}
	wg.Wait()

	assert.Equal(t, 100, factory.Count())

	var readWg sync.WaitGroup
	for i := 0; i < 100; i++ {
		readWg.Add(1)
		go func(i int) {
			defer readWg.Done()
			name := fmt.Sprintf("instance_%d", i)
			val, err := factory.Get(name)
			if err == nil {
				assert.Equal(t, i, val)
			}
		}(i)
	}
	readWg.Wait()
}

func TestFactoryConcurrentReadWrite(t *testing.T) {
	factory := NewBaseFactory[string]()
	err := factory.Register("default", "value")
	assert.NoError(t, err)

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, _ = factory.Get("default")
		}()
	}
	wg.Wait()

	got, err := factory.Get("default")
	assert.NoError(t, err)
	assert.Equal(t, "value", got)
}

func BenchmarkFactoryGet(b *testing.B) {
	factory := NewBaseFactory[int]()
	err := factory.Register("default", 42)
	assert.NoError(b, err)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err = factory.Get("default")
			assert.NoError(b, err)
		}
	})
}

func BenchmarkFactoryHas(b *testing.B) {
	factory := NewBaseFactory[int]()
	err := factory.Register("default", 42)
	assert.NoError(b, err)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			factory.Has("default")
		}
	})
}
