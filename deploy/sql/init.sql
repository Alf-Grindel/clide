create database if not exists clide;

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

-- 图片表
create table if not exists c_pictures
(
    id           bigint auto_increment primary key comment 'id',
    url          varchar(512)                       not null comment '图片url',
    pic_name     varchar(128)                       not null comment '图片名称',
    introduction varchar(512)                       null comment '简介',
    category     varchar(64)                        null comment '分类',
    tags         varchar(512)                       null comment '标签（JSON数组',
    pic_size     bigint                             null comment '图片体积',
    pic_width    int                                null comment '图片宽度',
    pic_height   int                                null comment '图片高度',
    pic_scale    double                             null comment '图片宽高比例',
    pic_format   varchar(32)                        null comment '图片格式',
    user_id      bigint                             not null comment '创建用户id',
    edit_time    datetime default current_timestamp not null comment '编辑时间',
    create_time  datetime default current_timestamp not null comment '创建时间',
    update_time  datetime default current_timestamp not null on update current_timestamp comment '更新时间',
    is_delete    tinyint  default 0                 not null comment '是否删除',
    index idx_pic_name (pic_name),
    index idx_pic_introduction (introduction),
    index idx_category (category),
    index idx_tags (tags),
    index idx_user_id (user_id)
) comment '图片', collate = utf8mb4_unicode_ci;