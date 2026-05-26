ALTER TABLE users
    ADD COLUMN username VARCHAR(255) UNIQUE NOT NULL,
    ADD COLUMN permission_role VARCHAR(50) NOT NULL DEFAULT 'user',
    DROP COLUMN password_hash