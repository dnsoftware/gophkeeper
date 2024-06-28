CREATE TABLE metainfo
(
    id SERIAL PRIMARY KEY,
    entity_id INTEGER,
    title CHARACTER VARYING(1024),
    value CHARACTER VARYING(1024)

);

CREATE INDEX meta_entity_id_index ON metainfo (entity_id);

