CREATE TABLE threads (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    categorie_title TEXT,
    title TEXT,
    user_username TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (categorie_title) REFERENCES categories(title),
    FOREIGN KEY (user_username) REFERENCES users(username)
);