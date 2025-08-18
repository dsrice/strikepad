-- Create "makers" table
-- メーカー情報を管理するテーブル

-- メーカー情報
-- * RestoreFromTempTable
create table "makers"
(
    "id"         serial                              not null,
    "name"       character varying(100)              not null,
    "created_at" timestamp default CURRENT_TIMESTAMP not null,
    "updated_at" timestamp default CURRENT_TIMESTAMP not null,
    "is_deleted" BOOLEAN   default false             not null,
    "deleted_at" timestamp,
    constraint "makers_PKC" primary key ("id")
);

-- Set comments
comment on table "makers" is 'メーカー情報';
comment on column "makers"."id" is 'ID';
comment on column "makers"."name" is 'メーカー名';
comment on column "makers"."created_at" is '作成日';
comment on column "makers"."updated_at" is '更新日';
comment on column "makers"."is_deleted" is '削除フラグ';
comment on column "makers"."deleted_at" is '削除日';