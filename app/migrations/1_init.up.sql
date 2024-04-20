CREATE TABLE IF NOT EXISTS users
(
    id INTEGER PRIMARY KEY,
    login TEXT UNIQUE NOT NULL,
    pass_hash BLOB NOT NULL
);

CREATE TABLE IF NOT EXISTS apps
(
    id INTEGER PRIMARY KEY,
    name TEXT UNIQUE NOT NULL,
    secret TEXT UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS expressions (
    id INTEGER PRIMARY KEY,
    expression TEXT NOT NULL,
    ttdo TEXT NOT NULL,
    status TEXT NOT NULL,
    result INTEGER DEFAULT -1 NOT NULL,
    user_id INTEGER NOT NULL REFERENCES users(id) default -1
);

CREATE TABLE IF NOT EXISTS arithmetics (
    sign TEXT NOT NULL,
    ttdo TEXT NOT NULL,
    user_id INTEGER NOT NULL REFERENCES users(id) default -1
);