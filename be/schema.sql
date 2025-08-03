-- Users table
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    provider_type VARCHAR(20) NOT NULL,
    provider_user_id VARCHAR(255),
    email VARCHAR(255),
    display_name VARCHAR(100) NOT NULL,
    password_hash VARCHAR(255),
    email_verified BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    is_deleted BOOLEAN NOT NULL DEFAULT false,
    deleted_at TIMESTAMP
);

COMMENT ON TABLE users IS 'ユーザー情報';
COMMENT ON COLUMN users.id IS 'ID:ID';
COMMENT ON COLUMN users.provider_type IS 'プロバイダー種別:プロバイダー種別';
COMMENT ON COLUMN users.provider_user_id IS 'プロバイダーユーザーID:プロバイダーユーザーID';
COMMENT ON COLUMN users.email IS 'Eメール:Eメール';
COMMENT ON COLUMN users.display_name IS '表示名:表示名';
COMMENT ON COLUMN users.password_hash IS 'パスワードハッシュ:パスワードハッシュ';
COMMENT ON COLUMN users.email_verified IS 'メール利用フラグ:メール利用フラグ';
COMMENT ON COLUMN users.created_at IS '作成日';
COMMENT ON COLUMN users.updated_at IS '更新日';
COMMENT ON COLUMN users.is_deleted IS '削除フラグ';
COMMENT ON COLUMN users.deleted_at IS '削除日';

-- User sessions table
CREATE TABLE user_sessions (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    session_token VARCHAR(255) NOT NULL UNIQUE,
    refresh_token VARCHAR(255),
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    is_deleted BOOLEAN NOT NULL DEFAULT false,
    deleted_at TIMESTAMP,
    CONSTRAINT fk_user_sessions_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

COMMENT ON TABLE user_sessions IS 'ユーザーセッション情報';
COMMENT ON COLUMN user_sessions.id IS 'セッションID:セッションID';
COMMENT ON COLUMN user_sessions.user_id IS 'ユーザーID:ユーザーID';
COMMENT ON COLUMN user_sessions.session_token IS 'セッショントークン:セッショントークン';
COMMENT ON COLUMN user_sessions.refresh_token IS 'リフレッシュトークン:リフレッシュトークン';
COMMENT ON COLUMN user_sessions.expires_at IS '有効期限:有効期限';
COMMENT ON COLUMN user_sessions.created_at IS '作成日';
COMMENT ON COLUMN user_sessions.updated_at IS '更新日';
COMMENT ON COLUMN user_sessions.is_deleted IS '削除フラグ';
COMMENT ON COLUMN user_sessions.deleted_at IS '削除日';

-- Create indexes
CREATE INDEX idx_user_sessions_user_id ON user_sessions(user_id);
CREATE INDEX idx_user_sessions_session_token ON user_sessions(session_token);
CREATE INDEX idx_user_sessions_expires_at ON user_sessions(expires_at);
CREATE INDEX idx_user_sessions_is_deleted ON user_sessions(is_deleted);