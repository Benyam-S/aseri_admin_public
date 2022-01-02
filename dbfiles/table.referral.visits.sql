create table referral_visits(
    id INTEGER  PRIMARY KEY UNIQUE AUTO_INCREMENT NOT NULL,
    code VARCHAR(50) NOT NULL,
    created_at DATETIME,
    FOREIGN KEY (code) REFERENCES referrals(code)
);