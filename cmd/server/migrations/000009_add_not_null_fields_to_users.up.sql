ALTER TABLE users
    ALTER COLUMN login SET NOT NULL,
    ALTER COLUMN password SET NOT NULL,
    ALTER COLUMN salt SET NOT NULL
