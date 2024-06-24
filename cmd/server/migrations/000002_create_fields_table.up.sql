CREATE TABLE fields
(
    id SERIAL PRIMARY KEY,
    etype CHARACTER VARYING(64),
    name CHARACTER VARYING(256),
    ftype CHARACTER VARYING(64),
    validate_rules CHARACTER VARYING(512),
    validate_messages CHARACTER VARYING(512)

);

CREATE INDEX etype_index ON fields (etype);
CREATE INDEX ftype_index ON fields (ftype);