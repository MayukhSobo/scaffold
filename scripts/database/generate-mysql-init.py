#!/usr/bin/env python3
"""
Generate MySQL initialization script from configuration files.
This ensures the MySQL init script uses the same credentials as the application.
Supports base64 password encoding for security.
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

def encode_password_base64(password):
    """Encode password to base64."""
    return base64.b64encode(password.encode('utf-8')).decode('utf-8')

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

def generate_mysql_init_sql(config):
    """Generate MySQL initialization SQL from configuration."""
    db_config = config.get('db', {}).get('mysql', {})
    
    if not db_config:
        raise ValueError("No database configuration found in config file")
    
    required_fields = ['user', 'password', 'database']
    for field in required_fields:
        if field not in db_config:
            raise ValueError(f"Missing required database field: {field}")
    
    user = db_config['user']
    password = db_config['password']
    database = db_config['database']
    
    # Decode password if it's base64 encoded (for use in SQL)
    decoded_password = decode_password_base64(password)
    
    # Generate SQL script
    sql_script = f"""-- MySQL initialization script
-- Generated from configuration file: {os.path.basename(config.get('_config_file', 'unknown'))}
-- This script sets up the database user and permissions
-- Note: Passwords are automatically decoded from base64 if encoded

-- Create a dedicated user with root permissions
CREATE USER IF NOT EXISTS '{user}'@'%' IDENTIFIED BY '{decoded_password}';
GRANT ALL PRIVILEGES ON *.* TO '{user}'@'%' WITH GRANT OPTION;

-- Ensure the database exists
CREATE DATABASE IF NOT EXISTS {database};

-- Grant specific database permissions (redundant but explicit)
GRANT ALL PRIVILEGES ON {database}.* TO '{user}'@'%';

-- Flush privileges to apply changes
FLUSH PRIVILEGES;

-- Display confirmation
SELECT 'MySQL initialization completed successfully' AS status;
"""
    
    return sql_script

def main():
    parser = argparse.ArgumentParser(description='Generate MySQL initialization script from config files')
    parser.add_argument('config_file', help='Path to the configuration file (e.g., configs/docker.yml)')
    parser.add_argument('--output', '-o', help='Output file path (default: db/mysql-init.sql)')
    parser.add_argument('--encode-password', action='store_true', help='Encode password to base64 before processing')
    
    args = parser.parse_args()
    
    try:
        # Determine output path
        if args.output:
            output_path = args.output
        else:
            # Default to migrations directory
            script_dir = Path(__file__).parent.parent.parent  # Go up to project root
            output_path = script_dir / "db" / "mysql-init.sql"
        
        config = load_config(args.config_file)
        
        # Store config file path for reference in the SQL
        config['_config_file'] = args.config_file
        
        # Encode password if requested
        if args.encode_password and 'db' in config and 'mysql' in config['db']:
            original_password = config['db']['mysql']['password']
            encoded_password = encode_password_base64(original_password)
            config['db']['mysql']['password'] = encoded_password
            print(f"Password encoded to base64: {encoded_password}")
        
        sql_script = generate_mysql_init_sql(config)
        
        # Ensure output directory exists
        output_path = Path(output_path)
        output_path.parent.mkdir(parents=True, exist_ok=True)
        
        with open(output_path, 'w') as f:
            f.write(sql_script)
        
        print(f"MySQL initialization script written to {output_path}")
    
    except Exception as e:
        print(f"Error: {e}", file=sys.stderr)
        sys.exit(1)

if __name__ == "__main__":
    main() 