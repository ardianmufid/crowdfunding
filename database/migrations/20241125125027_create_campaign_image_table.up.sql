CREATE TABLE IF NOT EXISTS campaign_images (
    id SERIAL PRIMARY KEY,
    campaign_id INT,
    file_name VARCHAR(255),
    is_primary SMALLINT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
) 