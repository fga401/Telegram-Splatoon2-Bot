create table permission
(
	uid BIGINT not null primary key,
	is_block BOOLEAN not null,
	max_account TINYINT not null,
	is_admin BOOLEAN not null,
	allow_polling BOOLEAN not null
);

insert into permission(uid, is_block, max_account, is_admin, allow_polling) select uid, is_block, max_account, is_admin, allow_polling from user;

create table user_dg_tmp
(
	uid BIGINT not null primary key,
	user_name VARCHAR(32) not null
);

insert into user_dg_tmp(uid, user_name) select uid, user_name from user;

drop table user;

alter table user_dg_tmp rename to user;

create index idx_user_name on user (user_name);

alter table runtime rename to status;
