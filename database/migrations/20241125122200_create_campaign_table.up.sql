CREATE TABLE IF NOT EXISTS campaigns (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    short_description VARCHAR(100),
    description TEXT,
    goal_amount INT,
    current_amount INT,
    perks TEXT,
    becker_count INT,
    user_id INT,
    slug VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
)