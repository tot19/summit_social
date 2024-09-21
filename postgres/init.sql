CREATE TABLE IF NOT EXISTS posts (
    id SERIAL PRIMARY KEY,
    date TIMESTAMP DEFAULT NOW(),
    content TEXT,
    image VARCHAR(255),
    likes INTEGER DEFAULT 0,
    poster VARCHAR(50),
    comments_count INTEGER DEFAULT 0
);
