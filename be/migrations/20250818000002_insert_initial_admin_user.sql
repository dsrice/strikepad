-- Insert initial admin user
-- Default credentials: login_id = 'admin', password = 'admin123' (hashed)
-- NOTE: Change the default password after first login in production environment

-- Create unique index on login_id if it doesn't exist (before inserting data)
CREATE UNIQUE INDEX IF NOT EXISTS "idx_admin_users_login_id" ON "public"."admin_users" ("login_id") WHERE "is_deleted" = false;

-- Insert initial admin user using INSERT WHERE NOT EXISTS pattern
INSERT INTO "public"."admin_users"
("login_id", "password", "name", "created_at", "updated_at", "is_deleted")
SELECT 'admin',
       '$2a$10$5PfHOyND.OIvQ6B3H/kauuswq5KCgGE3TCbzUbIoX850iFWbrKV1e',
       'System Administrator',
       CURRENT_TIMESTAMP,
       CURRENT_TIMESTAMP,
       false
WHERE NOT EXISTS (SELECT 1
                  FROM "public"."admin_users"
                  WHERE "login_id" = 'admin'
                    AND "is_deleted" = false);

-- Add comment about the default user
COMMENT ON TABLE "public"."admin_users" IS '管理者情報 - 初期管理者: login_id=admin, password=admin123';