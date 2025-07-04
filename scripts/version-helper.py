#!/usr/bin/env python3

"""
version-helper.py - Helper script to parse versions.yml and provide version information
Usage: 
  ./scripts/version-helper.py get go                    # Get specific version
  ./scripts/version-helper.py list                     # List all versions
  ./scripts/version-helper.py env                      # Export all versions as environment variables
  ./scripts/version-helper.py load                     # Load commonly used versions (for sourcing)
"""

import os
import sys
import argparse
import re
from pathlib import Path
from typing import Dict, Any, Optional

# Try to import yaml, fall back to simple parser if not available
try:
    import yaml
    HAS_YAML = True
except ImportError:
    HAS_YAML = False

# Colors for output
class Colors:
    RED = '\033[0;31m'
    GREEN = '\033[0;32m'
    YELLOW = '\033[1;33m'
    BLUE = '\033[0;34m'
    NC = '\033[0m'  # No Color

def log_info(message: str) -> None:
    """Log info message to stderr"""
    print(f"{Colors.GREEN}[INFO]{Colors.NC} {message}", file=sys.stderr)

def log_warn(message: str) -> None:
    """Log warning message to stderr"""
    print(f"{Colors.YELLOW}[WARN]{Colors.NC} {message}", file=sys.stderr)

def log_error(message: str) -> None:
    """Log error message to stderr"""
    print(f"{Colors.RED}[ERROR]{Colors.NC} {message}", file=sys.stderr)

def get_versions_file_path() -> Path:
    """Get the path to the versions.yml file"""
    script_dir = Path(__file__).parent
    return script_dir.parent / "versions.yml"

def simple_yaml_parser(content: str) -> Dict[str, Any]:
    """Simple YAML parser for our specific format"""
    result = {}
    current_section = None
    
    for line in content.split('\n'):
        original_line = line
        line = line.strip()
        
        # Skip empty lines and comments
        if not line or line.startswith('#'):
            continue
        
        # Check for section headers (key: with no value on non-indented lines)
        if line.endswith(':') and not original_line.startswith(' ') and not original_line.startswith('\t'):
            current_section = line[:-1].strip()
            result[current_section] = {}
            continue
        
        # Check for key-value pairs
        if ':' in line:
            # Handle indented items (part of a section)
            if original_line.startswith(' ') or original_line.startswith('\t'):
                if current_section is None:
                    continue
                    
                # Remove leading whitespace
                line = line.lstrip()
                key, value = line.split(':', 1)
                key = key.strip()
                value = value.strip()
                
                # Remove quotes from value
                if value.startswith('"') and value.endswith('"'):
                    value = value[1:-1]
                elif value.startswith("'") and value.endswith("'"):
                    value = value[1:-1]
                
                result[current_section][key] = value
            else:
                # Top-level key-value pair (not part of a section)
                key, value = line.split(':', 1)
                key = key.strip()
                value = value.strip()
                
                # Remove quotes from value
                if value.startswith('"') and value.endswith('"'):
                    value = value[1:-1]
                elif value.startswith("'") and value.endswith("'"):
                    value = value[1:-1]
                
                result[key] = value
                current_section = None  # Reset current section
    
    return result

def load_versions_file() -> Dict[str, Any]:
    """Load and parse the versions.yml file"""
    versions_file = get_versions_file_path()
    
    if not versions_file.exists():
        log_error(f"versions.yml not found at {versions_file}")
        sys.exit(1)
    
    try:
        with open(versions_file, 'r', encoding='utf-8') as f:
            content = f.read()
            
        if HAS_YAML:
            return yaml.safe_load(content) or {}
        else:
            # Use simple parser
            return simple_yaml_parser(content)
            
    except Exception as e:
        log_error(f"Error parsing versions.yml: {e}")
        sys.exit(1)

