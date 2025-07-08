#!/bin/bash

# Custom MySQL entrypoint that handles base64 encoded passwords
# This script decodes base64 encoded environment variables before starting MySQL

set -e

# Function to decode base64 if the value looks like base64
decode_if_base64() {
    local value="$1"
    # Simple check if it looks like base64 (alphanumeric + / + =)
    if [[ "$value" =~ ^[A-Za-z0-9+/]+=*$ ]] && [ ${#value} -gt 8 ]; then
        # Try to decode, if it fails, use original value
        decoded=$(echo "$value" | base64 -d 2>/dev/null) || decoded="$value"
        echo "$decoded"
    else
        echo "$value"
    fi
}

# Function to log with timestamp
log_info() {
    echo "[$(date '+%Y-%m-%d %H:%M:%S')] [mysql-entrypoint] INFO: $1"
}

log_info "Starting MySQL with custom entrypoint for base64 password support"

# Decode base64 encoded password if present
if [ -n "$MYSQL_ROOT_PASSWORD" ]; then
    ORIGINAL_ROOT_PASSWORD="$MYSQL_ROOT_PASSWORD"
    export MYSQL_ROOT_PASSWORD=$(decode_if_base64 "$MYSQL_ROOT_PASSWORD")
    if [ "$MYSQL_ROOT_PASSWORD" != "$ORIGINAL_ROOT_PASSWORD" ]; then
        log_info "🔓 Decoded base64 MySQL root password"
    else
        log_info "🔓 Using plain text MySQL root password"
    fi
fi

if [ -n "$MYSQL_PASSWORD" ]; then
    ORIGINAL_PASSWORD="$MYSQL_PASSWORD"
    export MYSQL_PASSWORD=$(decode_if_base64 "$MYSQL_PASSWORD")
    if [ "$MYSQL_PASSWORD" != "$ORIGINAL_PASSWORD" ]; then
        log_info "🔓 Decoded base64 MySQL user password"
    else
        log_info "🔓 Using plain text MySQL user password"
    fi
fi

# Log MySQL configuration (without showing passwords)
log_info "MySQL Configuration:"
log_info "  - Database: ${MYSQL_DATABASE:-default}"
log_info "  - User: ${MYSQL_USER:-root}"
log_info "  - Port: 3306"
log_info "  - Data Directory: /var/lib/mysql"

# Execute the original MySQL entrypoint
log_info "Starting MySQL server..."
exec docker-entrypoint.sh "$@" 