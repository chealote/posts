CREATE TABLE IF NOT EXISTS sessions (username text primary key, token text, expires timestamp);
CREATE TABLE IF NOT EXISTS users (username text primary key, password text, salt text);
CREATE TABLE IF NOT EXISTS posts (id text primary key, title text, post text, createdAt timestamp, createdBy username);
