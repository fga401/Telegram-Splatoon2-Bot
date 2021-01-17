create table account (
    uid bigint not null,
    session_token varchar(512) not null,
    tag varchar(64) not null,
    primary key (uid, tag)
);

create table user (
    uid bigint not null primary key,
    user_name varchar(32) not null,
    is_block boolean not null,
    max_account tinyint not null,
    n_account tinyint not null,
    is_admin boolean not null,
    allow_polling boolean not null
);

create table runtime (
    uid bigint not null primary key,
    session_token varchar(512) not null,
    iksm character(40) not null,
    language varchar(10) not null,
    timezone int not null
);

create index idx_user_name on user(user_name);