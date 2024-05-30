CREATE TABLE threads (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    forum_id INTEGER,
    title TEXT NOT NULL,
    user_id INTEGER,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (forum_id) REFERENCES forums(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);
