CREATE TABLE users (
    id VARCHAR(255) PRIMARY KEY UNIQUE NOT NULL,
    user_name VARCHAR(255),
    category VARCHAR(255),
    phone_number VARCHAR(255) UNIQUE NOT NULL,
    created_at DATETIME,
    updated_at DATETIME
);