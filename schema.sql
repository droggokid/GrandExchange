-- Create database (run this first if database doesn't exist)
-- CREATE DATABASE pagination;

-- Create the search_items table
CREATE TABLE IF NOT EXISTS search_items (
    id INTEGER PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    type VARCHAR(100),
    type_icon VARCHAR(255),
    icon VARCHAR(255),
    icon_large VARCHAR(255),
    members VARCHAR(10),
    current_trend VARCHAR(50),
    current_price TEXT,
    today_trend VARCHAR(50),
    today_price TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Create index on name for faster searches
CREATE INDEX IF NOT EXISTS idx_search_items_name ON search_items(name);
