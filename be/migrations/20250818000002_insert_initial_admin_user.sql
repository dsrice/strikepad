-- Insert initial admin user
-- Default credentials: login_id = 'admin', password = 'admin123' (hashed)
-- NOTE: Change the default password after first login in production environment

INSERT INTO "public"."admin_users"
("login_id", "password", "name", "created_at", "updated_at", "is_deleted")
VALUES ('admin', '$2a$10$BlW9F2BF4Q0AJ9OviEjlK.WDNCqHX7RzkgDVvaXx7dnrRDpLInT0O', 'System Administrator',
        CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, false)
ON CONFLICT
    ("login_id")
    DO NOTHING;

-- Add comment about the default user
COMMENT ON TABLE "public"."admin_users" IS '管理者情報 - 初期管理者: login_id=admin, password=admin123';

-- Create unique index on login_id if it doesn't exist
CREATE UNIQUE INDEX IF NOT EXISTS "idx_admin_users_login_id" ON "public"."admin_users" ("login_id") WHERE "is_deleted" = false;