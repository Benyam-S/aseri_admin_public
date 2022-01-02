Create table server_sessions(
    user_id VARCHAR(255) NOT NULL,
    session_id VARCHAR(255) UNIQUE NOT NULL,
    ip_address VARCHAR(255),
    device_info VARCHAR(255)
);