def get_version(key: str) -> Optional[str]:
    """Get version for a specific key"""
    versions = load_versions_file()
    
    # Handle different key formats
    if '.' in key:
        # Key with section (e.g., "tools.golangci-lint")
        section, subkey = key.split('.', 1)
        if section in versions and isinstance(versions[section], dict):
            return versions[section].get(subkey)
    else:
        # Direct key (e.g., "go")
        return versions.get(key)
    
    return None

def list_versions() -> None:
    """List all versions"""
    versions = load_versions_file()
    
    print(f"{Colors.BLUE}Available versions:{Colors.NC}")
    
    for key, value in versions.items():
        if isinstance(value, dict):
            print(f"\n{Colors.YELLOW}[{key}]{Colors.NC}")
            for subkey, subvalue in value.items():
                print(f"  {key}.{subkey}: {subvalue}")
        else:
            print(f"  {key}: {value}")

def export_versions_env() -> None:
    """Export all versions as environment variables"""
    versions = load_versions_file()
    
    def export_var(name: str, value: str) -> None:
        # Convert to uppercase and replace hyphens with underscores
        env_name = name.upper().replace('-', '_')
        print(f"export {env_name}_VERSION='{value}'")
    
    for key, value in versions.items():
        if isinstance(value, dict):
            for subkey, subvalue in value.items():
                export_var(f"{key}_{subkey}", str(subvalue))
        else:
            export_var(key, str(value))

def load_common_versions() -> None:
    """Load commonly used versions for sourcing"""
    versions = load_versions_file()
    
    # Define commonly used versions mapping
    common_versions = {
        'GO_VERSION': get_version('go'),
        'GOLANGCI_LINT_VERSION': get_version('tools.golangci-lint'),
        'GOTESTSUM_VERSION': get_version('tools.gotestsum'),
        'GOSEC_VERSION': get_version('tools.gosec'),
        'GOVULNCHECK_VERSION': get_version('tools.govulncheck'),
        'AIR_VERSION': get_version('tools.air'),
        'GOCOV_VERSION': get_version('tools.gocov'),
        'GOCOV_HTML_VERSION': get_version('tools.gocov-html'),
        'GO_COVER_TREEMAP_VERSION': get_version('tools.go-cover-treemap'),
        'TRIVY_VERSION': get_version('tools.trivy'),
        'CODEQL_CLI_VERSION': get_version('security.codeql-cli'),
        'TASK_VERSION': get_version('build.task'),
    }
    
    # Export variables
    for var_name, version in common_versions.items():
        if version:
            print(f"export {var_name}='{version}'")

def main() -> None:
    """Main function"""
    parser = argparse.ArgumentParser(
        description='Helper script to parse versions.yml and provide version information',
        formatter_class=argparse.RawDescriptionHelpFormatter
    )
    
    subparsers = parser.add_subparsers(dest='command', help='Available commands')
    
    # Get command
    get_parser = subparsers.add_parser('get', help='Get version for specific key')
    get_parser.add_argument('key', help='Version key (e.g., "go" or "tools.golangci-lint")')
    
    # List command
    subparsers.add_parser('list', help='List all versions')
    
    # Env command
    subparsers.add_parser('env', help='Export all versions as environment variables')
    
    # Load command
    subparsers.add_parser('load', help='Load commonly used versions (for sourcing)')
    
    args = parser.parse_args()
    
    if args.command == 'get':
        version = get_version(args.key)
        if version:
            print(version)
        else:
            log_error(f"Version not found for key: {args.key}")
            sys.exit(1)
    
    elif args.command == 'list':
        list_versions()
    
    elif args.command == 'env':
        export_versions_env()
    
    elif args.command == 'load':
        load_common_versions()
    
    else:
        parser.print_help()
        print("\nExamples:")
        print("  ./scripts/version-helper.py get go")
        print("  ./scripts/version-helper.py get tools.golangci-lint")
        print("  ./scripts/version-helper.py list")
        print("  eval $(./scripts/version-helper.py load)")

if __name__ == '__main__':
    main() 