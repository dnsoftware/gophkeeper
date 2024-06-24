DROP INDEX IF EXISTS login_index;

CREATE UNIQUE INDEX login_index
    ON users (login);
