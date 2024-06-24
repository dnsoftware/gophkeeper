CREATE TABLE entities
(
    id SERIAL PRIMARY KEY,
    user_id INTEGER,
    etype CHARACTER VARYING(64),
    created_at timestamp,
    updated_at timestamp

);

CREATE INDEX entity_user_id_index ON entities (user_id);
CREATE INDEX entity_etype_index ON entities (etype);

