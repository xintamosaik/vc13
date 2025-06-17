CREATE TABLE annotations (
    id INTEGER PRIMARY KEY,
    filename text NOT NULL,
    keyword text NOT NULL,
    content text NOT NULL,
    notes text,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);