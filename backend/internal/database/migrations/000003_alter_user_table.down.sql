ALTER TABLE users (
    DROP COLUMN username,
    DROP COLUMN permission_role,
    ADD COLUMN password_hash VARCHAR(255) NOT NULL
);