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
)

// Factory 通用工厂接口
type Factory[T any] interface {
	Get(name string) (T, error)
	Has(name string) bool
	Names() []string
}

// BaseFactory 工厂基础实现
type BaseFactory[T any] struct {
	instances sync.Map
	closers   map[string]func() error
	mu        sync.RWMutex
}

func NewBaseFactory[T any]() *BaseFactory[T] {
	return &BaseFactory[T]{
		closers: make(map[string]func() error),
	}
}

func (f *BaseFactory[T]) Get(name string) (T, error) {
	if instance, ok := f.instances.Load(name); ok {
		return instance.(T), nil
	}
	var zero T
	return zero, NewError(ParamInvalid, "factory: instance %q not found", name)
}

func (f *BaseFactory[T]) Has(name string) bool {
	_, ok := f.instances.Load(name)
	return ok
}

func (f *BaseFactory[T]) Names() []string {
	names := make([]string, 0)
	f.instances.Range(func(key, value interface{}) bool {
		names = append(names, key.(string))
		return true
	})
	return names
}

func (f *BaseFactory[T]) Register(name string, instance T) error {
	if _, loaded := f.instances.LoadOrStore(name, instance); loaded {
		return NewError(DuplicateCall, "factory: instance %q already registered", name)
	}
	return nil
}

func (f *BaseFactory[T]) RegisterWithCloser(name string, instance T, closer func() error) error {
	if err := f.Register(name, instance); err != nil {
		return err
	}
	if closer != nil {
		f.mu.Lock()
		f.closers[name] = closer
		f.mu.Unlock()
	}
	return nil
}

func (f *BaseFactory[T]) Close(name string) error {
	f.mu.RLock()
	closer, ok := f.closers[name]
	f.mu.RUnlock()

	if !ok {
		return NewError(ParamInvalid, "factory: instance %q not found or already closed", name)
	}

	if err := closer(); err != nil {
		return NewErrorWrapped(fmt.Sprintf("factory: fail to close instance %q", name), err)
	}

	f.mu.Lock()
	delete(f.closers, name)
	f.mu.Unlock()
	f.instances.Delete(name)

	return nil
}

func (f *BaseFactory[T]) CloseAll() error {
	var errs []error
	names := f.Names()

	for _, name := range names {
		if err := f.Close(name); err != nil {
			if e, ok := err.(ErrorTao); ok && e.Code() == ParamInvalid {
				continue
			}
			errs = append(errs, NewErrorWrapped(fmt.Sprintf("factory: fail to close instance %q during close all", name), err))
		}
	}

	if len(errs) > 0 {
		return NewError(InstancesNotClosed, "factory: %d errors during close all instances", len(errs))
	}
	return nil
}

func (f *BaseFactory[T]) Count() int {
	count := 0
	f.instances.Range(func(key, value interface{}) bool {
		count++
		return true
	})
	return count
}
