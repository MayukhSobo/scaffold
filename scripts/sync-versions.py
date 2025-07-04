#!/usr/bin/env python3

"""
sync-versions.py - Synchronize versions from versions.yml across all project files
Usage: python sync-versions.py [--dry-run] [--verbose]
"""

import re
import argparse
from pathlib import Path
from typing import List, Dict, Optional, Tuple, Set

# Import common utilities
from common import ScriptBase, has_rich, get_console

class VersionSyncer(ScriptBase):
    """Version synchronization manager with rich output"""
    
    def __init__(self, dry_run: bool = False, verbose: bool = False):
        super().__init__("VersionSyncer", dry_run)
        self.verbose = verbose
        self.versions = self.version_helper.load_versions_file()
        self.changes_made = []
        
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
    

    
    def _sync_file(self, file_path: Path, rules: List[Dict]) -> bool:
        """Sync versions in a single file"""
        relative_path = file_path.relative_to(self.project_root)
        
        if not file_path.exists():
            self.logger.warn(f"File not found: {relative_path}")
            return False
        
        content = self.file_manager.read_file(file_path)
        if content is None:
            return False
        
        original_content = content
        changes_in_file = []
        
        for rule in rules:
            pattern = rule['pattern']
            version_key = rule['version_key']
            replacement_func = rule['replacement']
            
            try:
                expected_version = self.version_helper.get_version(version_key)
                
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
                self.logger.error(f"Error processing {version_key} in {relative_path}: {e}")
                continue
        
        if changes_in_file:
            self.file_manager.write_file(file_path, content)
            self.changes_made.extend([{
                'file': str(relative_path),
                'changes': changes_in_file
            }])
            return True
        
        return False
    
    def sync_all_files(self) -> bool:
        """Sync versions across all configured files"""
        self.rich.print_panel(
            "Synchronizing versions from versions.yml",
            title="🔄 Version Synchronization"
        )
        
        if self.dry_run:
            self.rich.print_dry_run_warning()
        
        total_files = len(self.sync_rules)
        files_changed = 0
        
        if has_rich():
            console = get_console()
            from rich.progress import Progress, SpinnerColumn, TextColumn, BarColumn, TaskProgressColumn
            
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
                self.logger.verbose(f"Syncing {file_path}...", self.verbose)
                
                if self._sync_file(full_path, rules):
                    files_changed += 1
        
        self._show_summary(files_changed, total_files)
        return not self.has_errors()
    
    def _show_summary(self, files_changed: int, total_files: int):
        """Show synchronization summary"""
        if has_rich():
            console = get_console()
            from rich.table import Table
            from rich.panel import Panel
            
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
            
            if self.has_errors():
                console.print(Panel(
                    f"❌ {len(self.get_errors())} errors occurred during synchronization",
                    title="⚠️  Errors",
                    style="red"
                ))
        else:
            if files_changed > 0:
                action = "Updated" if not self.dry_run else "Would update"
                print(f"✅ {action} {files_changed} out of {total_files} files")
            else:
                print(f"ℹ️  All {total_files} files are already up to date")
            
            if self.has_errors():
                print(f"❌ {len(self.get_errors())} errors occurred during synchronization")
    
    def check_consistency(self) -> bool:
        """Check version consistency across all files"""
        self.rich.print_panel(
            "Checking version consistency across all files",
            title="🔍 Consistency Check"
        )
        
        inconsistencies = []
        
        for file_path, rules in self.sync_rules.items():
            full_path = self.project_root / file_path
            relative_path = full_path.relative_to(self.project_root)
            
            if not full_path.exists():
                continue
            
            content = self.file_manager.read_file(full_path)
            if content is None:
                continue
            
            for rule in rules:
                pattern = rule['pattern']
                version_key = rule['version_key']
                
                try:
                    expected_version = self.version_helper.get_version(version_key)
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
                    self.logger.error(f"Error checking {version_key} in {relative_path}: {e}")
        
        if inconsistencies:
            if has_rich():
                console = get_console()
                from rich.table import Table
                from rich.panel import Panel
                
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
            self.rich.print_panel(
                "✅ All versions are consistent across all files",
                title="🎉 All Good!",
                style="green"
            )
        
        return len(inconsistencies) == 0

def main():
    """Main entry point"""
    import sys
    
    parser = argparse.ArgumentParser(description='Synchronize versions from versions.yml')
    parser.add_argument('--dry-run', action='store_true', help='Show what would be changed without making changes')
    parser.add_argument('--verbose', '-v', action='store_true', help='Enable verbose output')
    parser.add_argument('--check', action='store_true', help='Check consistency without making changes')
    
    args = parser.parse_args()
    
    # Create syncer instance
    syncer = VersionSyncer(dry_run=args.dry_run, verbose=args.verbose)
    
    # Show header
    syncer.rich.print_panel(
        "Version Synchronization Tool",
        title="Keeping all project versions in sync",
        style="bold blue"
    )
    
    if args.check:
        success = syncer.check_consistency()
    else:
        success = syncer.sync_all_files()
    
    if not success:
        syncer.exit_with_error()
    else:
        syncer.exit_with_success()

if __name__ == '__main__':
    main() 