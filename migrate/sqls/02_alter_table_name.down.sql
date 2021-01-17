create table user_dg_tmp
(
	uid bigint not null primary key,
	user_name VARCHAR(32) not null,
    is_block boolean not null,
    max_account tinyint not null,
    n_account tinyint not null,
    is_admin boolean not null,
    allow_polling boolean not null
);

insert into user_dg_tmp(uid, user_name, is_block, max_account, n_account, is_admin, allow_polling) select user.uid, user_name, is_block, max_account, 0, is_admin, allow_polling from user join permission on user.uid=permission.uid;

-- todo: update from
-- update user_dg_tmp
-- set n_account=acct.n
-- from (select uid, count(uid) as n from account group by uid) as acct
-- where user_dg_tmp.uid=acct.uid;
update user_dg_tmp set n_account=(select count(uid) from account where user_dg_tmp.uid=account.uid);

drop table user;

drop table permission;

alter table user_dg_tmp rename to user;

create index idx_user_name on user (user_name);

alter table status rename to runtime;
