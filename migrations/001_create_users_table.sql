-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(255) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(255),
    last_name VARCHAR(255),
    avatar_url VARCHAR(255),
    bio TEXT,
    phone_number VARCHAR(50),
    address_street VARCHAR(255),
    address_city VARCHAR(100),
    address_state VARCHAR(100),
    address_postal_code VARCHAR(50),
    address_country VARCHAR(100),
    status ENUM('active', 'inactive', 'suspended', 'pending_verification') NOT NULL DEFAULT 'pending_verification',
    role ENUM('user', 'admin', 'moderator') NOT NULL DEFAULT 'user',
    email_verified_at TIMESTAMP NULL,
    last_login_at TIMESTAMP NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    INDEX idx_users_username (username),
    INDEX idx_users_email (email),
    INDEX idx_users_deleted_at (deleted_at),
    INDEX idx_users_status (status),
    INDEX idx_users_role (role)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd 