CREATE TABLE threads (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    categorie_id INTEGER NULL,
    title TEXT,
    user_username TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (categorie_id) REFERENCES categories(id),
    FOREIGN KEY (user_username) REFERENCES users(username)
);