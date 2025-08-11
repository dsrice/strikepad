-- Remove the unused expires_at column that conflicts with new JWT structure
ALTER TABLE user_sessions
    DROP COLUMN expires_at;