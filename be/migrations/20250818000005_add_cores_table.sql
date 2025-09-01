-- コア情報
-- * RestoreFromTempTable
create table "cores"
(
    "id"            serial                              not null,
    "name"          character varying(200)              not null,
    "rg"            real                                not null,
    "delta_rg"      real                                not null,
    "init_diff"     real,
    "symmetry_flag" BOOLEAN not null,
    "maker_id"      integer                             not null,
    "created_at"    timestamp default CURRENT_TIMESTAMP not null,
    "updated_at"    timestamp default CURRENT_TIMESTAMP not null,
    "is_deleted"    BOOLEAN   default false             not null,
    "deleted_at"    timestamp,
    constraint "cores_PKC" primary key ("id")
);

-- Set comments
comment on table "cores" is 'コア情報';
comment on column "cores"."id" is 'ID';
comment on column "cores"."name" is 'コア名';
comment on column "cores"."rg" is 'RG';
comment on column "cores"."delta_rg" is 'ΔRG';
comment on column "cores"."init_diff" is 'InitDiff';
comment on column "cores"."symmetry_flag" is '対称フラグ';
comment on column "cores"."maker_id" is 'メーカーID';
comment on column "cores"."created_at" is '作成日';
comment on column "cores"."updated_at" is '更新日';
comment on column "cores"."is_deleted" is '削除フラグ';
comment on column "cores"."deleted_at" is '削除日';

-- Create indexes
alter table "cores"
    add constraint "cores_FK1" foreign key ("maker_id") references "makers" ("id")
        on delete cascade
        on update cascade;