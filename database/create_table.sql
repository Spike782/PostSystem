-- 创建数据库（如果不存在）
CREATE DATABASE IF NOT EXISTS post;

-- 创建用户
CREATE USER IF NOT EXISTS 'tester'@'localhost' IDENTIFIED BY '123456';

-- 授权用户访问 post 数据库
GRANT ALL PRIVILEGES ON post.* TO 'tester'@'localhost';
FLUSH PRIVILEGES;

-- 使用数据库
USE post;

-- 建表
create table if not exists user(
    id int auto_increment comment '用户id，自增',
    name varchar(20) not null comment '用户名',
    password char(32) not null comment '密码,md5',
    create_time datetime default current_timestamp comment '用户注册时间',
    update_time datetime default current_timestamp on update current_timestamp comment '最后修改时间',
    primary key (id),
    unique key idx_name(name)
)default charset=utf8mb4 comment '用户信息';

create table if not exists news(
    id int auto_increment comment '帖子id，自增',
    user_id int not null comment '发布者id',
    title varchar(100) not null comment '标题',
    article text not null comment '正文',
    create_time datetime default current_timestamp comment '发布时间',
    update_time datetime default current_timestamp on update current_timestamp comment '最后修改时间',
    delete_time datetime default null comment '删除时间',
    primary key (id),
    key idx_user (user_id)

)default charset=utf8mb4 comment '帖子';