#!/usr/bin/env python3
"""
Generate Docker Compose configuration from config files.
This ensures the Docker service name matches the host field in the config.
Supports base64 password encoding for enhanced security.
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

def should_use_base64_passwords(config):
    """Determine if we should use base64 encoded passwords."""
    # Check if passwords are already base64 encoded
    password = config.get('db', {}).get('mysql', {}).get('password', '')
    # Ensure password is a string before regex
    password_str = str(password)
    import re
    return re.match(r'^[A-Za-z0-9+/]+=*$', password_str) and len(password_str) > 8

def generate_docker_compose(config, use_base64=False):
    """Generate Docker Compose configuration from config."""
    db_config = config.get('db', {}).get('mysql', {})
    adminer_config = config.get('db', {}).get('adminer', {})
    http_config = config.get('http', {})
    
    if not db_config:
        raise ValueError("No database configuration found in config file")
    
    required_fields = ['user', 'password', 'database']
    for field in required_fields:
        if field not in db_config:
            raise ValueError(f"Missing required database field: {field}")
    
    # Get the service name from the host field
    service_name = db_config.get('host', 'mysql')
    container_name = f"scaffold-{service_name}"
    
    # Get adminer configuration
    adminer_host = adminer_config.get('host', 'adminer')
    adminer_port = adminer_config.get('port', 8080)
    adminer_theme = adminer_config.get('theme', 'default')
    adminer_container_name = f"scaffold-{adminer_host}"
    
    # Get scaffold app configuration
    scaffold_host = http_config.get('host', 'scaffold')
    scaffold_port = http_config.get('port', 8000)
    scaffold_container_name = f"scaffold-{scaffold_host}"
    
    # Handle password encoding
    mysql_password = db_config['password']
    if use_base64 and not should_use_base64_passwords(config):
        mysql_password = encode_password_base64(mysql_password)
    
    # Generate the docker-compose configuration
    compose_config = {
        'services': {
            service_name: {
                'image': 'mysql:8.0',
                'container_name': container_name,
                'environment': {
                    'MYSQL_ROOT_PASSWORD': mysql_password,
                    'MYSQL_DATABASE': db_config['database'],
                    'MYSQL_USER': db_config['user'],
                    'MYSQL_PASSWORD': mysql_password
                },
                'ports': [f"{db_config.get('port', 3306)}:3306"],
                'volumes': [
                    'mysql_data:/var/lib/mysql',
                    './db/mysql-init.sql:/docker-entrypoint-initdb.d/mysql-init.sql:ro'
                ]
            },
            adminer_host: {
                'image': 'adminer:latest',
                'container_name': adminer_container_name,
                'ports': [f"{adminer_port}:8080"],
                'depends_on': [service_name],
                'environment': {
                    'ADMINER_DEFAULT_SERVER': service_name,
                    'ADMINER_DESIGN': adminer_theme
                }
            },
            scaffold_host: {
                'build': '.',
                'container_name': scaffold_container_name,
                'ports': [f"{scaffold_port}:8000"],
                'depends_on': [service_name],
                'volumes': [
                    './configs:/app/configs',
                    './logs:/app/logs'
                ],
                'command': ['./server', '--config=configs/docker.yml']
            }
        },
        'volumes': {
            'mysql_data': None
        }
    }
    
    # Add custom entrypoint if using base64 passwords
    if use_base64 or should_use_base64_passwords(config):
        compose_config['services'][service_name]['entrypoint'] = ['./scripts/docker/mysql-entrypoint.sh', 'mysqld']
        compose_config['services'][service_name]['volumes'].append(
            './scripts/docker/mysql-entrypoint.sh:/scripts/docker/mysql-entrypoint.sh:ro'
        )
    
    return compose_config

def main():
    parser = argparse.ArgumentParser(description='Generate Docker Compose configuration from config files')
    parser.add_argument('config_file', help='Path to the configuration file (e.g., configs/docker.yml)')
    parser.add_argument('--output', '-o', help='Output file (default: docker-compose.yml)', default='docker-compose.yml')
    parser.add_argument('--use-base64', action='store_true', help='Use base64 encoded passwords')
    
    args = parser.parse_args()
    
    try:
        config = load_config(args.config_file)
        compose_config = generate_docker_compose(config, use_base64=args.use_base64)
        
        # Write the docker-compose.yml file
        with open(args.output, 'w') as file:
            yaml.dump(compose_config, file, default_flow_style=False, sort_keys=False)
        
        service_name = config['db']['mysql'].get('host', 'mysql')
        password_method = "base64 encoded" if (args.use_base64 or should_use_base64_passwords(config)) else "plain text"
        print(f"Generated {args.output} with service name: {service_name} (passwords: {password_method})")
    
    except Exception as e:
        print(f"Error: {e}", file=sys.stderr)
        sys.exit(1)

if __name__ == "__main__":
    main() 