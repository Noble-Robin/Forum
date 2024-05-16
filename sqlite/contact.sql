CREATE TABLE contacts (
    ID TEXT PRIMARY KEY,
    user_id INTEGER,
    status TEXT,
    FOREIGN KEY (user_id) REFERENCES users(ID)
);