-- Create "user_sessions" table
create table user_sessions (
                               id serial not null
    , user_id integer not null
    , access_token text not null
    , "refresh-token" text
    , expires_at timestamp not null
    , created_at timestamp default CURRENT_TIMESTAMP not null
    , updated_at timestamp default CURRENT_TIMESTAMP not null
    , is_deleted BOOLEAN default false not null
    , deleted_at timestamp
    , constraint user_sessions_PKC primary key (id)
) ;

alter table "user_sessions"
    add constraint "user_sessions_FK1" foreign key ("user_id") references "users"("id")
        on delete cascade
        on update cascade;

comment on table user_sessions is 'ユーザーセッション情報';
comment on column user_sessions.id is 'セッションID:セッションID';
comment on column user_sessions.user_id is 'ユーザーID:ユーザーID';
comment on column user_sessions.access_token is 'アクセストークン:アクセストークン';
comment on column user_sessions."refresh-token" is 'リフレッシュトークン:リフレッシュトークン';
comment on column user_sessions.expires_at is '有効期限:有効期限';
comment on column user_sessions.created_at is '作成日';
comment on column user_sessions.updated_at is '更新日';
comment on column user_sessions.is_deleted is '削除フラグ';
comment on column user_sessions.deleted_at is '削除日';