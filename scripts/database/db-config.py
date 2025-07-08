#!/usr/bin/env python3
"""
Database configuration helper script for migration tasks.
Reads database configuration from YAML files and generates connection strings.
"""

import yaml
import sys
import os
import argparse
import base64
import re
from pathlib import Path

def decode_if_base64(value):
    """Decode base64 if the value looks like base64."""
    if not isinstance(value, str):
        return value
    
    # Simple check if it looks like base64 (alphanumeric + / + =)
    if re.match(r'^[A-Za-z0-9+/]+=*$', value) and len(value) > 8:
        try:
            # Try to decode, if it fails, use original value
            decoded = base64.b64decode(value).decode('utf-8')
            return decoded
        except Exception:
            return value
    else:
        return value

def load_config(config_path):
    """Load and parse the YAML configuration file."""
    if not os.path.exists(config_path):
        raise FileNotFoundError(f"Configuration file not found: {config_path}")
    
    with open(config_path, 'r') as file:
        config = yaml.safe_load(file)
    
    return config

def build_mysql_dsn(config):
    """Build MySQL DSN from configuration."""
    db_config = config.get('db', {}).get('mysql', {})
    
    if not db_config:
        raise ValueError("No database configuration found in config file")
    
    required_fields = ['host', 'port', 'user', 'password', 'database']
    for field in required_fields:
        if field not in db_config:
            raise ValueError(f"Missing required database field: {field}")
    
    host = db_config['host']
    port = db_config['port']
    user = db_config['user']
    password = decode_if_base64(db_config['password'])
    database = db_config['database']
    
    # Build MySQL DSN
    dsn = f"{user}:{password}@tcp({host}:{port})/{database}?charset=utf8mb4&parseTime=True&loc=Local"
    return dsn

def main():
    parser = argparse.ArgumentParser(description='Generate database connection strings from config files')
    parser.add_argument('config_file', help='Path to the configuration file (e.g., configs/local.yml)')
    parser.add_argument('--format', choices=['dsn', 'env'], default='dsn', 
                       help='Output format: dsn (connection string) or env (environment variables)')
    
    args = parser.parse_args()
    
    try:
        config = load_config(args.config_file)
        
        if args.format == 'dsn':
            dsn = build_mysql_dsn(config)
            print(dsn)
        elif args.format == 'env':
            db_config = config.get('db', {}).get('mysql', {})
            print(f"DB_HOST={db_config.get('host', '')}")
            print(f"DB_PORT={db_config.get('port', '')}")
            print(f"DB_USER={db_config.get('user', '')}")
            print(f"DB_PASSWORD={decode_if_base64(db_config.get('password', ''))}")
            print(f"DB_DATABASE={db_config.get('database', '')}")
    
    except Exception as e:
        print(f"Error: {e}", file=sys.stderr)
        sys.exit(1)

if __name__ == "__main__":
    main() 