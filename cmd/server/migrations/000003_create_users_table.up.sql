CREATE TABLE users
(
    id SERIAL PRIMARY KEY,
    login CHARACTER VARYING(256),
    password CHARACTER VARYING(256),
    salt CHARACTER VARYING(64),
    created_at timestamp

);

CREATE INDEX login_index ON users (login);
