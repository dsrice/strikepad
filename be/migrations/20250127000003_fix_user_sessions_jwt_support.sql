-- Fix user_sessions table structure for JWT support
-- Drop the problematic refresh-token column with hyphen
ALTER TABLE user_sessions
    DROP COLUMN "refresh-token";

-- Add new columns for JWT token support  
ALTER TABLE user_sessions
    ADD COLUMN refresh_token text;
ALTER TABLE user_sessions
    ADD COLUMN access_token_expires_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP;
ALTER TABLE user_sessions
    ADD COLUMN refresh_token_expires_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP;

-- Update comments
COMMENT ON COLUMN user_sessions.refresh_token IS 'リフレッシュトークン:リフレッシュトークン';
COMMENT ON COLUMN user_sessions.access_token_expires_at IS 'アクセストークン有効期限:アクセストークン有効期限';
COMMENT ON COLUMN user_sessions.refresh_token_expires_at IS 'リフレッシュトークン有効期限:リフレッシュトークン有効期限';
COMMENT ON COLUMN user_sessions.expires_at IS '非推奨：代わりにaccess_token_expires_atを使用';