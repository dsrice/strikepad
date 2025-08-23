-- Create "admin_users" table
CREATE TABLE "public"."admin_users"
(
    "id"         serial                 NOT NULL,
    "login_id"   character varying(30)  NOT NULL,
    "password"   character varying(200) NOT NULL,
    "name"       character varying(30)  NOT NULL,
    "created_at" timestamp              NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" timestamp              NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "is_deleted" boolean                NOT NULL DEFAULT false,
    "deleted_at" timestamp              NULL,
    CONSTRAINT "admin_users_PKC" PRIMARY KEY ("id")
);

-- Set comment to table: "admin_users"
COMMENT ON TABLE "public"."admin_users" IS '管理者情報';
-- Set comment to column: "id" on table: "admin_users"
COMMENT ON COLUMN "public"."admin_users"."id" IS 'ID';
-- Set comment to column: "login_id" on table: "admin_users"
COMMENT ON COLUMN "public"."admin_users"."login_id" IS 'ログインID';
-- Set comment to column: "password" on table: "admin_users"
COMMENT ON COLUMN "public"."admin_users"."password" IS 'パスワード';
-- Set comment to column: "name" on table: "admin_users"
COMMENT ON COLUMN "public"."admin_users"."name" IS 'ユーザー名';
-- Set comment to column: "created_at" on table: "admin_users"
COMMENT ON COLUMN "public"."admin_users"."created_at" IS '作成日';
-- Set comment to column: "updated_at" on table: "admin_users"
COMMENT ON COLUMN "public"."admin_users"."updated_at" IS '更新日';
-- Set comment to column: "is_deleted" on table: "admin_users"
COMMENT ON COLUMN "public"."admin_users"."is_deleted" IS '削除フラグ';
-- Set comment to column: "deleted_at" on table: "admin_users"
COMMENT ON COLUMN "public"."admin_users"."deleted_at" IS '削除日';