CREATE TABLE metainfo
(
    id SERIAL PRIMARY KEY,
    entity_id INTEGER,
    title CHARACTER VARYING(256),
    value CHARACTER VARYING(256)

);

CREATE INDEX meta_entity_id_index ON metainfo (entity_id);

