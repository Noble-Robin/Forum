CREATE TABLE Categories (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    description TEXT,
    view INTEGER DEFAULT 0
);

INSERT INTO Categories (title, description, view) VALUES ('vinyl', 'Disque se jouant sur un tourne disque', 14);
INSERT INTO Categories (title, description, view) VALUES ('rapFR', 'Rap Français', 4);
INSERT INTO Categories (title, description, view) VALUES ('Rap US', 'Rap Américain', 99);
INSERT INTO Categories (title, description, view) VALUES ('Rock', 'Musique Rock', 23);
INSERT INTO Categories (title, description, view) VALUES ('Pop', 'Musique Pop', 2);
INSERT INTO Categories (title, description, view) VALUES ('Jazz', 'Musique Jazz', 4);
INSERT INTO Categories (title, description, view) VALUES ('Blues', 'Musique Blues',0);
INSERT INTO Categories (title, description, view) VALUES ('House', 'Musique House',0);
INSERT INTO Categories (title, description, view) VALUES ('Reggae', 'Musique Reggae',0);
INSERT INTO Categories (title, description, view) VALUES ('Country', 'Musique Country',0);
INSERT INTO Categories (title, description, view) VALUES ('Indie', 'Musique Indie',0);
INSERT INTO Categories (title, description, view) VALUES ('RnB', 'Rhythm and Blues',0);
INSERT INTO Categories (title, description, view) VALUES ('Soul', 'Musique Soul',0);
