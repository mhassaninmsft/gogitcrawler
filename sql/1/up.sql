-- up.sql

CREATE TABLE IF NOT EXISTS contributors (
    id INTEGER PRIMARY KEY,
    login VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS repos (
    id INTEGER PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    full_name VARCHAR(255) NOT NULL,
    contributor_id INTEGER NOT NULL,
    FOREIGN KEY (contributor_id) REFERENCES contributors(id) ON DELETE CASCADE
);