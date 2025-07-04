#!/usr/bin/env python3

"""
verify-tool-versions.py - Verify actual installed tool versions against versions.yml
Usage: python verify-tool-versions.py [--fix] [--verbose]
"""

import re
import subprocess
import argparse
from pathlib import Path
from typing import Dict, Optional, List, Tuple

# Import common utilities
from common import ScriptBase, has_rich, get_console

class ToolVersionVerifier(ScriptBase):
    """Tool version verification with actual installed version checking"""
    
    def __init__(self, fix_versions: bool = False, verbose: bool = False):
        super().__init__("ToolVersionVerifier")
        self.fix_versions = fix_versions
        self.verbose = verbose
        self.versions = self.version_helper.load_versions_file()
        self.mismatches = []
        
        # Tool version checking commands
        self.version_commands = {
            'golangci-lint': {
                'cmd': ['golangci-lint', '--version'],
                'pattern': r'golangci-lint has version ([0-9]+\.[0-9]+\.[0-9]+)',
                'version_key': 'tools.golangci-lint',
                'go_package': 'github.com/golangci/golangci-lint/cmd/golangci-lint'
            },
            'gotestsum': {
                'cmd': ['gotestsum', '--version'],
                'pattern': r'gotestsum version (.+)',
                'version_key': 'tools.gotestsum',
                'go_package': 'gotest.tools/gotestsum'
            },
            'gosec': {
                'cmd': ['gosec', '--version'],
                'pattern': r'Version: (.+)',
                'version_key': 'tools.gosec',
                'go_package': 'github.com/securego/gosec/v2/cmd/gosec'
            },
            'govulncheck': {
                'cmd': ['govulncheck', '--version'],
                'pattern': r'govulncheck@v([0-9]+\.[0-9]+\.[0-9]+)',
                'version_key': 'tools.govulncheck',
                'go_package': 'golang.org/x/vuln/cmd/govulncheck'
            },
            'air': {
                'cmd': ['air', '-v'],
                'pattern': r'v([0-9]+\.[0-9]+\.[0-9]+)',
                'version_key': 'tools.air',
                'go_package': 'github.com/air-verse/air'
            },
            'trivy': {
                'cmd': ['trivy', '--version'],
                'pattern': r'Version: ([0-9]+\.[0-9]+\.[0-9]+)',
                'version_key': 'tools.trivy',
                'go_package': None  # trivy is not a Go tool, installed differently
            }
        }
    
    def _run_version_command(self, tool_name: str, config: Dict) -> Optional[str]:
        """Run version command for a tool and extract version"""
        try:
            result = subprocess.run(
                config['cmd'], 
                capture_output=True, 
                text=True, 
                timeout=10
            )
            
            if result.returncode != 0:
                self.logger.warn(f"Failed to get {tool_name} version: {result.stderr.strip()}")
                return None
            
            output = result.stdout + result.stderr
            match = re.search(config['pattern'], output)
            
            if match:
                version = match.group(1)
                # Clean up version (remove 'dev', handle special cases)
                if version == 'dev' or version == 'devel':
                    self.logger.warn(f"{tool_name} reports 'dev' version - installed from source")
                    return 'dev'
                return version
            else:
                self.logger.warn(f"Could not parse {tool_name} version from: {output.strip()}")
                return None
                
        except subprocess.TimeoutExpired:
            self.logger.error(f"Timeout getting {tool_name} version")
            return None
        except FileNotFoundError:
            self.logger.warn(f"{tool_name} not found in PATH")
            return None
        except Exception as e:
            self.logger.error(f"Error getting {tool_name} version: {e}")
            return None
    
    def _get_expected_version(self, version_key: str) -> Optional[str]:
        """Get expected version from versions.yml"""
        try:
            return self.version_helper.get_version(version_key)
        except Exception as e:
            self.logger.error(f"Error getting expected version for {version_key}: {e}")
            return None
    
    def _versions_match(self, installed: str, expected: str) -> bool:
        """Check if installed version matches expected version"""
        if installed == expected:
            return True
        
        # Handle version prefix variations (v1.2.3 vs 1.2.3)
        if installed.lstrip('v') == expected.lstrip('v'):
            return True
        
        # Handle dev versions
        if installed == 'dev':
            return False
        
        return False
    
    def verify_tool_versions(self) -> bool:
        """Verify all tool versions"""
        self.rich.print_panel(
            "Verifying actual installed tool versions",
            title="🔧 Tool Version Verification"
        )
        
        all_good = True
        
        if has_rich():
            console = get_console()
            from rich.table import Table
            from rich.progress import Progress, SpinnerColumn, TextColumn
            
            # Create results table
            table = Table(title="🔍 Tool Version Check Results")
            table.add_column("Tool", style="cyan", no_wrap=True)
            table.add_column("Installed", style="yellow")
            table.add_column("Expected", style="blue")
            table.add_column("Status", style="bold")
            
            with Progress(
                SpinnerColumn(),
                TextColumn("[bold blue]{task.description}"),
                console=console
            ) as progress:
                task = progress.add_task("Checking tool versions...", total=len(self.version_commands))
                
                for tool_name, config in self.version_commands.items():
                    progress.update(task, description=f"Checking {tool_name}...")
                    
                    installed_version = self._run_version_command(tool_name, config)
                    expected_version = self._get_expected_version(config['version_key'])
                    
                    if installed_version is None:
                        table.add_row(tool_name, "❌ Not found", expected_version or "?", "❌ Missing")
                        all_good = False
                        self.mismatches.append({
                            'tool': tool_name,
                            'installed': None,
                            'expected': expected_version,
                            'status': 'missing'
                        })
                    elif expected_version is None:
                        table.add_row(tool_name, installed_version, "❌ Not in versions.yml", "❌ Undefined")
                        all_good = False
                    elif self._versions_match(installed_version, expected_version):
                        table.add_row(tool_name, installed_version, expected_version, "✅ Match")
                    else:
                        table.add_row(tool_name, installed_version, expected_version, "❌ Mismatch")
                        all_good = False
                        self.mismatches.append({
                            'tool': tool_name,
                            'installed': installed_version,
                            'expected': expected_version,
                            'status': 'mismatch'
                        })
                    
                    progress.advance(task)
            
            console.print(table)
            console.print()
            
            # Show summary
            if all_good:
                from rich.panel import Panel
                console.print(Panel(
                    "✅ All tool versions match versions.yml",
                    title="🎉 All Good!",
                    style="green"
                ))
            else:
                from rich.panel import Panel
                console.print(Panel(
                    f"❌ Found {len(self.mismatches)} version mismatches",
                    title="⚠️  Mismatches Detected",
                    style="red"
                ))
        else:
            # Fallback for environments without rich
            for tool_name, config in self.version_commands.items():
                print(f"Checking {tool_name}...")
                installed_version = self._run_version_command(tool_name, config)
                expected_version = self._get_expected_version(config['version_key'])
                
                if installed_version and expected_version:
                    if self._versions_match(installed_version, expected_version):
                        print(f"  ✅ {tool_name}: {installed_version} (matches)")
                    else:
                        print(f"  ❌ {tool_name}: {installed_version} (expected {expected_version})")
                        all_good = False
                        self.mismatches.append({
                            'tool': tool_name,
                            'installed': installed_version,
                            'expected': expected_version,
                            'status': 'mismatch'
                        })
                elif not installed_version:
                    print(f"  ❌ {tool_name}: Not found")
                    all_good = False
        
        if self.fix_versions and self.mismatches:
            self._fix_tool_versions()
        
        return all_good
    
    def _install_go_tool(self, tool_name: str, go_package: str, version: str) -> bool:
        """Install a Go tool at the specified version"""
        try:
            # Special handling for gosec using official installer
            if tool_name == 'gosec':
                return self._install_gosec_with_script(version)
            
            # Standard go install for other tools
            install_target = f"{go_package}@v{version}"
            cmd = ['go', 'install', install_target]
            
            self.logger.info(f"Installing {tool_name} version {version}...")
            self.logger.verbose(f"Running: {' '.join(cmd)}", self.verbose)
            
            result = subprocess.run(
                cmd,
                capture_output=True,
                text=True,
                timeout=60  # Give it more time for installation
            )
            
            if result.returncode == 0:
                self.logger.success(f"✅ Successfully installed {tool_name} v{version}")
                return True
            else:
                self.logger.error(f"❌ Failed to install {tool_name}: {result.stderr.strip()}")
                return False
                
        except subprocess.TimeoutExpired:
            self.logger.error(f"❌ Timeout installing {tool_name}")
            return False
        except Exception as e:
            self.logger.error(f"❌ Error installing {tool_name}: {e}")
            return False

    def _install_gosec_with_script(self, version: str) -> bool:
        """Install gosec using the official installation script"""
        try:
            # Get GOPATH
            gopath_result = subprocess.run(
                ['go', 'env', 'GOPATH'],
                capture_output=True,
                text=True,
                timeout=10
            )
            
            if gopath_result.returncode != 0:
                self.logger.error("❌ Failed to get GOPATH")
                return False
            
            gopath = gopath_result.stdout.strip()
            bin_dir = f"{gopath}/bin"
            
            # Download and run gosec installer
            install_url = "https://raw.githubusercontent.com/securego/gosec/master/install.sh"
            version_with_prefix = f"v{version}"
            
            self.logger.info(f"Installing gosec version {version} using official installer...")
            self.logger.verbose(f"Running: curl -sfL {install_url} | sh -s -- -b {bin_dir} {version_with_prefix}", self.verbose)
            
            # Use shell=True to properly handle pipes
            cmd = f"curl -sfL {install_url} | sh -s -- -b {bin_dir} {version_with_prefix}"
            result = subprocess.run(
                cmd,
                shell=True,
                capture_output=True,
                text=True,
                timeout=120  # Give more time for download
            )
            
            if result.returncode == 0:
                self.logger.success(f"✅ Successfully installed gosec v{version}")
                return True
            else:
                self.logger.error(f"❌ Failed to install gosec: {result.stderr.strip()}")
                return False
                
        except subprocess.TimeoutExpired:
            self.logger.error(f"❌ Timeout installing gosec")
            return False
        except Exception as e:
            self.logger.error(f"❌ Error installing gosec: {e}")
            return False

    def _fix_tool_versions(self):
        """Install correct tool versions to match versions.yml"""
        if not self.mismatches:
            return
        
        self.rich.print_panel(
            "Installing correct tool versions to match versions.yml",
            title="🔧 Auto-fix Tool Versions"
        )
        
        if has_rich():
            console = get_console()
            from rich.progress import Progress, SpinnerColumn, TextColumn, BarColumn, TimeRemainingColumn
            
            # Filter mismatches to only Go tools that can be installed
            fixable_mismatches = []
            for mismatch in self.mismatches:
                if mismatch['status'] in ['mismatch', 'missing'] and mismatch['expected']:
                    tool_name = mismatch['tool']
                    if tool_name in self.version_commands:
                        go_package = self.version_commands[tool_name].get('go_package')
                        if go_package:
                            fixable_mismatches.append(mismatch)
                        else:
                            self.logger.warn(f"⚠️  Cannot auto-install {tool_name} (not a Go tool)")
                            
            if not fixable_mismatches:
                self.logger.warn("No fixable mismatches found (all tools are either correct or not Go tools)")
                return
            
            with Progress(
                SpinnerColumn(),
                TextColumn("[bold blue]{task.description}"),
                BarColumn(),
                TimeRemainingColumn(),
                console=console
            ) as progress:
                task = progress.add_task("Installing tools...", total=len(fixable_mismatches))
                
                success_count = 0
                for mismatch in fixable_mismatches:
                    tool_name = mismatch['tool']
                    expected_version = mismatch['expected']
                    go_package = self.version_commands[tool_name]['go_package']
                    
                    progress.update(task, description=f"Installing {tool_name} v{expected_version}...")
                    
                    if self._install_go_tool(tool_name, go_package, expected_version):
                        success_count += 1
                    
                    progress.advance(task)
                
                # Show summary
                if success_count == len(fixable_mismatches):
                    from rich.panel import Panel
                    console.print(Panel(
                        f"✅ Successfully installed {success_count} tools",
                        title="🎉 All Tools Fixed!",
                        style="green"
                    ))
                else:
                    from rich.panel import Panel
                    console.print(Panel(
                        f"⚠️  Installed {success_count}/{len(fixable_mismatches)} tools",
                        title="Partial Success",
                        style="yellow"
                    ))
                    
        else:
            # Fallback for environments without rich
            for mismatch in self.mismatches:
                if mismatch['status'] in ['mismatch', 'missing'] and mismatch['expected']:
                    tool_name = mismatch['tool']
                    if tool_name in self.version_commands:
                        go_package = self.version_commands[tool_name].get('go_package')
                        if go_package:
                            self._install_go_tool(tool_name, go_package, mismatch['expected'])

def main():
    """Main entry point"""
    parser = argparse.ArgumentParser(description='Verify tool versions against versions.yml')
    parser.add_argument('--fix', action='store_true', help='Auto-install correct tool versions to match versions.yml')
    parser.add_argument('--verbose', '-v', action='store_true', help='Enable verbose output')
    
    args = parser.parse_args()
    
    # Create verifier instance
    verifier = ToolVersionVerifier(fix_versions=args.fix, verbose=args.verbose)
    
    # Show header
    verifier.rich.print_panel(
        "Tool Version Verification System",
        title="Ensuring tool versions match versions.yml",
        style="bold blue"
    )
    
    success = verifier.verify_tool_versions()
    
    if not success:
        verifier.exit_with_error()
    else:
        verifier.exit_with_success()

if __name__ == '__main__':
    main() 