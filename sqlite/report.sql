CREATE TABLE reports (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    reporter_username TEXT,
    thread_id INTEGER,
    reported_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (reporter_username) REFERENCES users(username),
    FOREIGN KEY (thread_id) REFERENCES threads(id)
);