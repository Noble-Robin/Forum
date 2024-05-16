CREATE TABLE conversations (
    ID TEXT PRIMARY KEY,
    creation TIMESTAMP
);

CREATE TABLE conversation_participants (
    conversation_id TEXT,
    user_id INTEGER,
    FOREIGN KEY (conversation_id) REFERENCES conversations(ID),
    FOREIGN KEY (user_id) REFERENCES users(ID)
);