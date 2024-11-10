CREATE TABLE IF NOT EXISTS sessions (username text primary key, token text, expires timestamp);
CREATE TABLE IF NOT EXISTS users (username text primary key, password text, salt text, roles text);
CREATE TABLE IF NOT EXISTS posts (id text primary key, title text, post text, createdAt timestamp, createdBy username);
CREATE TABLE IF NOT EXISTS availableRoles (roles text primary key);
INSERT INTO availableRoles values ('read'),('create');
