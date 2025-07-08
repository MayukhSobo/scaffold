-- MySQL initialization script
-- Generated from configuration file: docker.yml
-- This script sets up the database user and permissions
-- Note: Passwords are automatically decoded from base64 if encoded

-- Create a dedicated user with root permissions
CREATE USER IF NOT EXISTS 'scaffold'@'%' IDENTIFIED BY 'my_secure_password_123';
GRANT ALL PRIVILEGES ON *.* TO 'scaffold'@'%' WITH GRANT OPTION;

-- Ensure the database exists
CREATE DATABASE IF NOT EXISTS user;

-- Grant specific database permissions (redundant but explicit)
GRANT ALL PRIVILEGES ON user.* TO 'scaffold'@'%';

-- Flush privileges to apply changes
FLUSH PRIVILEGES;

-- Display confirmation
SELECT 'MySQL initialization completed successfully' AS status;
