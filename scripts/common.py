#!/usr/bin/env python3

"""
common.py - Shared utilities for all Python scripts
Consolidates common patterns like logging, rich formatting, command execution, etc.
"""

import os
import sys
import subprocess
import importlib.util
from pathlib import Path
from typing import List, Dict, Optional, Tuple, Any
import shutil

# Rich imports with fallback
try:
    from rich.console import Console
    from rich.progress import Progress, SpinnerColumn, TextColumn, BarColumn, TaskProgressColumn
    from rich.panel import Panel
    from rich.table import Table
    from rich.text import Text
    from rich import print as rprint
    HAS_RICH = True
except ImportError:
    HAS_RICH = False

# Requests import with fallback
try:
    import requests
    HAS_REQUESTS = True
except ImportError:
    HAS_REQUESTS = False

# Initialize console
console = Console() if HAS_RICH else None

class Logger:
    """Unified logging interface with rich formatting"""
    
    def __init__(self, name: str = "Script"):
        self.name = name
        self.errors = []
    
    def info(self, message: str):
        """Log info message"""
        if HAS_RICH:
            console.print(f"✅ {message}", style="green")
        else:
            print(f"[INFO] {message}")
    
    def warn(self, message: str):
        """Log warning message"""
        if HAS_RICH:
            console.print(f"⚠️  {message}", style="yellow")
        else:
            print(f"[WARN] {message}")
    
    def error(self, message: str):
        """Log error message"""
        if HAS_RICH:
            console.print(f"❌ {message}", style="red")
        else:
            print(f"[ERROR] {message}")
        self.errors.append(message)
    
    def verbose(self, message: str, enabled: bool = True):
        """Log verbose message"""
        if enabled:
            if HAS_RICH:
                console.print(f"ℹ️  {message}", style="blue")
            else:
                print(f"[VERBOSE] {message}")
    
    def success(self, message: str):
        """Log success message"""
        if HAS_RICH:
            console.print(f"🎉 {message}", style="bold green")
        else:
            print(f"[SUCCESS] {message}")

class CommandRunner:
    """Unified command execution with logging"""
    
    def __init__(self, logger: Logger):
        self.logger = logger
    
    def run(self, cmd: List[str], description: str = "", capture_output: bool = True) -> Tuple[bool, str]:
        """Run a command and return success status and output"""
        try:
            if description:
                self.logger.verbose(f"Running: {description}")
            
            result = subprocess.run(
                cmd, 
                capture_output=capture_output, 
                text=True, 
                check=False
            )
            
            output = result.stdout + result.stderr if capture_output else ""
            return result.returncode == 0, output
        except Exception as e:
            self.logger.error(f"Command execution failed: {e}")
            return False, str(e)
    
    def run_with_status(self, cmd: List[str], status_msg: str) -> Tuple[bool, str]:
        """Run a command within a rich status context"""
        if not has_rich():
            self.logger.info(status_msg)
            return self.run(cmd)

        console = get_console()
        success, output = self.run(cmd)
        if success:
            console.print(f"✅ {status_msg} [green]Success[/green]")
        else:
            console.print(f"❌ {status_msg} [red]Failed[/red]")
        return success, output

class FileManager:
    """Unified file operations with error handling"""
    
    def __init__(self, logger: Logger, dry_run: bool = False):
        self.logger = logger
        self.dry_run = dry_run
    
    def read_file(self, file_path: Path) -> Optional[str]:
        """Read file content with error handling"""
        try:
            return file_path.read_text(encoding='utf-8')
        except Exception as e:
            self.logger.error(f"Failed to read {file_path}: {e}")
            return None
    
    def write_file(self, file_path: Path, content: str) -> bool:
        """Write file content with error handling"""
        if self.dry_run:
            self.logger.verbose(f"Would write {file_path}")
            return True
        
        try:
            # Ensure parent directory exists
            file_path.parent.mkdir(parents=True, exist_ok=True)
            file_path.write_text(content, encoding='utf-8')
            self.logger.verbose(f"Updated {file_path}")
            return True
        except Exception as e:
            self.logger.error(f"Failed to write {file_path}: {e}")
            return False
    
    def file_exists(self, file_path: Path) -> bool:
        """Check if file exists"""
        return file_path.exists()

