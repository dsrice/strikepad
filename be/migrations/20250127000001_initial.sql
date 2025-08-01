-- Create "users" table
create table users (
                       id serial not null
    , provider_type character varying(20) not null
    , provider_user_id character varying(255)
    , email character varying(255)
    , display_name character varying(100) not null
    , password_hash character varying(255)
    , email_verified BOOLEAN default false not null
    , created_at timestamp default CURRENT_TIMESTAMP not null
    , updated_at timestamp default CURRENT_TIMESTAMP not null
    , is_deleted BOOLEAN default false not null
    , deleted_at timestamp
    , constraint users_PKC primary key (id)
) ;

comment on table users is 'ユーザー情報';
comment on column users.id is 'ID:ID';
comment on column users.provider_type is 'プロバイダー種別:プロバイダー種別';
comment on column users.provider_user_id is 'プロバイダーユーザーID:プロバイダーユーザーID';
comment on column users.email is 'Eメール:Eメール';
comment on column users.display_name is '表示名:表示名';
comment on column users.password_hash is 'パスワードハッシュ:パスワードハッシュ';
comment on column users.email_verified is 'メール利用フラグ:メール利用フラグ';
comment on column users.created_at is '作成日';
comment on column users.updated_at is '更新日';
comment on column users.is_deleted is '削除フラグ';
comment on column users.deleted_at is '削除日';