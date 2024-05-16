CREATE TABLE messages (
    id_message TEXT PRIMARY KEY,
    id_conv TEXT,
    id_user INTEGER, 
    contenu TEXT,
    creation TIMESTAMP,
    FOREIGN KEY (id_conv) REFERENCES conversations(ID),
    FOREIGN KEY (id_user) REFERENCES users(ID)
);