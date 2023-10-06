-- Active: 1695894389485@@127.0.0.1@3306@infoflow_user

create database infoflow_user;

use infoflow_user;

CREATE TABLE
    `user` (
        `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
        `mtime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后修改时间',
        `ctime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
        `username` varchar(32) NOT NULL DEFAULT '' COMMENT '用户名',
        `avatar` varchar(256) NOT NULL DEFAULT '' COMMENT '头像',
        `phone` varchar(128) NOT NULL DEFAULT '' COMMENT '手机号',
        `password` varchar(128) NOT NULL DEFAULT '' COMMENT '密码',
        PRIMARY KEY (`id`),
        KEY `ix_mtime` (`mtime`)
    ) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_bin COMMENT = '用户表';