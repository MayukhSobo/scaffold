#!/usr/bin/env python3

"""
install-tools.py - Centralized tool installation script with rich formatting
Usage: python install-tools.py [tool1] [tool2] ...
If no tools specified, installs all tools
"""

import platform
import shutil
from pathlib import Path
from typing import List, Dict, Optional, Tuple

# Import common utilities
from common import ScriptBase, has_rich, get_console

class ToolInstaller(ScriptBase):
    """Tool installation manager with rich output"""
    
    def __init__(self):
        super().__init__("ToolInstaller")
        self.versions = self.version_helper.load_versions_file()
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
    
    def _check_binary_tool(self, tool_name: str) -> bool:
        """Check if a binary tool is available"""
        return self.check_binary_exists(tool_name)
    
    def _check_go_tool(self, tool_name: str) -> bool:
        """Check if a Go tool is available"""
        return self._check_binary_tool(tool_name)
    
    def _check_golangci_lint(self, tool_name: str) -> bool:
        """Check golangci-lint version"""
        if not self._check_binary_tool(tool_name):
            return False
        
        success, output = self.cmd_runner.run([tool_name, "version", "--short"])
        if success:
            expected_version = self.version_helper.get_version('tools.golangci-lint')
            return expected_version.lstrip('v') in output
        return True  # If we can't check version, assume it's okay
    
    def _install_go_tool(self, tool_name: str, config: Dict) -> bool:
        """Install a Go tool using go install"""
        version = self.version_helper.get_version(config['version_key'])
        package = config['package']
        
        package_with_version = f"{package}@{version}"
        
        success, output = self.cmd_runner.run_with_status(
            ['go', 'install', package_with_version],
            f"Installing {tool_name}@{version}..."
        )
        
        if success:
            self.logger.info(f"{tool_name} installed successfully")
            return True
        else:
            self.logger.error(f"Failed to install {tool_name}: {output}")
            return False
    
    def _install_golangci_lint(self, tool_name: str, config: Dict) -> bool:
        """Install golangci-lint"""
        version = self.version_helper.get_version(config['version_key'])
        package = f"github.com/golangci/golangci-lint/v2/cmd/golangci-lint@{version}"
        
        success, output = self.cmd_runner.run_with_status(
            ['go', 'install', package],
            f"Installing {tool_name}@{version}..."
        )
        
        if success:
            self.logger.info(f"{tool_name} installed successfully")
            return True
        else:
            self.logger.error(f"Failed to install {tool_name}: {output}")
            return False
    
    def _install_trivy(self, tool_name: str, config: Dict) -> bool:
        """Install trivy using platform-specific method"""
        system = platform.system().lower()
        
        with self.rich.status_context(f"Installing {tool_name}..."):
            if system == "darwin":
                success = self._install_trivy_macos()
            elif system == "linux":
                success = self._install_trivy_linux()
            else:
                self.logger.error(f"Unsupported OS: {system}. Please install Trivy manually.")
                return False
        
        if success:
            self.logger.info(f"{tool_name} installed successfully")
            return True
        else:
            self.logger.error(f"Failed to install {tool_name}")
            return False
    
    def _install_trivy_macos(self) -> bool:
        """Install trivy on macOS using Homebrew"""
        if not self.check_binary_exists('brew'):
            self.logger.error("Homebrew not found. Please install Homebrew first.")
            return False
        
        success, output = self.cmd_runner.run(['brew', 'install', 'aquasecurity/trivy/trivy'])
        return success
    
    def _install_trivy_linux(self) -> bool:
        """Install trivy on Linux"""
        # Try different package managers
        if self.check_binary_exists('apt-get'):
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
                success, output = self.cmd_runner.run(cmd)
                if not success:
                    return False
            return True
        
        self.logger.error("Unsupported Linux distribution. Please install Trivy manually.")
        return False
    
    def install_tool(self, tool_name: str) -> bool:
        """Install a specific tool"""
        if tool_name not in self.tools_config:
            self.logger.error(f"Unknown tool: {tool_name}")
            return False
        
        config = self.tools_config[tool_name]
        
        # Check if already installed
        if config['check_method'](tool_name):
            self.logger.info(f"{tool_name} is already installed and up to date")
            return True
        
        # Install the tool
        return config['install_method'](tool_name, config)
    
    def install_tools(self, tools: List[str] = None) -> bool:
        """Install multiple tools"""
        if not tools:
            tools = list(self.tools_config.keys())
        
        # Create a nice panel header
        self.rich.print_panel(
            f"Installing {len(tools)} development tools",
            title="🔧 Tool Installation"
        )
        
        # Create a table of tools to install
        if has_rich():
            console = get_console()
            from rich.table import Table
            from rich.progress import Progress, SpinnerColumn, TextColumn, BarColumn, TaskProgressColumn
            
            table = Table(title="Tools to Install")
            table.add_column("Tool", style="cyan", no_wrap=True)
            table.add_column("Description", style="white")
            table.add_column("Version", style="yellow")
            
            for tool in tools:
                if tool in self.tools_config:
                    config = self.tools_config[tool]
                    version = self.version_helper.get_version(config['version_key'])
                    table.add_row(tool, config['description'], version)
            
            console.print(table)
            console.print()
        else:
            print(f"Installing tools: {', '.join(tools)}")
        
        success_count = 0
        total_count = len(tools)
        
        if has_rich():
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
        if success_count == total_count:
            self.rich.print_panel(
                f"✅ All {success_count} tools installed successfully!",
                title="🎉 Installation Complete",
                style="green"
            )
        else:
            self.rich.print_panel(
                f"⚠️  {success_count}/{total_count} tools installed successfully",
                title="⚠️  Installation Completed with Issues",
                style="yellow"
            )
        
        return success_count == total_count

def main():
    """Main entry point"""
    import sys
    
    # Create installer instance
    installer = ToolInstaller()
    
    # Show header
    installer.rich.print_panel(
        "Development Tool Installer",
        title="Using versions from versions.yml",
        style="bold blue"
    )
    
    # Get tools from command line arguments
    tools = sys.argv[1:] if len(sys.argv) > 1 else None
    
    success = installer.install_tools(tools)
    
    if not success:
        installer.exit_with_error()
    else:
        installer.exit_with_success()

if __name__ == '__main__':
    main() 