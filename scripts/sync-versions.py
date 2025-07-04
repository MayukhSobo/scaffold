#!/usr/bin/env python3

"""
sync-versions.py - Synchronize versions from versions.yml across all project files
Usage: python sync-versions.py [--dry-run] [--verbose]
"""

import os
import sys
import re
import argparse
from pathlib import Path
from typing import List, Dict, Optional, Tuple, Set
import subprocess

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
    from version_helper import get_version, load_versions_file, list_versions
except ImportError:
    # If running through run.py, try different import
    import importlib.util
    spec = importlib.util.spec_from_file_location("version_helper", scripts_dir / "version-helper.py")
    version_helper = importlib.util.module_from_spec(spec)
    spec.loader.exec_module(version_helper)
    get_version = version_helper.get_version
    load_versions_file = version_helper.load_versions_file
    list_versions = version_helper.list_versions

# Initialize console
console = Console() if HAS_RICH else None

class VersionSyncer:
    """Version synchronization manager with rich output"""
    
    def __init__(self, dry_run: bool = False, verbose: bool = False):
        self.dry_run = dry_run
        self.verbose = verbose
        self.project_root = Path(__file__).parent.parent
        self.versions = load_versions_file()
        self.changes_made = []
        self.errors = []
        
        # Define file patterns and their sync rules
        self.sync_rules = {
            '.github/workflows/ci.yml': [
                {
                    'pattern': r'go-version:\s*[\'"]?([^\'"]+)[\'"]?',
                    'version_key': 'languages.go',
                    'replacement': lambda v: f'go-version: \'{v}\''
                },
                {
                    'pattern': r'golangci-lint-version:\s*[\'"]?([^\'"]+)[\'"]?',
                    'version_key': 'tools.golangci-lint',
                    'replacement': lambda v: f'golangci-lint-version: \'{v}\''
                }
            ],
            '.github/workflows/docker.yml': [
                {
                    'pattern': r'GO_VERSION:\s*[\'"]?([^\'"]+)[\'"]?',
                    'version_key': 'languages.go',
                    'replacement': lambda v: f'GO_VERSION: \'{v}\''
                }
            ],
            '.github/workflows/codeql.yml': [
                {
                    'pattern': r'go-version:\s*[\'"]?([^\'"]+)[\'"]?',
                    'version_key': 'languages.go',
                    'replacement': lambda v: f'go-version: \'{v}\''
                }
            ],
            '.github/workflows/dependencies.yml': [
                {
                    'pattern': r'go-version:\s*[\'"]?([^\'"]+)[\'"]?',
                    'version_key': 'languages.go',
                    'replacement': lambda v: f'go-version: \'{v}\''
                }
            ],
            '.github/workflows/release.yml': [
                {
                    'pattern': r'go-version:\s*[\'"]?([^\'"]+)[\'"]?',
                    'version_key': 'languages.go',
                    'replacement': lambda v: f'go-version: \'{v}\''
                }
            ],
            'Dockerfile': [
                {
                    'pattern': r'FROM\s+golang:([^\s]+)',
                    'version_key': 'languages.go',
                    'replacement': lambda v: f'FROM golang:{v}'
                }
            ],
            'go.mod': [
                {
                    'pattern': r'go\s+([0-9]+\.[0-9]+)',
                    'version_key': 'languages.go',
                    'replacement': lambda v: f'go {v}'
                }
            ],
            'scripts/install-tools.sh': [
                {
                    'pattern': r'GOLANGCI_LINT_VERSION="([^"]+)"',
                    'version_key': 'tools.golangci-lint',
                    'replacement': lambda v: f'GOLANGCI_LINT_VERSION="{v}"'
                }
            ]
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
        self.errors.append(message)
    
    def _log_verbose(self, message: str):
        """Log verbose message"""
        if self.verbose:
            if HAS_RICH:
                console.print(f"ℹ️  {message}", style="blue")
            else:
                print(f"[VERBOSE] {message}")
    
    def _read_file(self, file_path: Path) -> Optional[str]:
        """Read file content"""
        try:
            return file_path.read_text(encoding='utf-8')
        except Exception as e:
            self._log_error(f"Failed to read {file_path}: {e}")
            return None
    
    def _write_file(self, file_path: Path, content: str):
        """Write file content"""
        if self.dry_run:
            self._log_verbose(f"Would write {file_path}")
            return
        
        try:
            file_path.write_text(content, encoding='utf-8')
            self._log_verbose(f"Updated {file_path}")
        except Exception as e:
            self._log_error(f"Failed to write {file_path}: {e}")
    
    def _sync_file(self, file_path: Path, rules: List[Dict]) -> bool:
        """Sync versions in a single file"""
        relative_path = file_path.relative_to(self.project_root)
        
        if not file_path.exists():
            self._log_warn(f"File not found: {relative_path}")
            return False
        
        content = self._read_file(file_path)
        if content is None:
            return False
        
        original_content = content
        changes_in_file = []
        
        for rule in rules:
            pattern = rule['pattern']
            version_key = rule['version_key']
            replacement_func = rule['replacement']
            
            try:
                expected_version = get_version(version_key)
                
                def replace_match(match):
                    old_version = match.group(1)
                    if old_version != expected_version:
                        changes_in_file.append({
                            'version_key': version_key,
                            'old_version': old_version,
                            'new_version': expected_version
                        })
                        return replacement_func(expected_version)
                    return match.group(0)
                
                content = re.sub(pattern, replace_match, content)
                
            except Exception as e:
                self._log_error(f"Error processing {version_key} in {relative_path}: {e}")
                continue
        
        if changes_in_file:
            self._write_file(file_path, content)
            self.changes_made.extend([{
                'file': str(relative_path),
                'changes': changes_in_file
            }])
            return True
        
        return False
    
    def sync_all_files(self) -> bool:
        """Sync versions across all configured files"""
        if HAS_RICH:
            panel = Panel(
                "Synchronizing versions from versions.yml",
                title="🔄 Version Synchronization",
                title_align="left"
            )
            console.print(panel)
            
            if self.dry_run:
                console.print("🔍 [bold yellow]DRY RUN MODE[/bold yellow] - No files will be modified")
                console.print()
        else:
            print("Synchronizing versions from versions.yml")
            if self.dry_run:
                print("DRY RUN MODE - No files will be modified")
        
        total_files = len(self.sync_rules)
        files_changed = 0
        
        if HAS_RICH:
            with Progress(
                SpinnerColumn(),
                TextColumn("[bold blue]{task.description}"),
                BarColumn(),
                TaskProgressColumn(),
                console=console
            ) as progress:
                task = progress.add_task("Syncing files...", total=total_files)
                
                for file_path, rules in self.sync_rules.items():
                    full_path = self.project_root / file_path
                    progress.update(task, description=f"Syncing {file_path}...")
                    
                    if self._sync_file(full_path, rules):
                        files_changed += 1
                    
                    progress.advance(task)
        else:
            for file_path, rules in self.sync_rules.items():
                full_path = self.project_root / file_path
                print(f"Syncing {file_path}...")
                
                if self._sync_file(full_path, rules):
                    files_changed += 1
        
        self._show_summary(files_changed, total_files)
        return len(self.errors) == 0
    
    def _show_summary(self, files_changed: int, total_files: int):
        """Show synchronization summary"""
        if HAS_RICH:
            if self.changes_made:
                # Create a table of changes
                table = Table(title="📝 Changes Made")
                table.add_column("File", style="cyan", no_wrap=True)
                table.add_column("Version Key", style="yellow")
                table.add_column("Old Version", style="red")
                table.add_column("New Version", style="green")
                
                for file_changes in self.changes_made:
                    file_name = file_changes['file']
                    for change in file_changes['changes']:
                        table.add_row(
                            file_name,
                            change['version_key'],
                            change['old_version'],
                            change['new_version']
                        )
                
                console.print(table)
                console.print()
            
            if files_changed > 0:
                status = "green" if not self.dry_run else "yellow"
                action = "Updated" if not self.dry_run else "Would update"
                console.print(Panel(
                    f"✅ {action} {files_changed} out of {total_files} files",
                    title="🎉 Synchronization Complete",
                    style=status
                ))
            else:
                console.print(Panel(
                    f"ℹ️  All {total_files} files are already up to date",
                    title="ℹ️  No Changes Needed",
                    style="blue"
                ))
            
            if self.errors:
                console.print(Panel(
                    f"❌ {len(self.errors)} errors occurred during synchronization",
                    title="⚠️  Errors",
                    style="red"
                ))
        else:
            if files_changed > 0:
                action = "Updated" if not self.dry_run else "Would update"
                print(f"✅ {action} {files_changed} out of {total_files} files")
            else:
                print(f"ℹ️  All {total_files} files are already up to date")
            
            if self.errors:
                print(f"❌ {len(self.errors)} errors occurred during synchronization")
    
    def check_consistency(self) -> bool:
        """Check version consistency across all files"""
        if HAS_RICH:
            panel = Panel(
                "Checking version consistency across all files",
                title="🔍 Consistency Check",
                title_align="left"
            )
            console.print(panel)
        else:
            print("Checking version consistency across all files")
        
        inconsistencies = []
        
        for file_path, rules in self.sync_rules.items():
            full_path = self.project_root / file_path
            relative_path = full_path.relative_to(self.project_root)
            
            if not full_path.exists():
                continue
            
            content = self._read_file(full_path)
            if content is None:
                continue
            
            for rule in rules:
                pattern = rule['pattern']
                version_key = rule['version_key']
                
                try:
                    expected_version = get_version(version_key)
                    matches = re.findall(pattern, content)
                    
                    for match in matches:
                        if match != expected_version:
                            inconsistencies.append({
                                'file': str(relative_path),
                                'version_key': version_key,
                                'found_version': match,
                                'expected_version': expected_version
                            })
                
                except Exception as e:
                    self._log_error(f"Error checking {version_key} in {relative_path}: {e}")
        
        if inconsistencies:
            if HAS_RICH:
                table = Table(title="🚨 Version Inconsistencies Found")
                table.add_column("File", style="cyan", no_wrap=True)
                table.add_column("Version Key", style="yellow")
                table.add_column("Found", style="red")
                table.add_column("Expected", style="green")
                
                for inconsistency in inconsistencies:
                    table.add_row(
                        inconsistency['file'],
                        inconsistency['version_key'],
                        inconsistency['found_version'],
                        inconsistency['expected_version']
                    )
                
                console.print(table)
                console.print()
                console.print(Panel(
                    f"❌ Found {len(inconsistencies)} version inconsistencies",
                    title="⚠️  Inconsistencies Detected",
                    style="red"
                ))
            else:
                print("Version inconsistencies found:")
                for inconsistency in inconsistencies:
                    print(f"  {inconsistency['file']}: {inconsistency['version_key']} "
                          f"= {inconsistency['found_version']} (expected: {inconsistency['expected_version']})")
        else:
            if HAS_RICH:
                console.print(Panel(
                    "✅ All versions are consistent across all files",
                    title="🎉 All Good!",
                    style="green"
                ))
            else:
                print("✅ All versions are consistent across all files")
        
        return len(inconsistencies) == 0

def main():
    """Main entry point"""
    parser = argparse.ArgumentParser(description='Synchronize versions from versions.yml')
    parser.add_argument('--dry-run', action='store_true', help='Show what would be changed without making changes')
    parser.add_argument('--verbose', '-v', action='store_true', help='Enable verbose output')
    parser.add_argument('--check', action='store_true', help='Check consistency without making changes')
    
    args = parser.parse_args()
    
    if HAS_RICH:
        console.print(Panel(
            "Version Synchronization Tool",
            subtitle="Keeping all project versions in sync",
            style="bold blue"
        ))
    else:
        print("Version Synchronization Tool")
        print("Keeping all project versions in sync")
    
    syncer = VersionSyncer(dry_run=args.dry_run, verbose=args.verbose)
    
    if args.check:
        success = syncer.check_consistency()
    else:
        success = syncer.sync_all_files()
    
    if not success:
        sys.exit(1)

if __name__ == '__main__':
    main() 