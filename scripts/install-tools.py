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
import tarfile
import subprocess
import re
from pathlib import Path
from typing import List, Dict, Optional, Tuple

# Import common utilities
from common import ScriptBase, has_rich, get_console, has_requests

class ToolInstaller(ScriptBase):
    """Tool installation manager with rich output and proper version checking"""
    
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
                'check_method': self._check_go_tool_version,
                'description': 'Test runner',
                'package': 'gotest.tools/gotestsum'
            },
            'gosec': {
                'version_key': 'tools.gosec',
                'install_method': self._install_go_tool,
                'check_method': self._check_go_tool_version,
                'description': 'Security analyzer',
                'package': 'github.com/securego/gosec/v2/cmd/gosec'
            },
            'govulncheck': {
                'version_key': 'tools.govulncheck',
                'install_method': self._install_go_tool,
                'check_method': self._check_go_tool_version,
                'description': 'Vulnerability checker',
                'package': 'golang.org/x/vuln/cmd/govulncheck'
            },
            'air': {
                'version_key': 'tools.air',
                'install_method': self._install_go_tool,
                'check_method': self._check_go_tool_version,
                'description': 'Live reload',
                'package': 'github.com/air-verse/air'
            },
            'gocov': {
                'version_key': 'tools.gocov',
                'install_method': self._install_go_tool,
                'check_method': self._check_go_tool_version,
                'description': 'Coverage tool',
                'package': 'github.com/axw/gocov/gocov'
            },
            'gocov-html': {
                'version_key': 'tools.gocov-html',
                'install_method': self._install_go_tool,
                'check_method': self._check_go_tool_version,
                'description': 'Coverage HTML generator',
                'package': 'github.com/matm/gocov-html/cmd/gocov-html'
            },
            'go-cover-treemap': {
                'version_key': 'tools.go-cover-treemap',
                'install_method': self._install_go_tool,
                'check_method': self._check_go_tool_version,
                'description': 'Coverage treemap',
                'package': 'github.com/nikolaydubina/go-cover-treemap'
            },
            'trivy': {
                'version_key': 'tools.trivy',
                'install_method': self._install_trivy_versioned,
                'check_method': self._check_trivy_version,
                'description': 'Security scanner'
            },
            'goose': {
                'version_key': 'tools.goose',
                'install_method': self._install_go_tool,
                'check_method': self._check_go_tool_version,
                'description': 'Database migration tool',
                'package': 'github.com/pressly/goose/v3/cmd/goose'
            },
            'codeql-cli': {
                'version_key': 'tools.codeql-cli',
                'install_method': self._install_codeql_cli,
                'check_method': self._check_codeql_version,
                'description': 'Code analysis tool'
            },
            'sqlc': {
                'version_key': 'tools.sqlc',
                'install_method': self._install_sqlc,
                'check_method': self._check_sqlc_version,
                'description': 'SQL code generator'
            }
        }

    def _get_tool_version(self, tool_name: str, cmd: List[str], pattern: str) -> Optional[str]:
        """Get version of an installed tool"""
        try:
            result = subprocess.run(cmd, capture_output=True, text=True, timeout=10)
            if result.returncode == 0:
                output = result.stdout + result.stderr
                match = re.search(pattern, output)
                if match:
                    return match.group(1).lstrip('v')
        except:
            pass
        return None

    def _check_tool_version(self, tool_name: str, version_key: str, version_cmd: List[str], version_pattern: str) -> bool:
        """Check if tool is installed with correct version"""
        if not self.check_binary_exists(tool_name):
            return False
        
        expected_version = self.version_helper.get_version(version_key)
        if not expected_version:
            self.logger.warn(f"No expected version found for {tool_name}")
            return True  # If no version specified, assume current is fine
        
        installed_version = self._get_tool_version(tool_name, version_cmd, version_pattern)
        if not installed_version:
            self.logger.warn(f"Could not determine {tool_name} version")
            return False  # Can't determine version, assume needs reinstall
        
        # Normalize versions for comparison
        expected_clean = expected_version.lstrip('v')
        installed_clean = installed_version.lstrip('v')
        
        is_correct = expected_clean == installed_clean
        if not is_correct:
            self.logger.info(f"{tool_name}: installed={installed_clean}, expected={expected_clean}")
        
        return is_correct

    def _check_go_tool_version(self, tool_name: str) -> bool:
        """Check Go tool version using go version -m"""
        if not self.check_binary_exists(tool_name):
            return False
        
        expected_version = self.version_helper.get_version(self.tools_config[tool_name]['version_key'])
        if not expected_version:
            return True
        
        # Use go version -m to get module info
        tool_path = shutil.which(tool_name)
        if not tool_path:
            return False
        
        try:
            result = subprocess.run(['go', 'version', '-m', tool_path], 
                                    capture_output=True, text=True, timeout=10)
            if result.returncode == 0:
                # Look for version in the build info
                for line in result.stdout.split('\n'):
                    if 'mod' in line and self.tools_config[tool_name]['package'] in line:
                        parts = line.split()
                        if len(parts) >= 3:
                            installed_version = parts[2].lstrip('v')
                            expected_clean = expected_version.lstrip('v')
                            
                            # Handle special cases
                            if expected_version == 'dev' or installed_version == '(devel)':
                                return True
                            
                            if expected_clean == installed_version:
                                return True
                            
                            self.logger.info(f"{tool_name}: installed={installed_version}, expected={expected_clean}")
                            return False
        except:
            pass
        
        return False

    def _check_golangci_lint(self, tool_name: str) -> bool:
        """Check golangci-lint version"""
        return self._check_tool_version(
            tool_name, 'tools.golangci-lint',
            ['golangci-lint', '--version'],
            r'golangci-lint has version ([0-9]+\.[0-9]+\.[0-9]+)'
        )

    def _check_trivy_version(self, tool_name: str) -> bool:
        """Check trivy version"""
        return self._check_tool_version(
            tool_name, 'tools.trivy',
            ['trivy', '--version'],
            r'Version: ([0-9]+\.[0-9]+\.[0-9]+)'
        )

    def _check_codeql_version(self, tool_name: str) -> bool:
        """Check CodeQL version"""
        return self._check_tool_version(
            tool_name, 'tools.codeql-cli',
            ['codeql', 'version'],
            r'CodeQL command-line toolchain release ([0-9]+\.[0-9]+\.[0-9]+)'
        )

    def _check_sqlc_version(self, tool_name: str) -> bool:
        """Check SQLC version"""
        return self._check_tool_version(
            tool_name, 'tools.sqlc',
            ['sqlc', 'version'],
            r'v([0-9]+\.[0-9]+\.[0-9]+)'
        )

    def _install_go_tool(self, tool_name: str, config: Dict) -> bool:
        """Install a Go tool using go install"""
        version = self.version_helper.get_version(config['version_key'])
        package = config['package']
        
        # Ensure version is correctly formatted for go install
        if version:
            if version == 'dev':
                version = 'latest'
            elif version != 'latest' and not version.startswith('v'):
                version = f"v{version}"
            
        package_with_version = f"{package}@{version}" if version else package
        
        success, output = self.cmd_runner.run_with_status(
            ['go', 'install', package_with_version],
            f"Installing {tool_name}@{version}..."
        )

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

    def _install_trivy_versioned(self, tool_name: str, config: Dict) -> bool:
        """Install trivy with specific version from GitHub releases"""
        version = self.version_helper.get_version(config['version_key'])
        if not version:
            self.logger.error(f"No version specified for {tool_name}")
            return False
        
        system = platform.system().lower()
        arch = platform.machine().lower()
        
        # Map platform names
        if system == "darwin":
            platform_str = "macOS"
        elif system == "linux":
            platform_str = "Linux"
        else:
            self.logger.error(f"Unsupported OS for trivy: {system}")
            return False
        
        # Map architecture names
        if arch in ["x86_64", "amd64"]:
            arch_str = "64bit"
        elif arch in ["arm64", "aarch64"]:
            arch_str = "ARM64"
        else:
            self.logger.error(f"Unsupported architecture for trivy: {arch}")
            return False
        
        # Download URL format
        filename = f"trivy_{version}_{platform_str}-{arch_str}.tar.gz"
        download_url = f"https://github.com/aquasecurity/trivy/releases/download/v{version}/{filename}"
        
        # Get GOPATH/bin for install location
        try:
            result = subprocess.run(['go', 'env', 'GOPATH'], capture_output=True, text=True)
            if result.returncode != 0:
                self.logger.error("Failed to get GOPATH")
                return False
            install_dir = Path(result.stdout.strip()) / "bin"
            install_dir.mkdir(parents=True, exist_ok=True)
        except:
            self.logger.error("Failed to determine install directory")
            return False
        
        self.logger.info(f"Installing trivy v{version}...")
        
        try:
            with tempfile.TemporaryDirectory() as tmpdir:
                tar_path = Path(tmpdir) / filename
                
                # Download
                response = requests.get(download_url, stream=True, timeout=120)
                response.raise_for_status()
                with open(tar_path, 'wb') as f:
                    shutil.copyfileobj(response.raw, f)
                
                # Extract
                with tarfile.open(tar_path, 'r:gz') as tar_ref:
                    tar_ref.extractall(tmpdir)
                
                # Move binary to install directory
                binary_path = Path(tmpdir) / "trivy"
                if binary_path.exists():
                    target_path = install_dir / "trivy"
                    shutil.move(str(binary_path), str(target_path))
                    os.chmod(target_path, 0o755)
                    self.logger.success(f"✅ Trivy v{version} installed to {target_path}")
                    return True
                else:
                    self.logger.error("Trivy binary not found in archive")
                    return False
                    
        except Exception as e:
            self.logger.error(f"Failed to install trivy: {e}")
            return False

    def _install_golangci_lint_from_script(self, tool_name: str, config: Dict) -> bool:
        """Install golangci-lint using the official installer script"""
        version = self.version_helper.get_version(config['version_key'])
        if not version:
            self.logger.error(f"No version specified for {tool_name}")
            return False

        success, gopath_output = self.cmd_runner.run(['go', 'env', 'GOPATH'])
        if not success:
            self.logger.error("Failed to get GOPATH for golangci-lint installation.")
            return False

        install_dir = Path(gopath_output.strip().rstrip('/')) / "bin"
        install_dir.mkdir(parents=True, exist_ok=True)

        script_url = "https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh"
        
        try:
            cmd = f"curl -sSfL {script_url} | sh -s -- -b {install_dir} v{version}"
            result = subprocess.run(cmd, shell=True, capture_output=True, text=True, timeout=120)
            
            if result.returncode == 0:
                self.logger.success(f"✅ golangci-lint v{version} installed")
                return True
            else:
                self.logger.error(f"Failed to install golangci-lint: {result.stderr}")
                return False
                
        except Exception as e:
            self.logger.error(f"Error installing golangci-lint: {e}")
            return False

    def _install_codeql_cli(self, tool_name: str, config: Dict) -> bool:
        """Install CodeQL CLI from GitHub releases."""
        version = self.version_helper.get_version(config['version_key'])
        if not version:
            self.logger.error("CodeQL CLI version not found in versions.yml")
            return False

        system = platform.system().lower()
        
        if system == "darwin":
            platform_str = "osx64"
        elif system == "linux":
            platform_str = "linux64"
        else:
            self.logger.error(f"Unsupported OS for CodeQL CLI: {system}")
            return False

        download_url = f"https://github.com/github/codeql-cli-binaries/releases/download/{version}/codeql-{platform_str}.zip"
        
        # Get GOPATH/bin for install location
        try:
            result = subprocess.run(['go', 'env', 'GOPATH'], capture_output=True, text=True)
            if result.returncode != 0:
                self.logger.error("Failed to get GOPATH")
                return False
            install_dir = Path(result.stdout.strip()) / "bin"
            install_dir.mkdir(parents=True, exist_ok=True)
            codeql_install_base_dir = install_dir.parent / "codeql"
            codeql_install_base_dir.mkdir(parents=True, exist_ok=True)
        except:
            self.logger.error("Failed to determine install directory")
            return False

        self.logger.info(f"Installing CodeQL v{version}...")

        try:
            with tempfile.TemporaryDirectory() as tmpdir:
                zip_path = Path(tmpdir) / "codeql.zip"
                
                # Download
                response = requests.get(download_url, stream=True, timeout=120)
                response.raise_for_status()
                with open(zip_path, 'wb') as f:
                    shutil.copyfileobj(response.raw, f)
                
                # Extract
                with zipfile.ZipFile(zip_path, 'r') as zip_ref:
                    zip_ref.extractall(tmpdir)
                
                # Move codeql directory
                codeql_dir = Path(tmpdir) / "codeql"
                if codeql_dir.exists():
                    if codeql_install_base_dir.exists():
                        shutil.rmtree(codeql_install_base_dir)
                    shutil.move(str(codeql_dir), str(codeql_install_base_dir))
                    
                    # Create symlink in bin directory
                    codeql_bin = install_dir / "codeql"
                    if codeql_bin.exists():
                        codeql_bin.unlink()
                    codeql_bin.symlink_to(codeql_install_base_dir / "codeql")
                    
                    self.logger.success(f"✅ CodeQL v{version} installed")
                    return True
                else:
                    self.logger.error("CodeQL directory not found in archive")
                    return False
                    
        except Exception as e:
            self.logger.error(f"Failed to install CodeQL: {e}")
            return False

    def _install_sqlc(self, tool_name: str, config: Dict) -> bool:
        """Install SQLC from GitHub releases."""
        version = self.version_helper.get_version(config['version_key'])
        if not version:
            self.logger.error("SQLC version not found in versions.yml")
            return False

        system = platform.system().lower()
        arch = platform.machine().lower()

        if system == "darwin":
            platform_str = "darwin"
        elif system == "linux":
            platform_str = "linux"
        else:
            self.logger.error(f"Unsupported OS for SQLC: {system}")
            return False

        if arch == "arm64":
            arch_str = "arm64"
        elif arch == "x86_64":
            arch_str = "amd64"
        else:
            self.logger.error(f"Unsupported ARCH for SQLC: {arch}")
            return False

        version_clean = version.lstrip('v')
        download_url = f"https://github.com/sqlc-dev/sqlc/releases/download/{version}/sqlc_{version_clean}_{platform_str}_{arch_str}.tar.gz"
        
        # Get GOPATH/bin for install location
        try:
            result = subprocess.run(['go', 'env', 'GOPATH'], capture_output=True, text=True)
            if result.returncode != 0:
                self.logger.error("Failed to get GOPATH")
                return False
            install_dir = Path(result.stdout.strip()) / "bin"
            install_dir.mkdir(parents=True, exist_ok=True)
        except:
            self.logger.error("Failed to determine install directory")
            return False
        
        self.logger.info(f"Installing SQLC {version}...")
        
        try:
            with tempfile.TemporaryDirectory() as tmpdir:
                tar_path = Path(tmpdir) / "sqlc.tar.gz"
                
                # Download
                response = requests.get(download_url, stream=True, timeout=120)
                response.raise_for_status()
                with open(tar_path, 'wb') as f:
                    shutil.copyfileobj(response.raw, f)
                
                # Extract
                with tarfile.open(tar_path, 'r:gz') as tar_ref:
                    tar_ref.extractall(tmpdir)

                # Move sqlc binary to install_dir
                binary_path = Path(tmpdir) / "sqlc"
                if binary_path.exists():
                    target_path = install_dir / "sqlc"
                    shutil.move(str(binary_path), str(target_path))
                    os.chmod(target_path, 0o755)
                    self.logger.success(f"✅ SQLC {version} installed")
                    return True
                else:
                    self.logger.error("SQLC binary not found in archive")
                    return False
                    
        except Exception as e:
            self.logger.error(f"Failed to install SQLC: {e}")
            return False

    def install_tool(self, tool_name: str) -> bool:
        """Install a specific tool"""
        if tool_name not in self.tools_config:
            self.logger.error(f"Unknown tool: {tool_name}")
            return False
        
        config = self.tools_config[tool_name]
        return config['install_method'](tool_name, config)
    
    def install_tools(self, tools: List[str] = None) -> bool:
        """Install multiple tools with version awareness"""
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
            table.add_column("Expected Version", style="yellow")
            table.add_column("Status", style="bold")
            
            for tool in tools:
                if tool in self.tools_config:
                    config = self.tools_config[tool]
                    version = self.version_helper.get_version(config['version_key'])
                    
                    # Check current status
                    if config['check_method'](tool):
                        status = "✅ Current"
                    else:
                        status = "🔄 Needs Install/Update"
                    
                    table.add_row(tool, config['description'], version, status)
            
            self.rich.print_table(table)
            console.print()

        # --- Phase 1: Check all tools ---
        self.rich.print_panel("1. Checking Tool Versions", style="bold blue")
        tools_to_install = []
        for tool_name in tools:
            if tool_name not in self.tools_config:
                continue
                
            config = self.tools_config[tool_name]
            if config['check_method'](tool_name):
                self.logger.info(f"{tool_name} is already installed with correct version")
            else:
                self.logger.warn(f"{tool_name} needs installation or update")
                tools_to_install.append(tool_name)

        # --- Phase 2: Install missing/outdated tools ---
        if not tools_to_install:
            self.logger.success("✅ All tools are already installed with correct versions")
            return True
        
        self.rich.print_panel(f"2. Installing {len(tools_to_install)} Tool(s)", style="bold blue")
        
        success_count = 0
        for tool_name in tools_to_install:
            if self.install_tool(tool_name):
                success_count += 1
            else:
                self.logger.error(f"Failed to install {tool_name}")

        # --- Phase 3: Final Report ---
        total_success = len(tools) - len(tools_to_install) + success_count
        
        if total_success == len(tools):
            self.rich.print_panel(
                f"✅ All {len(tools)} tools installed successfully",
                title="🎉 Installation Complete!",
                style="green"
            )
            return True
        else:
            self.rich.print_panel(
                f"⚠️  {total_success}/{len(tools)} tools installed successfully",
                title="Partial Success",
                style="yellow"
            )
            return False

def main():
    """Main entry point"""
    import sys
    
    # Create installer instance
    installer = ToolInstaller()
    
    # Show header
    installer.rich.print_panel(
        "Development Tool Installer",
        title="Version-aware tool installation using versions.yml",
        style="bold blue"
    )
    
    # Get tools from command line arguments
    tools = sys.argv[1:] if len(sys.argv) > 1 else None
    
    success = installer.install_tools(tools)
    
    if not success:
        installer.exit_with_error()
    else:
        installer.exit_with_success()

if __name__ == "__main__":
    main() 