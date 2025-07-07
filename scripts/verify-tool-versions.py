#!/usr/bin/env python3

"""
verify-tool-versions.py - Verify actual installed tool versions against versions.yml
Usage: python verify-tool-versions.py [--fix] [--verbose]
"""

import re
import subprocess
import argparse
import json
import shutil
from pathlib import Path
from typing import Dict, Optional, List, Tuple, Set

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
                'go_package': None  # trivy is not a Go tool
            },
            'goose': {
                'cmd': ['goose', '--version'],
                'pattern': r'goose version: v([0-9]+\.[0-9]+\.[0-9]+)',
                'version_key': 'tools.goose',
                'go_package': 'github.com/pressly/goose/v3/cmd/goose'
            },
            'gocov': {
                'cmd': ['go', 'version', '-m'], # Special check
                'pattern': r'\tmod\t.*\t(v[0-9]+\.[0-9]+\.[0-9]+)',
                'version_key': 'tools.gocov',
                'go_package': 'github.com/axw/gocov/gocov',
                'use_go_version_m': True
            },
            'gocov-html': {
                'cmd': ['go', 'version', '-m'], # Special check
                'pattern': r'\tmod\t.*\t(v[0-9]+\.[0-9]+\.[0-9]+)',
                'version_key': 'tools.gocov-html',
                'go_package': 'github.com/matm/gocov-html/cmd/gocov-html',
                'use_go_version_m': True
            },
            'go-cover-treemap': {
                'cmd': ['go', 'version', '-m'], # Special check
                'pattern': r'\tmod\t.*\t(v[0-9]+\.[0-9]+\.[0-9]+)',
                'version_key': 'tools.go-cover-treemap',
                'go_package': 'github.com/nikolaydubina/go-cover-treemap',
                'use_go_version_m': True
            },
            'codeql-cli': {
                'cmd': ['codeql', 'version', '--format=terse'],
                'pattern': r'([0-9]+\.[0-9]+\.[0-9]+)',
                'version_key': 'tools.codeql-cli',
                'go_package': None
            }
        }
    
    def _run_version_command(self, tool_name: str, config: Dict) -> Tuple[Optional[str], Optional[str]]:
        """Run version command for a tool and extract version and any warnings."""
        warning = None
        cmd = config['cmd']

        # For tools without a version flag, check binary build info
        if config.get('use_go_version_m'):
            tool_path = shutil.which(tool_name)
            if not tool_path:
                return None, f"{tool_name} not found in PATH"
            cmd = cmd + [tool_path]

        try:
            result = subprocess.run(
                cmd,
                capture_output=True, 
                text=True, 
                timeout=10
            )
            
            if result.returncode != 0:
                warning = f"Failed to get {tool_name} version: {result.stderr.strip()}"
                return None, warning
            
            output = result.stdout + result.stderr
            match = re.search(config['pattern'], output)
            
            if match:
                version = match.group(1)
                if version == 'dev' or version == 'devel':
                    warning = f"{tool_name} reports 'dev' version - installed from source"
                    return 'dev', warning
                return version.lstrip('v'), None
            else:
                warning = f"Could not parse {tool_name} version from: {output.strip()}"
                return None, warning
                
        except subprocess.TimeoutExpired:
            return None, f"Timeout getting {tool_name} version"
        except FileNotFoundError:
            return None, f"{tool_name} not found in PATH"
        except Exception as e:
            return None, f"Error getting {tool_name} version: {e}"
    
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

        # --- Phase 1: Gather all data silently ---
        results = []
        all_good = True
        
        for tool_name, config in self.version_commands.items():
            installed_version, warning = self._run_version_command(tool_name, config)
            expected_version = self._get_expected_version(config['version_key'])
            
            status = "✅ Match"
            if installed_version is None:
                status = "❌ Missing"
                all_good = False
                self.mismatches.append({
                    'tool': tool_name, 'installed': None, 'expected': expected_version, 'status': 'missing'
                })
            elif expected_version is None:
                status = "❌ Undefined"
                all_good = False
            elif not self._versions_match(installed_version, expected_version):
                status = "❌ Mismatch"
                all_good = False
                self.mismatches.append({
                    'tool': tool_name, 'installed': installed_version, 'expected': expected_version, 'status': 'mismatch'
                })

            results.append({
                "tool": tool_name,
                "installed": installed_version or "N/A",
                "expected": expected_version or "?",
                "status": status,
                "warning": warning
            })

        # --- Phase 2: Print clean report ---
        table = self.rich.create_table(title="🔍 Tool Version Check Results")
        table.add_column("Tool", style="cyan", no_wrap=True)
        table.add_column("Installed", style="yellow")
        table.add_column("Expected", style="blue")
        table.add_column("Status", style="bold")
        
        warnings_to_show = []
        for res in results:
            table.add_row(res['tool'], res['installed'], res['expected'], res['status'])
            if res['warning']:
                warnings_to_show.append(res['warning'])
        
        self.rich.print_table(table)
        print() # Add a newline for spacing

        if warnings_to_show:
            self.rich.print_panel("\n".join(f"⚠️  {w}" for w in warnings_to_show), title="Notes", style="yellow")
            print()

        # --- Phase 3: Show summary and handle fixes ---
        if all_good:
            self.rich.print_panel("✅ All tool versions match versions.yml", title="🎉 All Good!", style="green")
        else:
            self.rich.print_panel(f"❌ Found {len(self.mismatches)} version mismatches", title="⚠️  Mismatches Detected", style="red")
        
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

    def check_outdated_tools(self):
        """Check for outdated Go tools."""
        self.rich.print_panel("Checking for outdated tools...", style="bold blue")
        
        outdated_tools = []
        for tool_name, config in self.version_commands.items():
            if not config.get('go_package'):
                self.logger.verbose(f"Skipping update check for {tool_name} (not a Go tool).")
                continue

            package = config['go_package']
            # Strip subdirectories like /cmd/ for the check
            repo_path = re.sub(r'/(v[0-9]+|cmd)/.*', '', package)

            current_version = self.version_helper.get_version(config['version_key'])
            
            self.logger.verbose(f"Checking {repo_path} for updates...")
            success, result = self.cmd_runner.run(['go', 'list', '-m', '-u', '-json', f'{repo_path}@latest'])
            
            if not success:
                self.logger.error(f"Failed to check for updates for {repo_path}: {result}")
                continue

            try:
                data = json.loads(result)
                if 'Update' in data:
                    latest_version = data['Update']['Version']
                    if latest_version != current_version:
                        outdated_tools.append({
                            "tool": tool_name,
                            "package": repo_path,
                            "current": current_version,
                            "latest": latest_version
                        })
            except json.JSONDecodeError:
                self.logger.error(f"Failed to parse JSON for {repo_path}.")
        
        if not outdated_tools:
            self.logger.success("All Go tools in versions.yml are up to date.")
        else:
            table = self.rich.create_table(title="Outdated Go Tools")
            table.add_column("Tool", style="cyan")
            table.add_column("Go Package", style="white")
            table.add_column("Current Version", style="yellow")
            table.add_column("Latest Version", style="green")
            for tool in outdated_tools:
                table.add_row(tool['tool'], tool['package'], tool['current'], tool['latest'])
            self.rich.print_table(table)

def main():
    """Main entry point"""
    parser = argparse.ArgumentParser(description="Verify tool versions against versions.yml")
    parser.add_argument(
        '--fix',
        action='store_true',
        help='Automatically install the correct version of mismatched tools.'
    )
    parser.add_argument(
        '--check-outdated',
        action='store_true',
        help='Check for outdated versions of Go tools.'
    )
    parser.add_argument(
        '--verbose',
        action='store_true',
        help='Enable verbose output for debugging.'
    )
    args = parser.parse_args()

    verifier = ToolVersionVerifier(fix_versions=args.fix, verbose=args.verbose)
    
    if args.check_outdated:
        verifier.check_outdated_tools()
    else:
        success = verifier.verify_tool_versions()
        if not success:
            verifier.exit_with_error("Version verification failed.")
        else:
            verifier.exit_with_success("All tool versions are correct.")

if __name__ == '__main__':
    main() 