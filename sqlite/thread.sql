CREATE TABLE threads (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    categorie_id TEXT,
    title TEXT NOT NULL,
    user_username TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (categorie_id) REFERENCES categories(title),
    FOREIGN KEY (user_username) REFERENCES users(username)
);