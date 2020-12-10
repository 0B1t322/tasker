CREATE TABLE IF NOT EXISTS tasks (
    id          INTEGER PRIMARY KEY ,
    user_id                     INT,
    name                       TEXT,
    description                TEXT,
    create_time                TEXT,
    done                       BOOL
);

CREATE TABLE IF NOT EXISTS users (
    id      INTEGER PRIMARY KEY ,
    token                   TEXT
);

CREATE TABLE IF NOT EXISTS archive_tasks (
    id          INTEGER PRIMARY KEY ,
    user_id                     INT,
    name                       TEXT,
    description                TEXT,
    create_time                TEXT,
    done                       BOOL
);