class VersionHelper:
    """Unified version helper with fallback import logic"""
    
    def __init__(self, logger: Logger):
        self.logger = logger
        self._version_module = None
        self._load_version_module()
    
    def _load_version_module(self):
        """Load version helper module with fallback logic"""
        scripts_dir = Path(__file__).parent
        
        # Try direct import first
        try:
            sys.path.insert(0, str(scripts_dir))
            from version_helper import get_version, load_versions_file, list_versions
            self._version_module = sys.modules['version_helper']
            return
        except ImportError:
            pass
        
        # Try importing from version-helper.py file
        try:
            spec = importlib.util.spec_from_file_location("version_helper", scripts_dir / "version-helper.py")
            version_helper = importlib.util.module_from_spec(spec)
            spec.loader.exec_module(version_helper)
            self._version_module = version_helper
            return
        except Exception as e:
            self.logger.error(f"Failed to load version helper: {e}")
            self._version_module = None
    
    def get_version(self, key: str) -> Optional[str]:
        """Get version by key"""
        if self._version_module and hasattr(self._version_module, 'get_version'):
            try:
                return self._version_module.get_version(key)
            except Exception as e:
                self.logger.error(f"Failed to get version for {key}: {e}")
        return None
    
    def load_versions_file(self) -> Dict[str, Any]:
        """Load versions.yml file"""
        if self._version_module and hasattr(self._version_module, 'load_versions_file'):
            try:
                return self._version_module.load_versions_file()
            except Exception as e:
                self.logger.error(f"Failed to load versions file: {e}")
        return {}
    
    def list_versions(self) -> List[Tuple[str, str]]:
        """List all versions"""
        if self._version_module and hasattr(self._version_module, 'list_versions'):
            try:
                return self._version_module.list_versions()
            except Exception as e:
                self.logger.error(f"Failed to list versions: {e}")
        return []

class RichHelper:
    """Rich formatting utilities"""
    
    @staticmethod
    def create_panel(content: str, title: str = "", style: str = ""):
        """Create a rich panel if available"""
        if HAS_RICH:
            return Panel(content, title=title, style=style)
        return None
    
    @staticmethod
    def create_table(title: str = "") -> Optional[Table]:
        """Create a rich table if available"""
        if HAS_RICH:
            table = Table(title=title)
            return table
        return None
    
    @staticmethod
    def print_panel(content: str, title: str = "", style: str = ""):
        """Print a rich panel or fallback to plain text"""
        if HAS_RICH:
            console.print(Panel(content, title=title, style=style))
        else:
            print(f"=== {title} ===")
            print(content)
            print("=" * (len(title) + 8))
    
    @staticmethod
    def print_table(table: Table):
        """Print a rich table if available"""
        if HAS_RICH and table:
            console.print(table)
    
    @staticmethod
    def status_context(message: str):
        """Create a status context manager"""
        if HAS_RICH:
            return console.status(message)
        else:
            class PlainStatus:
                def __init__(self, msg):
                    self.msg = msg
                def __enter__(self):
                    print(self.msg)
                    return self
                def __exit__(self, *args):
                    pass
            return PlainStatus(message)

class ScriptBase:
    """Base class for all scripts with common functionality"""
    
    def __init__(self, name: str, dry_run: bool = False):
        self.name = name
        self.dry_run = dry_run
        self.logger = Logger(name)
        self.cmd_runner = CommandRunner(self.logger)
        self.file_manager = FileManager(self.logger, dry_run)
        self.version_helper = VersionHelper(self.logger)
        self.rich = RichHelper()
        
        # Set up project root
        self.project_root = Path(__file__).parent.parent
        self.scripts_dir = Path(__file__).parent
    
    def check_binary_exists(self, binary_name: str) -> bool:
        """Check if a binary exists in PATH"""
        return shutil.which(binary_name) is not None
    
    def get_project_root(self) -> Path:
        """Get project root directory"""
        return self.project_root
    
    def get_scripts_dir(self) -> Path:
        """Get scripts directory"""
        return self.scripts_dir
    
    def has_errors(self) -> bool:
        """Check if any errors occurred"""
        return len(self.logger.errors) > 0
    
    def get_errors(self) -> List[str]:
        """Get all errors that occurred"""
        return self.logger.errors.copy()
    
    def exit_with_error(self, message: str = ""):
        """Exit with error status"""
        if message:
            self.logger.error(message)
        sys.exit(1)
    
    def exit_with_success(self, message: str = ""):
        """Exit with success status"""
        if message:
            self.logger.success(message)
        sys.exit(0)

# Global convenience functions for backward compatibility
def get_console() -> Optional[Console]:
    """Get rich console instance"""
    return console

def has_rich() -> bool:
    """Check if rich is available"""
    return HAS_RICH

def has_requests() -> bool:
    """Check if requests is available"""
    return HAS_REQUESTS 