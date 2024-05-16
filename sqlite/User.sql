CREATE TABLE users (
    ID INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT,
    name TEXT,
    email TEXT,
    password TEXT
);

CREATE TABLE user_conversation (
    user_id INTEGER,
    conversation_id TEXT,
    FOREIGN KEY (user_id) REFERENCES users(ID),
    FOREIGN KEY (conversation_id) REFERENCES conversations(ID)
);