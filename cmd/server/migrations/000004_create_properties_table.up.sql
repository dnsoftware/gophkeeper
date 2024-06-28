CREATE TABLE properties
(
    id SERIAL PRIMARY KEY,
    entity_id INTEGER,
    field_id INTEGER,
    value CHARACTER VARYING(1024)

);

CREATE INDEX entity_id_index ON properties (entity_id);
CREATE INDEX field_id_index ON properties (field_id);
