CREATE TABLE Categories (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    description TEXT,
    user_username TEXT,
    post INTEGER DEFAULT 0,
    FOREIGN KEY (user_username) REFERENCES users(username)
);

INSERT INTO Categories (title, description, post) VALUES ('vinyl', 'Disque se jouant sur un tourne disque', 14);
INSERT INTO Categories (title, description, post) VALUES ('rapFR', 'Rap Français', 4);
INSERT INTO Categories (title, description, post) VALUES ('RapUS', 'Rap Américain', 99);
INSERT INTO Categories (title, description, post) VALUES ('Rock', 'Musique Rock', 23);
INSERT INTO Categories (title, description, post) VALUES ('Pop', 'Musique Pop', 2);
INSERT INTO Categories (title, description, post) VALUES ('Jazz', 'Musique Jazz', 4);
INSERT INTO Categories (title, description, post) VALUES ('Blues', 'Musique Blues',0);
INSERT INTO Categories (title, description, post) VALUES ('House', 'Musique House',0);
INSERT INTO Categories (title, description, post) VALUES ('Reggae', 'Musique Reggae',0);
INSERT INTO Categories (title, description, post) VALUES ('Country', 'Musique Country',0);
INSERT INTO Categories (title, description, post) VALUES ('Indie', 'Musique Indie',0);
INSERT INTO Categories (title, description, post) VALUES ('RnB', 'Rhythm and Blues',0);
INSERT INTO Categories (title, description, post) VALUES ('Soul', 'Musique Soul',0);
