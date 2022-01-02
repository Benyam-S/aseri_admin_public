create table referrals(
    code VARCHAR(50) PRIMARY KEY UNIQUE NOT NULL,
    owner VARCHAR(50),
    created_at DATETIME,
    updated_at DATETIME
);