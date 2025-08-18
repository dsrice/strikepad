-- ボール情報
-- * RestoreFromTempTable
create table "balls"
(
    "id"         serial                              not null,
    "name"       character varying(100)              not null,
    "maker_id"   integer                             not null,
    "core_id"    integer                             not null,
    "cover_id"   integer                             not null,
    "url"        character varying(200)              not null,
    "created_at" timestamp default CURRENT_TIMESTAMP not null,
    "updated_at" timestamp default CURRENT_TIMESTAMP not null,
    "is_deleted" BOOLEAN   default false             not null,
    "deleted_at" timestamp,
    constraint "balls_PKC" primary key ("id")
);

-- Set comments
comment on table "balls" is 'ボール情報';
comment on column "balls"."id" is 'ID';
comment on column "balls"."name" is 'ボール名';
comment on column "balls"."maker_id" is 'メーカーID';
comment on column "balls"."core_id" is 'コアID';
comment on column "balls"."cover_id" is 'カバーID';
comment on column "balls"."url" is 'ボールURL';
comment on column "balls"."created_at" is '作成日';
comment on column "balls"."updated_at" is '更新日';
comment on column "balls"."is_deleted" is '削除フラグ';
comment on column "balls"."deleted_at" is '削除日';

-- Create indexes
alter table "balls"
    add constraint "balls_FK1" foreign key ("cover_id") references "covers" ("id")
        on delete cascade
        on update cascade;

alter table "balls"
    add constraint "balls_FK2" foreign key ("core_id") references "cores" ("id")
        on delete cascade
        on update cascade;

alter table "balls"
    add constraint "balls_FK3" foreign key ("maker_id") references "makers" ("id")
        on delete cascade
        on update cascade;