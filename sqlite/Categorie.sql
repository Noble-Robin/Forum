CREATE TABLE Categories (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    description TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO Categories (title, description) VALUES ('vinyl', 'Disque se jouant sur un tourne disque');
INSERT INTO Categories (title, description) VALUES ('rapFR', 'Rap Français');
INSERT INTO Categories (title, description) VALUES ('Rap US', 'Rap Américain');
INSERT INTO Categories (title, description) VALUES ('Rock', 'Musique Rock');
INSERT INTO Categories (title, description) VALUES ('Pop', 'Musique Pop');
INSERT INTO Categories (title, description) VALUES ('Jazz', 'Musique Jazz');
INSERT INTO Categories (title, description) VALUES ('Blues', 'Musique Blues');
INSERT INTO Categories (title, description) VALUES ('House', 'Musique House');
INSERT INTO Categories (title, description) VALUES ('Reggae', 'Musique Reggae');
INSERT INTO Categories (title, description) VALUES ('Country', 'Musique Country');
INSERT INTO Categories (title, description) VALUES ('Indie', 'Musique Indie');
INSERT INTO Categories (title, description) VALUES ('RnB', 'Rhythm and Blues');
INSERT INTO Categories (title, description) VALUES ('Soul', 'Musique Soul');
