# Python Scripts & Development Setup

Modern Python-based automation scripts with rich formatting and enhanced functionality.

## 🐍 Overview

All development scripts are written in Python with enhanced features:
- **Enhanced output**: Progress bars, tables, and colored output for better clarity
- **Auto venv activation**: Scripts automatically use the Python virtual environment
- **Enhanced error handling**: Clear error messages and recovery suggestions
- **Cross-platform support**: Works on macOS, Linux, and Windows
- **Graceful fallbacks**: Works without enhanced formatting if dependencies aren't available

## 📦 Python Environment Setup

### Quick Setup

```bash
# Setup Python virtual environment with all dependencies
task python
# or
task shared:setup:python

# Show shell-specific activation instructions
task shared:setup:python:activate
```

### What Gets Installed

The setup process:
1. **Detects and installs Python 3** if missing (macOS/Linux)
2. **Creates `.venv` virtual environment** 
3. **Installs dependencies** from `requirements.txt`:
   - **PyYAML** ≥ 6.0 - YAML parsing
   - **rich** ≥ 13.0.0 - Terminal formatting and progress bars
   - **requests** ≥ 2.31.0 - HTTP requests for downloads and API calls
4. **Adds `.venv/` to .gitignore**
5. **Cross-platform support** (Homebrew, apt, yum, dnf)

### Manual Activation

```bash
# Show shell-specific activation instructions
task shared:setup:python:activate

# Or activate manually based on your shell:
# Bash/Zsh: source .venv/bin/activate
# Fish: source .venv/bin/activate.fish
# C shell: source .venv/bin/activate.csh

# Now scripts use enhanced features
python scripts/version-helper.py list

# Deactivate when done
deactivate
```

## 🔧 Script Architecture

### Smart Script Runner (`scripts/run.py`)

Central script runner that:
- **Automatically detects and uses** the `.venv` Python interpreter
- **Provides helpful error messages** if venv isn't set up
- **Handles script execution** with proper error handling
- **Cross-platform support** for Windows, macOS, and Linux

Usage:
```bash
python scripts/run.py <script_name> [args...]
python scripts/run.py version-helper list
python scripts/run.py install-tools gotestsum
python scripts/run.py sync-versions --dry-run
```

### Available Scripts

| Script | Purpose | Features |
|--------|---------|----------|
| `version-helper.py` | Version management and YAML parsing | Formatted version display |
| `install-tools.py` | Development tool installation | Progress bars, status tables |
| `sync-versions.py` | Version synchronization across files | Change previews, diff tables |
| `install-task.py` | Task runner installation | Platform detection, progress tracking |

## 🎨 Output Features

### Progress Bars & Spinners

```python
# Example from install-tools.py
with Progress(
    SpinnerColumn(),
    TextColumn("[bold blue]{task.description}"),
    BarColumn(),
    TaskProgressColumn(),
    console=console
) as progress:
    task = progress.add_task("Installing tools...", total=total_count)
    # ... installation logic
```

### Tables

```python
# Example from sync-versions.py
table = Table(title="📝 Changes Made")
table.add_column("File", style="cyan", no_wrap=True)
table.add_column("Version Key", style="yellow")
table.add_column("Old Version", style="red")
table.add_column("New Version", style="green")
```

### Informative Panels

```python
# Example from install-tools.py
console.print(Panel(
    f"✅ All {success_count} tools installed successfully!",
    title="🎉 Installation Complete",
    style="green"
))
```

## 🔄 Script Integration

### Task Integration

All scripts are integrated with Task via the `run.py` wrapper:

```yaml
# In tasks/shared.yml
shared:setup:tools:
  cmds:
    - python scripts/run.py install-tools "{{.CLI_ARGS}}"

shared:version:show:
  cmds:
    - python scripts/run.py version-helper list
```

### Automatic venv Usage

Scripts automatically use the virtual environment when:
1. Called through `python scripts/run.py`
2. Virtual environment exists in `.venv/`
3. Required dependencies are installed

### Error Handling

Smart error handling with helpful messages:

```python
def ensure_venv_setup():
    if not check_venv_exists():
        print("❌ Python virtual environment not found!")
        print("Please set it up first by running:")
        print("  task python")
        print("  # or")
        print("  task shared:setup:python")
        sys.exit(1)
```

## 📋 Script Details

### Version Helper (`version-helper.py`)

**Purpose**: Parse and manage versions from `versions.yml`

**Features**:
- Robust YAML parsing with PyYAML
- Fallback parser for environments without PyYAML
- Command-line interface for getting specific versions
- List all versions in organized format

**Usage**:
```bash
python scripts/run.py version-helper list
python scripts/run.py version-helper get go
python scripts/run.py version-helper get tools.golangci-lint
```

### Tool Installer (`install-tools.py`)

**Purpose**: Centralized installation of development tools

**Features**:
- Progress bars during installation
- Tool status checking (already installed vs needs update)
- Version-aware installation from `versions.yml`
- Support for selective tool installation
- Cross-platform tool installation methods

**Usage**:
```bash
# Install all tools
task shared:setup:tools

# Install specific tools
task shared:setup:tools -- gotestsum golangci-lint
```

### Version Synchronizer (`sync-versions.py`)

**Purpose**: Sync versions from `versions.yml` across all project files

**Features**:
- Tables showing what will be changed
- Dry-run mode for safe previewing
- Consistency checking across files
- Detailed change reporting
- File-by-file synchronization status

**Usage**:
```bash
task shared:version:sync:dry-run  # Preview changes
task shared:version:sync          # Apply changes
task shared:version:check         # Check consistency
```

### Task Installer (`install-task.py`)

**Purpose**: Install the Task runner tool

**Features**:
- Cross-platform detection (macOS, Linux, Windows)
- Multiple installation methods (package managers, direct download)
- Version checking and updates
- Progress tracking for downloads
- Installation verification

## 🛡️ Safety & Reliability

### Fallback Support

Scripts work even without enhanced formatting dependencies:

```python
try:
    from rich.console import Console
    from rich.progress import Progress
    HAS_RICH = True
except ImportError:
    HAS_RICH = False

# Later in code:
if HAS_RICH:
    console.print(f"✅ {message}", style="green")
else:
    print(f"[INFO] {message}")
```

### Error Recovery

- **Clear error messages** with suggested solutions
- **Graceful degradation** when dependencies missing
- **Cross-platform compatibility** checks
- **Network error handling** for downloads
- **File permission checks** for installations

### Idempotent Operations

All scripts are designed to be run multiple times safely:
- **Tool installation**: Checks if already installed
- **venv setup**: Skips if already exists
- **Version sync**: Only changes what's needed

## 🎯 Best Practices

1. **Always use via Task**: Scripts designed for Task integration
2. **Setup venv first**: Run `task python` for best experience
3. **Install dependencies**: Install dependencies for enhanced output
4. **Cross-platform**: Scripts work on macOS, Linux, and Windows
5. **Error handling**: Scripts provide clear guidance on issues
6. **Version management**: Let scripts handle tool versions automatically 