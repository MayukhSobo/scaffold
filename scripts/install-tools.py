#!/usr/bin/env python3

"""
install-tools.py - Centralized tool installation script with rich formatting
Usage: python install-tools.py [tool1] [tool2] ...
If no tools specified, installs all tools
"""

import platform
import shutil
import os
import tempfile
import requests
import zipfile
from pathlib import Path
from typing import List, Dict, Optional, Tuple

# Import common utilities
from common import ScriptBase, has_rich, get_console, has_requests

class ToolInstaller(ScriptBase):
    """Tool installation manager with rich output"""
    
    def __init__(self):
        super().__init__("ToolInstaller")
        self.versions = self.version_helper.load_versions_file()
        self.tools_config = {
            'golangci-lint': {
                'version_key': 'tools.golangci-lint',
                'install_method': self._install_golangci_lint_from_script,
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
            },
            'goose': {
                'version_key': 'tools.goose',
                'install_method': self._install_go_tool,
                'check_method': self._check_go_tool,
                'description': 'Database migration tool',
                'package': 'github.com/pressly/goose/v3/cmd/goose'
            },
            'codeql-cli': {
                'version_key': 'tools.codeql-cli',
                'install_method': self._install_codeql_cli,
                'check_method': self._check_binary_tool,
                'description': 'CodeQL CLI'
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
        
        # Ensure version is correctly formatted for go install
        if version:
            if version == 'dev':
                # 'dev' often maps to 'latest' for go install
                version = 'latest'
            elif version != 'latest' and not version.startswith('v'):
                version = f"v{version}"
            
        package_with_version = f"{package}@{version}" if version else package
        
        success, output = self.cmd_runner.run_with_status(
            ['go', 'install', package_with_version],
            f"Installing {tool_name}@{version}..."
        )

        # Fallback for specific tools that might have versioning quirks
        if not success and tool_name == 'go-cover-treemap':
            self.logger.warn(f"Version '{version}' for {tool_name} failed. Trying @latest...")
            package_with_version = f"{package}@latest"
            success, output = self.cmd_runner.run_with_status(
                ['go', 'install', package_with_version],
                f"Installing {tool_name}@latest..."
            )
        
        if success:
            self.logger.info(f"{tool_name} installed successfully")
            return True
        else:
            self.logger.error(f"Failed to install {tool_name}: {output}")
            return False
    
    def _install_golangci_lint_from_script(self, tool_name: str, config: Dict) -> bool:
        """Install golangci-lint using the official install script."""
        version = self.version_helper.get_version(config['version_key'])
        if version and not version.startswith('v'):
            version = f"v{version}"

        if not has_requests():
            self.logger.error("The 'requests' library is required to download the installer. Please run 'pip install requests'.")
            return False

        # Get GOPATH/bin
        success, gopath_output = self.cmd_runner.run(['go', 'env', 'GOPATH'])
        if not success:
            self.logger.error("Failed to get GOPATH for golangci-lint installation.")
            return False
        install_dir = str(Path(gopath_output.strip()) / "bin")

        installer_path = None
        try:
            # 1. Download the script content
            url = "https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh"
            self.logger.verbose(f"Downloading installer from {url}")
            response = requests.get(url, timeout=30)
            response.raise_for_status()
            
            # 2. Write to a temporary file
            with tempfile.NamedTemporaryFile(mode='w', delete=False, suffix='.sh', prefix='golangci-lint-installer-') as f:
                f.write(response.text)
                installer_path = f.name
            
            # 3. Make executable
            Path(installer_path).chmod(0o755)
            self.logger.verbose(f"Installer script saved to {installer_path}")

            # 4. Run the script
            install_cmd = [installer_path, '-b', install_dir, version]
            success, output = self.cmd_runner.run_with_status(
                install_cmd,
                f"Installing {tool_name} {version}"
            )
            
            if not success:
                self.logger.error(f"Failed to install {tool_name}: {output}")
                return False

        except Exception as e:
            self.logger.error(f"An error occurred during {tool_name} installation: {e}")
            return False
        finally:
            # 5. Clean up installer
            if installer_path and Path(installer_path).exists():
                os.unlink(installer_path)
                self.logger.verbose(f"Removed temporary installer {installer_path}")
        
        return True
    
    def _install_trivy(self, tool_name: str, config: Dict) -> bool:
        """Install trivy using platform-specific method"""
        system = platform.system().lower()
        
        self.logger.info(f"Installing {tool_name}...")
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
    
    def _install_codeql_cli(self, tool_name: str, config: Dict) -> bool:
        """Install CodeQL CLI from GitHub releases."""
        version = self.version_helper.get_version(config['version_key'])
        if not version:
            self.logger.error("CodeQL CLI version not found in versions.yml")
            return False

        system = platform.system().lower()
        arch = platform.machine().lower()

        if system == "darwin":
            platform_str = "osx64"
        elif system == "linux":
            platform_str = "linux64"
        else:
            self.logger.error(f"Unsupported OS for CodeQL CLI: {system}")
            return False

        download_url = f"https://github.com/github/codeql-cli-binaries/releases/download/v{version}/codeql-{platform_str}.zip"
        
        # Get GOPATH/bin to determine install location
        success, gopath_output = self.cmd_runner.run(['go', 'env', 'GOPATH'])
        if not success:
            self.logger.error("Failed to get GOPATH for CodeQL CLI installation.")
            return False
        install_dir = Path(gopath_output.strip().rstrip('/')) / "bin"
        install_dir.mkdir(parents=True, exist_ok=True)
        codeql_install_base_dir = install_dir.parent / "codeql-cli"
        codeql_install_base_dir.mkdir(parents=True, exist_ok=True)

        self.logger.info(f"Downloading CodeQL CLI v{version} to {codeql_install_base_dir}...")
        
        try:
            # Check if already installed
            if (codeql_install_base_dir / "codeql" / "codeql").exists():
                self.logger.info("CodeQL CLI already seems to be installed.")
            else:
                with tempfile.TemporaryDirectory() as tmpdir:
                    zip_path = Path(tmpdir) / "codeql.zip"
                    
                    # Download
                    response = requests.get(download_url, stream=True, timeout=120)
                    response.raise_for_status()
                    with open(zip_path, 'wb') as f:
                        shutil.copyfileobj(response.raw, f)
                    
                    # Unzip
                    with zipfile.ZipFile(zip_path, 'r') as zip_ref:
                        zip_ref.extractall(codeql_install_base_dir)
            
            # Create symlink in gopath/bin
            binary_path = codeql_install_base_dir / "codeql" / "codeql"
            symlink_path = install_dir / "codeql"
            
            if binary_path.exists():
                if symlink_path.exists() and symlink_path.is_symlink():
                    self.logger.info("Symlink already exists.")
                else:
                    os.symlink(binary_path, symlink_path)
                
                # Ensure the binary and all its tools are executable
                binary_path.chmod(0o755)
                
                # Recursively set execute permissions on all files in the tools dir
                tools_dir = codeql_install_base_dir / "codeql" / "tools"
                if tools_dir.is_dir():
                    for f in tools_dir.rglob('*'):
                        if f.is_file():
                            f.chmod(0o755)

                self.logger.success(f"✅ CodeQL CLI symlinked to {symlink_path}")
                return True
            else:
                self.logger.error(f"Could not find 'codeql' binary in the extracted archive.")
                return False
        except Exception as e:
            self.logger.error(f"An error occurred during CodeQL CLI installation: {e}")
            return False

    def install_tool(self, tool_name: str) -> bool:
        """Install a specific tool"""
        if tool_name not in self.tools_config:
            self.logger.error(f"Unknown tool: {tool_name}")
            return False
        
        config = self.tools_config[tool_name]
        
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
            table = self.rich.create_table(title="Tools to Install")
            table.add_column("Tool", style="cyan", no_wrap=True)
            table.add_column("Description", style="white")
            table.add_column("Version", style="yellow")
            
            for tool in tools:
                if tool in self.tools_config:
                    config = self.tools_config[tool]
                    version = self.version_helper.get_version(config['version_key'])
                    table.add_row(tool, config['description'], version)
            
            self.rich.print_table(table)
            console.print()
        else:
            print(f"Installing tools: {', '.join(tools)}")

        # --- Phase 1: Check all tools ---
        self.rich.print_panel("1. Checking Tool Status", style="bold blue")
        tools_to_install = []
        for tool_name in tools:
            config = self.tools_config[tool_name]
            if config['check_method'](tool_name):
                self.logger.info(f"{tool_name} is already installed and up to date.")
            else:
                self.logger.warn(f"{tool_name} is not installed or out of date.")
                tools_to_install.append(tool_name)

        # --- Phase 2: Install missing tools ---
        if not tools_to_install:
            success_count = len(tools)
        else:
            self.rich.print_panel(f"2. Installing {len(tools_to_install)} Missing/Outdated Tool(s)", style="bold blue")
            
            initial_success_count = len(tools) - len(tools_to_install)
            installed_count = 0

            if has_rich():
                from rich.progress import Progress, SpinnerColumn, TextColumn, BarColumn
                from rich_custom_columns import TaskProgressColumn
                with Progress(
                    SpinnerColumn(),
                    TextColumn("[bold blue]{task.description}"),
                    BarColumn(),
                    TaskProgressColumn(),
                    console=get_console()
                ) as progress:
                    task = progress.add_task("Installing...", total=len(tools_to_install))
                    
                    for tool_name in tools_to_install:
                        progress.update(task, description=f"Installing {tool_name}...")
                        if self.install_tool(tool_name):
                            installed_count += 1
                        progress.advance(task)
            else:
                for tool_name in tools_to_install:
                    if self.install_tool(tool_name):
                        installed_count += 1
            
            success_count = initial_success_count + installed_count
        
        total_count = len(tools)
        
        # Summary
        if success_count == total_count:
            self.rich.print_panel(
                f"✅ All {success_count} tools are installed and up to date!",
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