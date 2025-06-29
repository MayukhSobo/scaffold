<div align="center">
  <img src="scaffold.png" alt="Scaffold Banner" width="400" height="400"/>
</div>

# 🚀 Scaffold: High-Performance Go Application Boilerplate

<div align="center">

[![Go Version](https://img.shields.io/badge/go-1.22+-blue.svg)](https://golang.org/dl/)
[![Build Status](https://img.shields.io/badge/build-passing-brightgreen.svg)](https://github.com/MayukhSobo/scaffold)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/MayukhSobo/scaffold)](https://goreportcard.com/report/github.com/MayukhSobo/scaffold)
[![Last Commit](https://img.shields.io/github/last-commit/MayukhSobo/scaffold.svg)](https://github.com/MayukhSobo/scaffold/commits/main)

A production-ready Go application boilerplate, engineered for performance, developer efficiency, and robust tooling. Featuring a modular build system with Task, advanced testing & coverage, hot-reloading, and optimized cross-platform builds.

</div>

---

## ✨ Features

- **🔩 Modular Build System**: Powered by a modular `Taskfile` structure for streamlined, organized, and maintainable builds and development tasks.
- **⚡ Hot-Reloading**: Uses `air` for live-reloading during development, boosting productivity.
- **🧪 Comprehensive Testing**: Integrated with `gotestsum` for beautiful, readable test outputs. Supports unit, integration, and benchmark tests.
- **📊 Advanced Code Coverage**:
  - **Standard HTML reports**.
  - **Enhanced function-level reports** via `gocov` and `gocov-html`.
  - **Interactive visual treemaps** via `go-cover-treemap`.
- **🏆 Code Quality Assurance**:
  - **Linting** with `golangci-lint` using a comprehensive ruleset with smart version management.
  - **Formatting** with `gofmt`.
  - **Static analysis** with `go vet`.
- **⚙️ Configuration Management**: Flexible configuration loading for different environments using `viper`.
- **🚀 Optimized Production Builds**:
  - **Cross-platform compilation** for Linux, macOS, and Windows.
  - **Aggressive size reduction** using `ldflags` (`-s -w`).
  - **UPX compression** for ultra-small binaries (up to 85% size reduction).
- **🐳 Docker Ready**: Multi-stage `Dockerfile` for small, secure production images.
- **🎛️ Centralized Binary Naming**: Easily manage binary names from a single variable in the `Taskfile`.
- **📖 Self-Documenting**: Includes a `task help:help` command for a detailed overview of all available tasks.

---

## 🛠️ Key Dependencies

### Core Libraries

| Library                               | Description                               |
| ------------------------------------- | ----------------------------------------- |
| [`github.com/gin-gonic/gin`](https://github.com/gin-gonic/gin) | High-performance HTTP web framework. |
| [`github.com/spf13/viper`](https://github.com/spf13/viper) | Complete configuration solution. |
| [`github.com/rs/zerolog`](https://github.com/rs/zerolog) | Blazing fast, structured JSON logger. |
| [`gorm.io/gorm`](https://github.com/go-gorm/gorm) | The fantastic ORM library for Go. |
| [`gopkg.in/natefinch/lumberjack.v2`](https://github.com/natefinch/lumberjack) | Log rotation for file-based logging. |

### Development & Tooling

| Tool | Description |
|---|---|
| [`github.com/air-verse/air`](https://github.com/air-verse/air) | Live-reloading for Go applications. |
| [`gotest.tools/gotestsum`](https://github.com/gotestyourself/gotestsum) | 'go test' runner with custom output formatting. |
| [`github.com/golangci/golangci-lint`](https://github.com/golangci/golangci-lint) | Fast Go linters runner. |
| [`github.com/axw/gocov`](https://github.com/axw/gocov) | Coverage reporting tool. |
| [`github.com/matm/gocov-html`](https://github.com/matm/gocov-html) | Generates HTML reports from `gocov` data. |
| [`github.com/nikolaydubina/go-cover-treemap`](https://github.com/nikolaydubina/go-cover-treemap)| Generates visual treemaps for coverage. |

---

## 🏗️ Project Structure

```
scaffold/
├── build/                   # Build artifacts (binaries)
├── cmd/
│   └── server/              # Main application entrypoint
├── configs/                 # Configuration files (local.yml, prod.yml)
├── internal/                # Internal application code
│   ├── handler/
│   ├── middleware/
│   ├── model/
│   ├── repository/
│   └── service/
├── pkg/                     # Public packages
│   ├── config/
│   ├── helper/
│   ├── http/
│   └── log/
├── tasks/                   # Modular task definitions
│   ├── build.yml            # Build-related tasks
│   ├── clean.yml            # Cleanup tasks
│   ├── config.yml           # Configuration validation
│   ├── deps.yml             # Dependency management
│   ├── dev.yml              # Development tasks
│   ├── docker.yml           # Docker operations
│   ├── help.yml             # Help documentation
│   ├── quality.yml          # Code quality (lint, fmt, vet)
│   └── test.yml             # Testing tasks
├── .air.toml                # Configuration for hot-reloading (air)
├── .golangci.yml            # Configuration for golangci-lint
├── Dockerfile               # Multi-stage Dockerfile
├── go.mod
├── go.sum
└── Taskfile.yml             # Main task configuration with includes
```

---

## 🏁 Getting Started

### Prerequisites

- [Go](https://golang.org/dl/) (version 1.22+ recommended)
- [Task](https://taskfile.dev/installation/)
- [Docker](https://www.docker.com/get-started) (for containerized builds)
- [UPX](https://upx.github.io/) (optional, for binary compression)

### Installation

1.  **Clone the repository:**
    ```bash
    git clone https://github.com/thedatageek/scaffold.git
    cd scaffold
    ```

2.  **Install dependencies:**
    The project uses Go Modules. The required tools and dependencies are installed automatically when you run a task for the first time. To install them manually:
    ```bash
    task deps:deps
    ```

---

## 🚀 Usage: Available Tasks

This project uses a **modular `Taskfile.yml` structure** as a modern alternative to `Makefile`. All commands are managed through `task` with organized namespaces.

Run `task --list` for a quick overview or `task help:help` for detailed descriptions.

### 📦 Build Tasks

| Command                               | Description                                                     |
| ------------------------------------- | --------------------------------------------------------------- |
| `task build:build`                    | Build a development binary with debug symbols and race detection. |
| `task build:build:release:linux`      | Build an optimized, compressed release binary for Linux.        |
| `task build:build:release:darwin`     | Build an optimized, compressed release binary for macOS.        |
| `task build:build:release:windows`    | Build an optimized, compressed release binary for Windows.      |
| `task build:build:release:all`        | Build release binaries for all platforms.                       |

### 🏃 Development Tasks

| Command                 | Description                                         |
| ----------------------- | --------------------------------------------------- |
| `task dev:run`          | Run the application using `configs/local.yml`.      |
| `task dev:run:prod`     | Run the application using `configs/prod.yml`.       |
| `task dev:dev`          | Run with **hot-reloading** using `air`.                 |

### 🧪 Test Tasks

| Command                           | Description                                                                  |
| --------------------------------- | ---------------------------------------------------------------------------- |
| `task test:test`                  | Run all tests using `gotestsum`.                                             |
| `task test:test:unit`             | Run only unit tests.                                                         |
| `task test:test:integration`      | Run only integration tests.                                                  |
| `task test:test:race`             | Run tests with the race detector enabled.                                    |
| `task test:test:benchmark`        | Run benchmark tests.                                                         |
| `task test:test:coverage`         | Generate a standard HTML coverage report.                                    |
| `task test:test:coverage:open`    | Generate **comprehensive coverage reports** (HTML, gocov, treemap) and open in browser. |

### 🔍 Code Quality Tasks

| Command                    | Description                                  |
| -------------------------- | -------------------------------------------- |
| `task quality:lint`        | Run `golangci-lint` to find code issues.     |
| `task quality:lint skip=true` | Run linter without checking/installing golangci-lint. |
| `task quality:lint force=true`| Force reinstall golangci-lint and run linter. |
| `task quality:fmt`         | Format all Go source files with `gofmt`.     |
| `task quality:vet`         | Run `go vet` to analyze source code.         |
| `task quality:check`       | Run all quality checks (`fmt`, `vet`, `lint`). |

### 🧹 Cleanup Tasks

| Command                               | Description                                  |
| ------------------------------------- | -------------------------------------------- |
| `task clean:clean`                    | Clean all build artifacts and caches.        |
| `task clean:clean:debug`              | Clean development build artifacts only.      |
| `task clean:clean:release:linux`      | Clean Linux release build artifacts only.    |
| `task clean:clean:release:darwin`     | Clean macOS release build artifacts only.    |
| `task clean:clean:release:windows`    | Clean Windows release build artifacts only.  |
| `task clean:clean:release:all`        | Clean all release build artifacts.           |

### 📦 Dependency Management

| Command                    | Description                                  |
| -------------------------- | -------------------------------------------- |
| `task deps:deps`           | Download and tidy Go module dependencies.    |
| `task deps:deps:install`   | Install/update dependencies.                 |
| `task deps:deps:update`    | Update all dependencies to the latest versions. |

### 🐳 Docker Tasks

| Command                     | Description                                  |
| --------------------------- | -------------------------------------------- |
| `task docker:docker:build` | Build a production-ready Docker image.       |
| `task docker:docker:run`   | Run the application in a Docker container.   |

### ⚙️ Configuration & Help

| Command                           | Description                                  |
| --------------------------------- | -------------------------------------------- |
| `task config:config:validate`    | Validate all `.yml` files in the `configs` directory. |
| `task help:help`                 | Show detailed help for all tasks.            |

---

## 🔧 Modular Build System

This project features a **modular Taskfile structure** that organizes tasks into logical namespaces:

### 📁 Task Organization

```
Taskfile.yml              # Main configuration with includes
├── tasks/build.yml        # Build operations
├── tasks/clean.yml        # Cleanup operations  
├── tasks/config.yml       # Configuration validation
├── tasks/deps.yml         # Dependency management
├── tasks/dev.yml          # Development workflow
├── tasks/docker.yml       # Container operations
├── tasks/help.yml         # Documentation
├── tasks/quality.yml      # Code quality assurance
└── tasks/test.yml         # Testing operations
```

### 🎯 Benefits

- **Modularity**: Each file focuses on a specific domain
- **Maintainability**: Easier to find and modify specific tasks
- **Collaboration**: Team members can work on different task files simultaneously
- **Reusability**: Individual task files can be shared across projects

---

## ⚙️ Configuration

Application configuration is managed by `viper` and loaded from the `configs/` directory.

-   **`configs/local.yml`**: Used for local development (`task dev:run`, `task dev:dev`).
-   **`configs/prod.yml`**: Used for production runs (`task dev:run:prod`).

You can specify a configuration file using the `--config` flag:
```bash
go run ./cmd/server --config=configs/local.yml
```

The system also supports a `--validate-config` flag to check if a configuration file is valid without running the server, used in the `task config:config:validate` task.

---

## 📦 Build & Deployment

### Development Build

For a quick debug build with race detection enabled:
```bash
task build:build
```
This creates a binary at `build/debug/scaffold`.

### Production Release Builds

To create highly optimized and compressed binaries for distribution:
```bash
task build:build:release:all
```
This generates binaries for Linux, macOS, and Windows in their respective `build/` subdirectories (e.g., `build/linux/scaffold-amd64-linux`).

**Optimization Highlights:**
- **Stripped Symbols (`-s -w`):** Removes debug information to reduce size.
- **Static Linking:** Creates self-contained binaries where possible.
- **UPX Compression:** Further compresses the binary, often resulting in an **80-85% size reduction**. A 10MB binary can become ~1.5MB.

---

## 🔬 Testing and Coverage

This boilerplate offers a rich testing and coverage experience.

To run all tests:
```bash
task test:test
```

To generate and view the full suite of coverage reports:
```bash
task test:test:coverage:open
```
This command:
1.  Runs tests and generates coverage data.
2.  Creates three different reports in the `reports/` directory:
    - `coverage.html` (standard)
    - `coverage-enhanced.html` (detailed)
    - `coverage-treemap.svg` (visual)
3.  Starts a local web server on port `8080`.
4.  Opens your browser to view the reports.

---

## 🔍 Advanced Linting

The project includes smart `golangci-lint` management with version control:

```bash
# Normal linting (auto-installs if needed)
task quality:lint

# Skip installation check (faster if you know it's installed)
task quality:lint skip=true

# Force reinstall golangci-lint
task quality:lint force=true
```

The linter automatically:
- Checks if the correct version is installed
- Installs `golangci-lint` v2.2.0 if missing or outdated
- Uses your project's `.golangci.yml` configuration

---

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1.  Fork the repository.
2.  Create your feature branch (`git checkout -b feature/AmazingFeature`).
3.  Commit your changes (`git commit -m 'Add some AmazingFeature'`).
4.  Push to the branch (`git push origin feature/AmazingFeature`).
5.  Open a Pull Request.

---

## 📜 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details. 
