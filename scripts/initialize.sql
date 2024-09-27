CREATE TABLE IF NOT EXISTS sessions (username text primary key, token text, expires timestamp);
CREATE TABLE IF NOT EXISTS users (username text primary key, password text);
CREATE TABLE IF NOT EXISTS posts (title text primary key, link text, createdAt timestamp, createdBy username);
