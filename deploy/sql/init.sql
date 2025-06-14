use clide;

-- 用户表
create table if not exists c_users
(
    id            bigint auto_increment primary key comment 'id',
    user_account  varchar(256)                                    not null comment '账号',
    user_password varchar(512)                                    not null comment '密码',
    user_avatar   varchar(1024)                                   null comment '头像',
    user_profile  varchar(512)                                    null comment '简介',
    user_role     enum ('user','admin') default 'user'            not null comment '权限',
    edit_time     datetime              default current_timestamp not null comment '编辑时间',
    create_time   datetime              default current_timestamp not null comment '创建时间',
    update_time   datetime              default current_timestamp not null on update current_timestamp comment '更新时间',
    is_delete     tinyint               default 0                 not null comment '是否删除',
    index idx_user_account (user_account)
) comment '用户' collate = utf8mb4_unicode_ci;