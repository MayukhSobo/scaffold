#!/usr/bin/env python3

"""
install-tools.py - Centralized tool installation script with rich formatting
Usage: python install-tools.py [tool1] [tool2] ...
If no tools specified, installs all tools
"""

import os
import sys
import subprocess
import platform
import shutil
from pathlib import Path
from typing import List, Dict, Optional, Tuple

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

# Import our version helper
scripts_dir = Path(__file__).parent
sys.path.insert(0, str(scripts_dir))
try:
    from version_helper import get_version, load_versions_file
except ImportError:
    # If running through run.py, try different import
    import importlib.util
    spec = importlib.util.spec_from_file_location("version_helper", scripts_dir / "version-helper.py")
    version_helper = importlib.util.module_from_spec(spec)
    spec.loader.exec_module(version_helper)
    get_version = version_helper.get_version
    load_versions_file = version_helper.load_versions_file

# Initialize console
console = Console() if HAS_RICH else None

class ToolInstaller:
    """Tool installation manager with rich output"""
    
    def __init__(self):
        self.versions = load_versions_file()
        self.tools_config = {
            'golangci-lint': {
                'version_key': 'tools.golangci-lint',
                'install_method': self._install_golangci_lint,
                'check_method': self._check_golangci_lint,
                'description': 'Go linter'
            },
            'gotestsum': {
                'version_key': 'tools.gotestsum',
                'install_method': self._install_go_tool,
                'check_method': self._check_go_tool,
                'description': 'Test runner',
                'package': 'gotest.tools/gotestsum'
            },
            'gosec': {
                'version_key': 'tools.gosec',
                'install_method': self._install_go_tool,
                'check_method': self._check_go_tool,
                'description': 'Security analyzer',
                'package': 'github.com/securego/gosec/v2/cmd/gosec'
            },
            'govulncheck': {
                'version_key': 'tools.govulncheck',
                'install_method': self._install_go_tool,
                'check_method': self._check_go_tool,
                'description': 'Vulnerability checker',
                'package': 'golang.org/x/vuln/cmd/govulncheck'
            },
            'air': {
                'version_key': 'tools.air',
                'install_method': self._install_go_tool,
                'check_method': self._check_go_tool,
                'description': 'Live reload',
                'package': 'github.com/air-verse/air'
            },
            'gocov': {
                'version_key': 'tools.gocov',
                'install_method': self._install_go_tool,
                'check_method': self._check_go_tool,
                'description': 'Coverage tool',
                'package': 'github.com/axw/gocov/gocov'
            },
            'gocov-html': {
                'version_key': 'tools.gocov-html',
                'install_method': self._install_go_tool,
                'check_method': self._check_go_tool,
                'description': 'Coverage HTML generator',
                'package': 'github.com/matm/gocov-html/cmd/gocov-html'
            },
            'go-cover-treemap': {
                'version_key': 'tools.go-cover-treemap',
                'install_method': self._install_go_tool,
                'check_method': self._check_go_tool,
                'description': 'Coverage treemap',
                'package': 'github.com/nikolaydubina/go-cover-treemap'
            },
            'trivy': {
                'version_key': 'tools.trivy',
                'install_method': self._install_trivy,
                'check_method': self._check_binary_tool,
                'description': 'Security scanner'
            }
        }
    
    def _log_info(self, message: str):
        """Log info message"""
        if HAS_RICH:
            console.print(f"✅ {message}", style="green")
        else:
            print(f"[INFO] {message}")
    
    def _log_warn(self, message: str):
        """Log warning message"""
        if HAS_RICH:
            console.print(f"⚠️  {message}", style="yellow")
        else:
            print(f"[WARN] {message}")
    
    def _log_error(self, message: str):
        """Log error message"""
        if HAS_RICH:
            console.print(f"❌ {message}", style="red")
        else:
            print(f"[ERROR] {message}")
    
    def _run_command(self, cmd: List[str], description: str = "") -> Tuple[bool, str]:
        """Run a command and return success status and output"""
        try:
            result = subprocess.run(
                cmd, 
                capture_output=True, 
                text=True, 
                check=False
            )
            return result.returncode == 0, result.stdout + result.stderr
        except Exception as e:
            return False, str(e)
    
    def _check_binary_tool(self, tool_name: str) -> bool:
        """Check if a binary tool is available"""
        return shutil.which(tool_name) is not None
    
    def _check_go_tool(self, tool_name: str) -> bool:
        """Check if a Go tool is available"""
        return self._check_binary_tool(tool_name)
    
    def _check_golangci_lint(self, tool_name: str) -> bool:
        """Check golangci-lint version"""
        if not self._check_binary_tool(tool_name):
            return False
        
        success, output = self._run_command([tool_name, "version", "--short"])
        if success:
            expected_version = get_version('tools.golangci-lint')
            return expected_version.lstrip('v') in output
        return True  # If we can't check version, assume it's okay
    
    def _install_go_tool(self, tool_name: str, config: Dict) -> bool:
        """Install a Go tool using go install"""
        version = get_version(config['version_key'])
        package = config['package']
        
        package_with_version = f"{package}@{version}"
        
        if HAS_RICH:
            with console.status(f"Installing {tool_name}@{version}..."):
                success, output = self._run_command(['go', 'install', package_with_version])
        else:
            print(f"Installing {tool_name}@{version}...")
            success, output = self._run_command(['go', 'install', package_with_version])
        
        if success:
            self._log_info(f"{tool_name} installed successfully")
            return True
        else:
            self._log_error(f"Failed to install {tool_name}: {output}")
            return False
    
    def _install_golangci_lint(self, tool_name: str, config: Dict) -> bool:
        """Install golangci-lint"""
        version = get_version(config['version_key'])
        package = f"github.com/golangci/golangci-lint/v2/cmd/golangci-lint@{version}"
        
        if HAS_RICH:
            with console.status(f"Installing {tool_name}@{version}..."):
                success, output = self._run_command(['go', 'install', package])
        else:
            print(f"Installing {tool_name}@{version}...")
            success, output = self._run_command(['go', 'install', package])
        
        if success:
            self._log_info(f"{tool_name} installed successfully")
            return True
        else:
            self._log_error(f"Failed to install {tool_name}: {output}")
            return False
    
    def _install_trivy(self, tool_name: str, config: Dict) -> bool:
        """Install trivy using platform-specific method"""
        system = platform.system().lower()
        
        if HAS_RICH:
            with console.status(f"Installing {tool_name}..."):
                if system == "darwin":
                    success, output = self._install_trivy_macos()
                elif system == "linux":
                    success, output = self._install_trivy_linux()
                else:
                    self._log_error(f"Unsupported OS: {system}. Please install Trivy manually.")
                    return False
        else:
            print(f"Installing {tool_name}...")
            if system == "darwin":
                success, output = self._install_trivy_macos()
            elif system == "linux":
                success, output = self._install_trivy_linux()
            else:
                self._log_error(f"Unsupported OS: {system}. Please install Trivy manually.")
                return False
        
        if success:
            self._log_info(f"{tool_name} installed successfully")
            return True
        else:
            self._log_error(f"Failed to install {tool_name}: {output}")
            return False
    
    def _install_trivy_macos(self) -> Tuple[bool, str]:
        """Install trivy on macOS using Homebrew"""
        if not shutil.which('brew'):
            return False, "Homebrew not found. Please install Homebrew first."
        
        return self._run_command(['brew', 'install', 'aquasecurity/trivy/trivy'])
    
    def _install_trivy_linux(self) -> Tuple[bool, str]:
        """Install trivy on Linux"""
        # Try different package managers
        if shutil.which('apt-get'):
            # Debian/Ubuntu
            commands = [
                ['sudo', 'apt-get', 'update'],
                ['sudo', 'apt-get', 'install', '-y', 'wget', 'apt-transport-https', 'gnupg', 'lsb-release'],
                ['wget', '-qO', '-', 'https://aquasecurity.github.io/trivy-repo/deb/public.key'],
                ['sudo', 'apt-key', 'add', '-'],
                ['sudo', 'apt-get', 'update'],
                ['sudo', 'apt-get', 'install', '-y', 'trivy']
            ]
            
            for cmd in commands:
                success, output = self._run_command(cmd)
                if not success:
                    return False, output
            return True, "Trivy installed via apt-get"
        
        return False, "Unsupported Linux distribution. Please install Trivy manually."
    
    def install_tool(self, tool_name: str) -> bool:
        """Install a specific tool"""
        if tool_name not in self.tools_config:
            self._log_error(f"Unknown tool: {tool_name}")
            return False
        
        config = self.tools_config[tool_name]
        
        # Check if already installed
        if config['check_method'](tool_name):
            self._log_info(f"{tool_name} is already installed and up to date")
            return True
        
        # Install the tool
        return config['install_method'](tool_name, config)
    
    def install_tools(self, tools: List[str] = None) -> bool:
        """Install multiple tools"""
        if not tools:
            tools = list(self.tools_config.keys())
        
        if HAS_RICH:
            # Create a nice panel header
            panel = Panel(
                f"Installing {len(tools)} development tools",
                title="🔧 Tool Installation",
                title_align="left"
            )
            console.print(panel)
            
            # Create a table of tools to install
            table = Table(title="Tools to Install")
            table.add_column("Tool", style="cyan", no_wrap=True)
            table.add_column("Description", style="white")
            table.add_column("Version", style="yellow")
            
            for tool in tools:
                if tool in self.tools_config:
                    config = self.tools_config[tool]
                    version = get_version(config['version_key'])
                    table.add_row(tool, config['description'], version)
            
            console.print(table)
            console.print()
        else:
            print(f"Installing tools: {', '.join(tools)}")
        
        success_count = 0
        total_count = len(tools)
        
        if HAS_RICH:
            with Progress(
                SpinnerColumn(),
                TextColumn("[bold blue]{task.description}"),
                BarColumn(),
                TaskProgressColumn(),
                console=console
            ) as progress:
                task = progress.add_task("Installing tools...", total=total_count)
                
                for tool in tools:
                    progress.update(task, description=f"Installing {tool}...")
                    if self.install_tool(tool):
                        success_count += 1
                    progress.advance(task)
        else:
            for tool in tools:
                if self.install_tool(tool):
                    success_count += 1
        
        # Summary
        if HAS_RICH:
            if success_count == total_count:
                console.print(Panel(
                    f"✅ All {success_count} tools installed successfully!",
                    title="🎉 Installation Complete",
                    style="green"
                ))
            else:
                console.print(Panel(
                    f"⚠️  {success_count}/{total_count} tools installed successfully",
                    title="⚠️  Installation Completed with Issues",
                    style="yellow"
                ))
        else:
            if success_count == total_count:
                print(f"✅ All {success_count} tools installed successfully!")
            else:
                print(f"⚠️  {success_count}/{total_count} tools installed successfully")
        
        return success_count == total_count

def main():
    """Main entry point"""
    if HAS_RICH:
        console.print(Panel(
            "Development Tool Installer",
            subtitle="Using versions from versions.yml",
            style="bold blue"
        ))
    else:
        print("Development Tool Installer")
        print("Using versions from versions.yml")
    
    installer = ToolInstaller()
    
    # Get tools from command line arguments
    tools = sys.argv[1:] if len(sys.argv) > 1 else None
    
    success = installer.install_tools(tools)
    
    if not success:
        sys.exit(1)

if __name__ == '__main__':
    main() 