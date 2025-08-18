-- カバー情報
-- * RestoreFromTempTable
create table "covers"
(
    "id"            serial                              not null,
    "name"          character varying(200)              not null,
    "material_type" integer                             not null,
    "maker_id"      integer                             not null,
    "rank"          integer                             not null,
    "created_at"    timestamp default CURRENT_TIMESTAMP not null,
    "updated_at"    timestamp default CURRENT_TIMESTAMP not null,
    "is_deleted"    BOOLEAN   default false             not null,
    "deleted_at"    timestamp,
    constraint "covers_PKC" primary key ("id")
);

-- Set comments
comment on table "covers" is 'カバー情報';
comment on column "covers"."id" is 'ID';
comment on column "covers"."name" is 'カバー名';
comment on column "covers"."material_type" is '材質種別';
comment on column "covers"."maker_id" is 'メーカーID';
comment on column "covers"."rank" is 'ランク';
comment on column "covers"."created_at" is '作成日';
comment on column "covers"."updated_at" is '更新日';
comment on column "covers"."is_deleted" is '削除フラグ';
comment on column "covers"."deleted_at" is '削除日';
-- Create indexes
alter table "covers"
    add constraint "covers_FK1" foreign key ("maker_id") references "makers" ("id")
        on delete cascade
        on update cascade;