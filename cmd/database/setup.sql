CREATE TABLE users(
    id INTEGER PRIMARY KEY, 
    username TEXT NOT NULL UNIQUE, 
    hashed_password TEXT NOT NULL, 
    created DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP);
CREATE TABLE sessions (
	token TEXT PRIMARY KEY,
	data BLOB NOT NULL,
	expiry REAL NOT NULL
);
CREATE INDEX sessions_expiry_idx ON sessions(expiry);