#!/usr/bin/env python3
"""
Generate Docker Compose environment variables from configuration files.
This ensures consistency between application config and Docker Compose setup.
Supports base64 password decoding for security.
"""

import yaml
import sys
import os
import argparse
import base64
from pathlib import Path

def load_config(config_path):
    """Load and parse the YAML configuration file."""
    if not os.path.exists(config_path):
        raise FileNotFoundError(f"Configuration file not found: {config_path}")
    
    with open(config_path, 'r') as file:
        config = yaml.safe_load(file)
    
    return config

def decode_password_base64(encoded_password):
    """Decode base64 password. Returns original if not valid base64."""
    try:
        # Ensure input is a string
        password_str = str(encoded_password)
        # Simple check if it looks like base64 (alphanumeric + / + =)
        import re
        if re.match(r'^[A-Za-z0-9+/]+=*$', password_str) and len(password_str) > 8:
            decoded = base64.b64decode(password_str).decode('utf-8')
            return decoded
        return password_str
    except Exception:
        return str(encoded_password)

def generate_env_vars(config, decode_passwords=True):
    """Generate environment variables from configuration."""
    db_config = config.get('db', {}).get('mysql', {})
    adminer_config = config.get('db', {}).get('adminer', {})
    http_config = config.get('http', {})
    
    if not db_config:
        raise ValueError("No database configuration found in config file")
    
    required_fields = ['user', 'password', 'database']
    for field in required_fields:
        if field not in db_config:
            raise ValueError(f"Missing required database field: {field}")
    
    # Decode passwords if requested (useful for shell access)
    password = db_config['password']
    if decode_passwords:
        password = decode_password_base64(password)
    
    # Generate environment variables for Docker Compose
    env_vars = {
        # MySQL variables
        'MYSQL_ROOT_PASSWORD': password,
        'MYSQL_DATABASE': db_config['database'],
        'MYSQL_USER': db_config['user'],
        'MYSQL_PASSWORD': password,
        'DB_HOST': db_config.get('host', 'mysql'),
        'DB_PORT': str(db_config.get('port', 3306)),
        'DB_USER': db_config['user'],
        'DB_PASSWORD': password,
        'DB_DATABASE': db_config['database'],
        
        # HTTP/App variables
        'HTTP_HOST': http_config.get('host', 'scaffold'),
        'HTTP_PORT': str(http_config.get('port', 8000)),
        
        # Adminer variables
        'ADMINER_HOST': adminer_config.get('host', 'adminer'),
        'ADMINER_PORT': str(adminer_config.get('port', 8080)),
        'ADMINER_THEME': adminer_config.get('theme', 'default')
    }
    
    return env_vars

def main():
    parser = argparse.ArgumentParser(description='Generate Docker Compose environment variables from config files')
    parser.add_argument('config_file', help='Path to the configuration file (e.g., configs/docker.yml)')
    parser.add_argument('--format', choices=['env', 'export', 'compose'], default='env', 
                       help='Output format: env (KEY=VALUE), export (export KEY=VALUE), compose (for .env file)')
    parser.add_argument('--no-decode', action='store_true', help='Do not decode base64 passwords (keep them encoded)')
    
    args = parser.parse_args()
    
    try:
        config = load_config(args.config_file)
        env_vars = generate_env_vars(config, decode_passwords=not args.no_decode)
        
        if args.format == 'env':
            for key, value in env_vars.items():
                print(f"{key}={value}")
        elif args.format == 'export':
            for key, value in env_vars.items():
                print(f"export {key}={value}")
        elif args.format == 'compose':
            print("# Generated environment variables for Docker Compose")
            print(f"# Source: {args.config_file}")
            for key, value in env_vars.items():
                print(f"{key}={value}")
    
    except Exception as e:
        print(f"Error: {e}", file=sys.stderr)
        sys.exit(1)

if __name__ == "__main__":
    main() 