CREATE TABLE passwords(
    id VARCHAR(255) PRIMARY KEY UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    salt VARCHAR(255) NOT NULL,
	created_at timestamp,
    updated_at timestamp
);