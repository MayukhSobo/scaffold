#!/usr/bin/env python3

"""
run.py - Script runner that ensures virtual environment is used
This script automatically uses the .venv Python interpreter for all development scripts
"""

import os
import sys
import subprocess
from pathlib import Path

def get_project_root() -> Path:
    """Get the project root directory"""
    return Path(__file__).parent.parent

def get_venv_python() -> Path:
    """Get the path to the virtual environment Python executable"""
    project_root = get_project_root()
    venv_dir = project_root / ".venv"
    
    if os.name == 'nt':  # Windows
        return venv_dir / "Scripts" / "python.exe"
    else:  # Unix-like (macOS, Linux)
        return venv_dir / "bin" / "python"

def check_venv_exists() -> bool:
    """Check if virtual environment exists"""
    venv_python = get_venv_python()
    return venv_python.exists()

def ensure_venv_setup():
    """Ensure virtual environment is set up"""
    if not check_venv_exists():
        print("❌ Python virtual environment not found!")
        print("Please set it up first by running:")
        print("  task python")
        print("  # or")
        print("  task shared:setup:python")
        sys.exit(1)

def run_script_with_venv(script_name: str, args: list = None) -> int:
    """Run a Python script using the virtual environment"""
    ensure_venv_setup()
    
    venv_python = get_venv_python()
    scripts_dir = Path(__file__).parent
    script_path = scripts_dir / script_name
    
    if not script_path.exists():
        print(f"❌ Script not found: {script_path}")
        sys.exit(1)
    
    # Build command
    cmd = [str(venv_python), str(script_path)]
    if args:
        cmd.extend(args)
    
    # Run the script
    try:
        result = subprocess.run(cmd, cwd=get_project_root())
        return result.returncode
    except KeyboardInterrupt:
        print("\n🛑 Interrupted by user")
        return 130
    except Exception as e:
        print(f"❌ Error running script: {e}")
        return 1

def main():
    """Main entry point"""
    if len(sys.argv) < 2:
        print("Usage: python run.py <script_name> [args...]")
        print("Example: python run.py version-helper.py list")
        sys.exit(1)
    
    script_name = sys.argv[1]
    args = sys.argv[2:] if len(sys.argv) > 2 else []
    
    # Ensure script has .py extension
    if not script_name.endswith('.py'):
        script_name += '.py'
    
    return run_script_with_venv(script_name, args)

if __name__ == '__main__':
    sys.exit(main()) 