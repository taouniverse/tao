# tao

[![MainTest](https://github.com/taouniverse/tao/actions/workflows/main.yml/badge.svg)](https://github.com/taouniverse/tao/actions/workflows/main.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/taouniverse/tao)](https://goreportcard.com/report/github.com/taouniverse/tao)
[![codecov](https://codecov.io/gh/taouniverse/tao/branch/main/graph/badge.svg?token=NC37IY2D8S)](https://codecov.io/gh/taouniverse/tao)
[![GoDoc](https://pkg.go.dev/badge/github.com/taouniverse/tao?status.svg)](https://pkg.go.dev/github.com/taouniverse/tao?tab=doc)
![GitHub release (latest by date)](https://img.shields.io/github/v/release/taouniverse/tao?style=flat-square)

```go
___________              
\__    ___/____    ____  
  |    |  \__  \  /  _ \ 
  |    |   / __ \(  <_> )
  |____|  (____  /\____/ 
               \/
```

> The Tao produced One; One produced Two; Two produced Three; Three produced All things.

**Tao** is a lightweight, modular Go application framework designed for building scalable and maintainable applications. It provides a pipeline-based task scheduling system, factory-pattern component management, configuration management, and graceful shutdown capabilities.

## Features

- 🏭 **Generic Factory Pattern** - Type-safe component factories with multi-instance support via `BaseFactory[T]`
- ⚙️ **Flexible Configuration** - Support for YAML/JSON config files, auto-detecting single/multi-instance modes
- 📝 **Structured Logging** - Built-in multi-level logging system with customizable outputs
- 🔄 **Graceful Shutdown** - Proper resource cleanup on application termination
- 🧩 **Modular Design** - Easy to extend with custom components using `Register[T, C]()`
- 🎯 **Context Aware** - Full support for Go's context package for cancellation and timeouts

## Installation

```bash
go get github.com/taouniverse/tao
```

## Quick Start

### 1. Create a Configuration File

Create `conf/config.yaml`:

```yaml
tao:
  log:
    level: debug
    type: console
    flag: std|short
    call_depth: 3
    disable: false
  banner:
    hide: false
```

### 2. Define Your Component (with Factory Pattern)

```go
package main

import (
    "context"
    "github.com/taouniverse/tao"
)

// InstanceConfig stores per-instance fields
type MyInstanceConfig struct {
    Message string `json:"message"`
}

// Config implements tao.MultiConfig for multi-instance support
type MyConfig struct {
    tao.BaseMultiConfig[MyInstanceConfig]
}

func (c *MyConfig) Name() string      { return "myapp" }
func (c *MyConfig) ValidSelf()        {}
func (c *MyConfig) ToTask() tao.Task  { return nil }
func (c *MyConfig) RunAfter() []string { return nil }

var Factory *tao.BaseFactory[string]

func init() {
    var err error
    Factory, err = tao.Register("myapp", &MyConfig{}, NewMyApp)
    if err != nil {
        panic(err)
    }
}

func NewMyApp(name string, cfg MyInstanceConfig) (string, func() error, error) {
    closer := func() error { return nil }
    return cfg.Message, closer, nil
}

func main() {
    if err := tao.Run(context.Background(), nil); err != nil {
        panic(err)
    }
}
```

### 3. Run Your Application

```bash
go run main.go
```

## Core Concepts

### Register — Generic Component Registration

`Register[T any, C any]()` is the unified entry point for all components:

```go
// Multi-instance mode (config implements MultiConfig[C])
factory, _ := tao.Register[*gorm.DB, InstanceConfig]("mysql", &Config{}, NewMySQL)

// Single-instance mode (config does NOT implement MultiConfig)
factory, _ := tao.Register[*gin.Engine, Config]("gin", &Config{}, NewGin)

// Pure config registration (constructor = nil)
_, _ := tao.Register[struct{}, struct{}]("mytask", &TaskConfig{}, nil)
```

### BaseFactory[T] — Type-Safe Component Factory

```go
factory := tao.NewBaseFactory[*gorm.DB]()
factory.RegisterWithCloser("master", dbMaster, dbMaster.Close)
factory.RegisterWithCloser("replica", dbReplica, dbReplica.Close)

db, _ := factory.Get("master")
names := factory.Names()
count := factory.Count()

factory.CloseAll()
```

### MultiConfig[C] — Multi-Instance Configuration Interface

```go
type MultiConfig[C any] interface {
    tao.Config
    GetInstances() []Instance[C]
    SetInstances(instances []Instance[C])
    GetDefaultInstanceName() string
}
```

Embed `BaseMultiConfig[C]` to get the implementation for free:

```go
type Config struct {
    tao.BaseMultiConfig[InstanceConfig]
    RunAfters []string `json:"run_after,omitempty"`
}
```

Instances are stored as an ordered slice to preserve creation order:

```go
type Instance[C any] struct {
    Name string
    Cfg  C
}
```

Use `GetInstanceByName(name string) (C, bool)` to look up a specific instance by name.

### Task

The basic unit of work in Tao:

```go
task := tao.NewTask("mytask", func(ctx context.Context, param tao.Parameter) (tao.Parameter, error) {
    // Your logic here
    return param, nil
}, tao.SetClose(func() error {
    // Cleanup logic
    return nil
}))
```

### Pipeline

A `Pipeline` orchestrates multiple tasks with dependencies:

```go
pipeline := tao.NewPipeline("mypipeline")
pipeline.Register(tao.NewPipeTask(task1))
pipeline.Register(tao.NewPipeTask(task2, "task1")) // task2 depends on task1
```

### Configuration

Implement the `Config` interface for your module:

```go
type Config interface {
    Name() string           // Config identifier
    ValidSelf()             // Set default values
    ToTask() Task           // Convert to executable task
    RunAfter() []string     // Dependency task names
}
```

For multi-instance components, implement `MultiConfig[C]` instead.

### Logging

Tao provides multiple logging levels:

```go
tao.Debug("debug message")
tao.Info("info message")
tao.Warn("warning message")
tao.Error("error message")
```

## Project Structure

```
tao/
├── tao.go          # Core Universe, Run(), Register[T,C]()
├── factory.go      # Factory[T] interface + BaseFactory[T]
├── config.go       # Config/MultiConfig interfaces + parseMultiConfig
├── pipeline.go     # Pipeline implementation for task orchestration
├── task.go         # Task interface and implementation
├── init.go         # Initialization and config loading
├── log.go          # Logging system
├── param.go        # Parameter passing between tasks
├── error.go        # Error handling utilities
├── factory_test.go # Factory & BaseMultiConfig tests
└── config_test.go  # Config & parseMultiConfig tests
```

## Unit Tests

The core library has comprehensive test coverage:

| Test File | Coverage |
|-----------|----------|
| `factory_test.go` | `BaseFactory[T]` — CRUD, close, concurrent access, benchmarks (14 tests) |
| `config_test.go` | `BaseMultiConfig[C]`, `parseMultiConfig` — single/multi-instance parsing, nested structs, edge cases (16 tests) |
| `init_test.go` | Pre-init registration, config path handling |
| `tao_test.go` | Register, Run lifecycle |
| `log_test.go` | LogLevel, LogType, LogFlag marshaling |
| `pipeline_test.go` | Pipeline orchestration, dependencies, errors |
| `task_test.go` | Task lifecycle, result propagation |

### Running Tests

```bash
# All tests
cd tao && go test -v ./...

# Specific category
cd tao && go test -v -run "TestFactory" ./...
cd tao && go test -v -run "TestParseMultiConfig" ./...

# With coverage
cd tao && go test -race -coverprofile=coverage.out ./...
```

## Advanced Usage

### Multi-Instance Database Pattern

```yaml
mysql:
  default_instance: master
  master:
    host: 10.0.0.1
    port: 3306
    user: app
    password: secret
    db: production
  replica:
    host: 10.0.0.2
    port: 3306
    user: readonly
    db: production
```

```go
masterDB, _ := mysql.GetDB("master")
replicaDB, _ := mysql.GetDB("replica")
```

### Custom Logger

```go
writer := os.Stdout
logger := &myCustomLogger{}
tao.SetWriter("myapp", writer)
tao.SetLogger("myapp", logger)
```

### Development Mode

For quick testing without configuration files:

```go
func init() {
    tao.DevelopMode() // Uses default configurations
}
```

### Programmatic Configuration

```go
configData := []byte(`
tao:
  log:
    level: info
    type: console
`)
tao.SetAllConfigBytes(configData, tao.Yaml)
```

## Ecosystem Components

| Component | Description | Factory Type |
|-----------|-------------|-------------|
| [tao-mysql](https://github.com/taouniverse/tao-mysql) | MySQL via GORM | `*BaseFactory[*gorm.DB]` |
| [tao-postgres](https://github.com/taouniverse/tao-postgres) | PostgreSQL via GORM | `*BaseFactory[*gorm.DB]` |
| [tao-sqlite](https://github.com/taouniverse/tao-sqlite) | SQLite via GORM (pure Go) | `*BaseFactory[*gorm.DB]` |
| [tao-redis](https://github.com/taouniverse/tao-redis) | Redis via go-redis | `*BaseFactory[redis.UniversalClient]` |
| [tao-mongodb](https://github.com/taouniverse/tao-mongodb) | MongoDB via mongo-driver | `*BaseFactory[*mongo.Client]` |
| [tao-gin](https://github.com/taouniverse/tao-gin) | Gin HTTP framework | `*BaseFactory[*gin.Engine]` |
| [tao-websocket](https://github.com/taouniverse/tao-websocket) | WebSocket support | `*BaseFactory[ws.Upgrader]` |
| [tao-zap](https://github.com/taouniverse/tao-zap) | Zap structured logging | `*BaseFactory[*zap.SugaredLogger]` |
| [tao-hello](https://github.com/taouniverse/tao-hello) | Example component template | `*BaseFactory[struct{}]` |

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.
