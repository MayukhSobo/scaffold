#!/usr/bin/env python3

"""
install-task.py - Install Task runner with rich formatting
Usage: python install-task.py
"""

import os
import sys
import subprocess
import platform
import shutil
import tempfile
import urllib.request
import json
from pathlib import Path
from typing import Optional, Tuple

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

try:
    import requests
    HAS_REQUESTS = True
except ImportError:
    HAS_REQUESTS = False

# Initialize console
console = Console() if HAS_RICH else None

class TaskInstaller:
    """Task runner installation manager with rich output"""
    
    def __init__(self):
        self.github_api_url = "https://api.github.com/repos/go-task/task/releases/latest"
        self.github_release_url = "https://github.com/go-task/task/releases/latest"
        
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
    
    def _run_command(self, cmd: list, description: str = "") -> Tuple[bool, str]:
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
    
    def _get_latest_version(self) -> Optional[str]:
        """Get the latest Task version from GitHub"""
        try:
            if HAS_REQUESTS:
                response = requests.get(self.github_api_url, timeout=10)
                if response.status_code == 200:
                    data = response.json()
                    return data.get('tag_name', '').lstrip('v')
            else:
                # Fallback to urllib
                with urllib.request.urlopen(self.github_api_url) as response:
                    data = json.loads(response.read().decode())
                    return data.get('tag_name', '').lstrip('v')
        except Exception as e:
            self._log_warn(f"Could not get latest version from GitHub: {e}")
            return None
    
    def _check_task_installed(self) -> Tuple[bool, Optional[str]]:
        """Check if Task is already installed and return version"""
        if not shutil.which('task'):
            return False, None
        
        success, output = self._run_command(['task', '--version'])
        if success:
            # Extract version from output like "Task version: v3.21.0"
            for line in output.split('\n'):
                if 'version:' in line.lower():
                    version = line.split(':')[-1].strip().lstrip('v')
                    return True, version
        
        return True, "unknown"
    
    def _get_download_info(self, version: str) -> Optional[Tuple[str, str]]:
        """Get download URL and filename for the current platform"""
        system = platform.system().lower()
        machine = platform.machine().lower()
        
        # Map platform names
        if system == 'darwin':
            os_name = 'darwin'
        elif system == 'linux':
            os_name = 'linux'
        elif system == 'windows':
            os_name = 'windows'
        else:
            self._log_error(f"Unsupported operating system: {system}")
            return None
        
        # Map architecture names
        if machine in ['x86_64', 'amd64']:
            arch = 'amd64'
        elif machine in ['aarch64', 'arm64']:
            arch = 'arm64'
        elif machine in ['i386', 'i686']:
            arch = '386'
        else:
            self._log_error(f"Unsupported architecture: {machine}")
            return None
        
        # Build filename
        if system == 'windows':
            filename = f"task_{os_name}_{arch}.zip"
        else:
            filename = f"task_{os_name}_{arch}.tar.gz"
        
        url = f"https://github.com/go-task/task/releases/download/v{version}/{filename}"
        
        return url, filename
    
    def _download_file(self, url: str, filename: str) -> Optional[Path]:
        """Download a file with progress bar"""
        try:
            temp_dir = Path(tempfile.gettempdir())
            file_path = temp_dir / filename
            
            if HAS_RICH:
                with console.status(f"Downloading {filename}..."):
                    if HAS_REQUESTS:
                        response = requests.get(url, stream=True, timeout=30)
                        response.raise_for_status()
                        
                        with open(file_path, 'wb') as f:
                            for chunk in response.iter_content(chunk_size=8192):
                                if chunk:
                                    f.write(chunk)
                    else:
                        urllib.request.urlretrieve(url, file_path)
            else:
                print(f"Downloading {filename}...")
                if HAS_REQUESTS:
                    response = requests.get(url, stream=True, timeout=30)
                    response.raise_for_status()
                    
                    with open(file_path, 'wb') as f:
                        for chunk in response.iter_content(chunk_size=8192):
                            if chunk:
                                f.write(chunk)
                else:
                    urllib.request.urlretrieve(url, file_path)
            
            return file_path
            
        except Exception as e:
            self._log_error(f"Failed to download {filename}: {e}")
            return None
    
    def _install_from_archive(self, archive_path: Path) -> bool:
        """Install Task from downloaded archive"""
        try:
            temp_dir = Path(tempfile.gettempdir()) / "task-install"
            temp_dir.mkdir(exist_ok=True)
            
            if HAS_RICH:
                with console.status("Extracting archive..."):
                    if archive_path.suffix == '.zip':
                        success, output = self._run_command(['unzip', '-o', str(archive_path), '-d', str(temp_dir)])
                    else:
                        success, output = self._run_command(['tar', '-xzf', str(archive_path), '-C', str(temp_dir)])
            else:
                print("Extracting archive...")
                if archive_path.suffix == '.zip':
                    success, output = self._run_command(['unzip', '-o', str(archive_path), '-d', str(temp_dir)])
                else:
                    success, output = self._run_command(['tar', '-xzf', str(archive_path), '-C', str(temp_dir)])
            
            if not success:
                self._log_error(f"Failed to extract archive: {output}")
                return False
            
            # Find the task binary
            task_binary = None
            for file in temp_dir.rglob('task*'):
                if file.is_file() and file.stat().st_mode & 0o111:  # executable
                    task_binary = file
                    break
            
            if not task_binary:
                self._log_error("Task binary not found in archive")
                return False
            
            # Install to /usr/local/bin
            install_dir = Path('/usr/local/bin')
            if not install_dir.exists():
                install_dir.mkdir(parents=True, exist_ok=True)
            
            install_path = install_dir / 'task'
            
            if HAS_RICH:
                with console.status("Installing Task..."):
                    success, output = self._run_command(['sudo', 'cp', str(task_binary), str(install_path)])
                    if success:
                        self._run_command(['sudo', 'chmod', '+x', str(install_path)])
            else:
                print("Installing Task...")
                success, output = self._run_command(['sudo', 'cp', str(task_binary), str(install_path)])
                if success:
                    self._run_command(['sudo', 'chmod', '+x', str(install_path)])
            
            # Cleanup
            shutil.rmtree(temp_dir, ignore_errors=True)
            archive_path.unlink(missing_ok=True)
            
            return success
            
        except Exception as e:
            self._log_error(f"Failed to install Task: {e}")
            return False
    
    def _install_via_package_manager(self) -> bool:
        """Install Task using system package manager"""
        system = platform.system().lower()
        
        if system == 'darwin':
            # macOS - use Homebrew
            if shutil.which('brew'):
                if HAS_RICH:
                    with console.status("Installing Task via Homebrew..."):
                        success, output = self._run_command(['brew', 'install', 'go-task/tap/go-task'])
                else:
                    print("Installing Task via Homebrew...")
                    success, output = self._run_command(['brew', 'install', 'go-task/tap/go-task'])
                
                if success:
                    self._log_info("Task installed successfully via Homebrew")
                    return True
                else:
                    self._log_error(f"Failed to install via Homebrew: {output}")
            else:
                self._log_warn("Homebrew not found. Will try direct download.")
        
        elif system == 'linux':
            # Linux - try different package managers
            if shutil.which('snap'):
                if HAS_RICH:
                    with console.status("Installing Task via snap..."):
                        success, output = self._run_command(['sudo', 'snap', 'install', 'task', '--classic'])
                else:
                    print("Installing Task via snap...")
                    success, output = self._run_command(['sudo', 'snap', 'install', 'task', '--classic'])
                
                if success:
                    self._log_info("Task installed successfully via snap")
                    return True
                else:
                    self._log_warn(f"Failed to install via snap: {output}")
            
            # Try other package managers if snap fails
            if shutil.which('apt-get'):
                if HAS_RICH:
                    with console.status("Installing Task via apt..."):
                        # Add Task repository
                        commands = [
                            ['sudo', 'sh', '-c', 'echo "deb [trusted=yes] https://repo.goreleaser.com/apt/ /" > /etc/apt/sources.list.d/goreleaser.list'],
                            ['sudo', 'apt', 'update'],
                            ['sudo', 'apt', 'install', '-y', 'task']
                        ]
                        
                        for cmd in commands:
                            success, output = self._run_command(cmd)
                            if not success:
                                self._log_warn(f"Command failed: {' '.join(cmd)}")
                                break
                        else:
                            self._log_info("Task installed successfully via apt")
                            return True
                else:
                    print("Installing Task via apt...")
                    # Similar implementation without rich
                    commands = [
                        ['sudo', 'sh', '-c', 'echo "deb [trusted=yes] https://repo.goreleaser.com/apt/ /" > /etc/apt/sources.list.d/goreleaser.list'],
                        ['sudo', 'apt', 'update'],
                        ['sudo', 'apt', 'install', '-y', 'task']
                    ]
                    
                    for cmd in commands:
                        success, output = self._run_command(cmd)
                        if not success:
                            break
                    else:
                        self._log_info("Task installed successfully via apt")
                        return True
        
        return False
    
    def install(self) -> bool:
        """Install Task runner"""
        if HAS_RICH:
            panel = Panel(
                "Installing Task Runner",
                subtitle="Cross-platform task runner and build tool",
                style="bold blue"
            )
            console.print(panel)
        else:
            print("Installing Task Runner")
            print("Cross-platform task runner and build tool")
        
        # Check if already installed
        is_installed, current_version = self._check_task_installed()
        
        if is_installed:
            latest_version = self._get_latest_version()
            
            if latest_version and current_version and current_version != "unknown":
                if current_version == latest_version:
                    self._log_info(f"Task is already installed and up to date (v{current_version})")
                    return True
                else:
                    if HAS_RICH:
                        table = Table(title="Task Version Information")
                        table.add_column("Status", style="yellow")
                        table.add_column("Version", style="white")
                        table.add_row("Current", f"v{current_version}")
                        table.add_row("Latest", f"v{latest_version}")
                        console.print(table)
                        console.print()
                    else:
                        print(f"Current version: v{current_version}")
                        print(f"Latest version: v{latest_version}")
                        print("Updating Task...")
            else:
                self._log_info(f"Task is already installed (v{current_version})")
                return True
        
        # Try package manager first
        if self._install_via_package_manager():
            return True
        
        # Fall back to direct download
        latest_version = self._get_latest_version()
        if not latest_version:
            self._log_error("Could not determine latest Task version")
            return False
        
        download_info = self._get_download_info(latest_version)
        if not download_info:
            return False
        
        url, filename = download_info
        
        if HAS_RICH:
            console.print(f"📥 Downloading Task v{latest_version} from GitHub...")
        else:
            print(f"Downloading Task v{latest_version} from GitHub...")
        
        archive_path = self._download_file(url, filename)
        if not archive_path:
            return False
        
        if self._install_from_archive(archive_path):
            self._log_info(f"Task v{latest_version} installed successfully!")
            
            # Verify installation
            is_installed, installed_version = self._check_task_installed()
            if is_installed:
                if HAS_RICH:
                    console.print(Panel(
                        f"✅ Task v{installed_version} is ready to use!",
                        title="🎉 Installation Complete",
                        style="green"
                    ))
                else:
                    print(f"✅ Task v{installed_version} is ready to use!")
                return True
            else:
                self._log_error("Task installation verification failed")
                return False
        else:
            return False

def main():
    """Main entry point"""
    installer = TaskInstaller()
    
    success = installer.install()
    
    if not success:
        if HAS_RICH:
            console.print(Panel(
                "Installation failed. Please check the errors above and try again.",
                title="❌ Installation Failed",
                style="red"
            ))
        else:
            print("❌ Installation failed. Please check the errors above and try again.")
        sys.exit(1)

if __name__ == '__main__':
    main() 