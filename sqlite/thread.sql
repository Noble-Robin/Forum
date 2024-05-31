CREATE TABLE threads (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    categorie_id INTEGER,
    title TEXT NOT NULL,
    user_id INTEGER,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (categorie_id) REFERENCES category(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);
