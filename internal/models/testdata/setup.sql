CREATE TABLE users(id integer primary key, username TEXT NOT NULL UNIQUE, hashed_password TEXT NOT NULL, created DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP);
INSERT INTO users values(NULL, "admin", "cheesecake123", "2024");
