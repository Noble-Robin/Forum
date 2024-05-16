CREATE TABLE groups (
    ID TEXT PRIMARY KEY,
    groupname TEXT,
    creation TIMESTAMP
);

CREATE TABLE group_members (
    group_id TEXT,
    user_id INTEGER,
    FOREIGN KEY (group_id) REFERENCES groups(ID),
    FOREIGN KEY (user_id) REFERENCES users(ID)
);